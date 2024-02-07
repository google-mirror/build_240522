// Driver for c++ extractor. Operates in two modes:
// - DEPS, scans build graph for active files and reports targets that need to
// be build for analyzing that file.
// - INPUTS, scans the source code for active files and returns all the sources
// required for analyzing that file.
//
// Uses stdin/stdout to take in requests and provide responses.
#include <unistd.h>

#include "extractor.h"
#include "extractor.pb.h"
#include "llvm/ADT/StringRef.h"
#include "llvm/Support/CommandLine.h"
#include "llvm/Support/InitLLVM.h"
#include "llvm/Support/TargetSelect.h"

namespace {
enum class OpMode {
  DEPS = 0,
  INPUTS = 1,
};
llvm::cl::opt<OpMode> mode{
    "mode",
    llvm::cl::values(clEnumValN(OpMode::DEPS, "deps",
                                "Figure out targets that need to be build"),
                     clEnumValN(OpMode::INPUTS, "inputs",
                                "Figure out generated files used")),
    llvm::cl::desc("Print the list of headers to insert and remove"),
};

cider::services::build::companion::IdeAnalysis ReturnError(
    llvm::StringRef message) {
  cider::services::build::companion::IdeAnalysis result;
  result.mutable_status()->set_code(
      cider::services::build::companion::Status::FAILURE);
  result.mutable_status()->set_message(message.str());
  return result;
}

}  // namespace

int main(int argc, char* argv[]) {
  llvm::InitializeAllTargetInfos();
  llvm::cl::ParseCommandLineOptions(argc, argv);

  cider::services::build::companion::IdeAnalysis result;
  switch (mode) {
    case OpMode::DEPS: {
      cider::services::build::companion::DepsRequest req;
      if (!req.ParseFromFileDescriptor(STDIN_FILENO)) {
        result = ReturnError("Failed to parse input.");
      } else if (!tools::ide_query::cc_extractor::GetDeps(req)
                      .SerializeToFileDescriptor(STDOUT_FILENO)) {
        result = ReturnError("Failed to write output.");
      }
      break;
    }
    case OpMode::INPUTS: {
      cider::services::build::companion::InputsRequest req;
      if (!req.ParseFromFileDescriptor(STDIN_FILENO)) {
        result = ReturnError("Failed to parse input.");
      } else if (!tools::ide_query::cc_extractor::GetBuildInputs(req)
                      .SerializeToFileDescriptor(STDOUT_FILENO)) {
        result = ReturnError("Failed to write output.");
      }
      break;
    }
  }

  return 0;
}
