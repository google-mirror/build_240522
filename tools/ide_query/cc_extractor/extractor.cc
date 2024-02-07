#include "extractor.h"

#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "clang/Tooling/CompilationDatabase.h"
#include "clang/Tooling/JSONCompilationDatabase.h"
#include "extractor.pb.h"
#include "include_scanner.h"
#include "llvm/ADT/SmallString.h"
#include "llvm/ADT/StringRef.h"
#include "llvm/ADT/Twine.h"
#include "llvm/Support/Error.h"
#include "llvm/Support/Path.h"
#include "llvm/Support/VirtualFileSystem.h"

namespace tools::ide_query::cc_extractor {
namespace companion = cider::services::build::companion;
namespace {
llvm::Expected<std::unique_ptr<clang::tooling::CompilationDatabase>> LoadCompDB(
    llvm::StringRef comp_db_path) {
  std::string err;
  std::unique_ptr<clang::tooling::CompilationDatabase> db =
      clang::tooling::JSONCompilationDatabase::loadFromFile(
          comp_db_path, err, clang::tooling::JSONCommandLineSyntax::AutoDetect);
  if (!db) {
    return llvm::createStringError(llvm::inconvertibleErrorCode(),
                                   "Failed to load CDB: " + err);
  }
  // Provide some heuristic support for missing files.
  return inferMissingCompileCommands(std::move(db));
}
}  // namespace

companion::DepsResponse GetDeps(const companion::DepsRequest& req) {
  companion::DepsResponse results;
  auto db = LoadCompDB(req.state().comp_db_path());
  if (!db) {
    results.mutable_status()->set_code(companion::Status::FAILURE);
    results.mutable_status()->set_message(llvm::toString(db.takeError()));
    return results;
  }
  for (llvm::StringRef active_file : req.state().active_file_path()) {
    auto& result = *results.add_deps();

    llvm::SmallString<256> abs_file(req.state().repo_dir());
    llvm::sys::path::append(abs_file, active_file);
    auto cmds = db->get()->getCompileCommands(abs_file);
    if (cmds.empty()) {
      result.mutable_status()->set_code(companion::Status::FAILURE);
      result.mutable_status()->set_message(
          llvm::Twine("Can't find compile flags for file: ", abs_file).str());
      continue;
    }
    result.set_source_file(active_file.str());
    // TODO: Query ninja graph to figure out a minimal set of targets to build.
    result.add_build_target(cmds[0].Filename + "^");
  }
  return results;
}

companion::IdeAnalysis GetBuildInputs(const companion::InputsRequest& req) {
  auto db = LoadCompDB(req.state().comp_db_path());
  companion::IdeAnalysis results;
  if (!db) {
    results.mutable_status()->set_code(companion::Status::FAILURE);
    results.mutable_status()->set_message(llvm::toString(db.takeError()));
    return results;
  }
  auto workspace_path = llvm::sys::path::parent_path(req.state().repo_dir());
  results.set_build_artifact_root(req.state().out_dir());
  for (llvm::StringRef active_file : req.state().active_file_path()) {
    auto& result = *results.add_sources();

    llvm::SmallString<256> abs_file(req.state().repo_dir());
    llvm::sys::path::append(abs_file, active_file);
    auto cmds = db->get()->getCompileCommands(abs_file);
    if (cmds.empty()) {
      result.mutable_status()->set_code(companion::Status::FAILURE);
      result.mutable_status()->set_message(
          llvm::Twine("Can't find compile flags for file: ", abs_file).str());
      continue;
    }
    const auto& cmd = cmds.front();
    llvm::StringRef working_dir = cmd.Directory;
    if (!working_dir.consume_front(workspace_path) ||
        !working_dir.consume_front("/")) {
      result.mutable_status()->set_code(companion::Status::FAILURE);
      result.mutable_status()->set_message("Working dir outside workspace: " +
                                           cmd.Directory);
      continue;
    }
    auto includes =
        ScanIncludes(cmds.front(), llvm::vfs::createPhysicalFileSystem());
    if (!includes) {
      result.mutable_status()->set_code(companion::Status::FAILURE);
      result.mutable_status()->set_message(
          llvm::toString(includes.takeError()));
      continue;
    }

    result.set_working_dir(working_dir.str());
    result.set_path(active_file.str());
    for (auto& arg : cmd.CommandLine) {
      result.add_compiler_arguments(arg);
    }
    for (auto& [req_input, contents] : *includes) {
      llvm::StringRef req_input_ref(req_input);
      // We're only interested in generated files.
      if (!req_input_ref.consume_front(results.build_artifact_root())) continue;
      auto& genfile = *result.add_generated();
      genfile.set_path(req_input_ref.str());
      genfile.set_contents(std::move(contents));
    }
  }
  return results;
}
}  // namespace tools::ide_query::cc_extractor
