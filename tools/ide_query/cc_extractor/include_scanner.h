#ifndef _TOOLS_IDE_QUERY_CC_EXTRACTOR_INCLUDE_SCANNER_H_
#define _TOOLS_IDE_QUERY_CC_EXTRACTOR_INCLUDE_SCANNER_H_

#include <string>
#include <utility>
#include <vector>

#include "clang/Tooling/CompilationDatabase.h"
#include "llvm/ADT/IntrusiveRefCntPtr.h"
#include "llvm/Support/Error.h"
#include "llvm/Support/VirtualFileSystem.h"

namespace tools::ide_query::cc_extractor {

// Returns absolute paths and contents for all the includes necessary for
// compiling source file in command.
llvm::Expected<std::vector<std::pair<std::string, std::string>>> ScanIncludes(
    const clang::tooling::CompileCommand &cmd,
    llvm::IntrusiveRefCntPtr<llvm::vfs::FileSystem> fs);

}  // namespace tools::ide_query::cc_extractor

#endif
