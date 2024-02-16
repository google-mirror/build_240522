#ifndef _TOOLS_IDE_QUERY_CC_ANALYZER_ANALYZER_H_
#define _TOOLS_IDE_QUERY_CC_ANALYZER_ANALYZER_H_

#include "ide_query.pb.h"

namespace tools::ide_query::cc_analyzer {

// Scans the build graph and returns target names from the build graph to
// generate all the dependencies for the active files.
::ide_query::DepsResponse GetDeps(::ide_query::RepoState state);

// Scans the sources and returns all the source files required for analyzing the
// active files.
::ide_query::IdeAnalysis GetBuildInputs(::ide_query::RepoState state);

}  // namespace tools::ide_query::cc_analyzer

#endif
