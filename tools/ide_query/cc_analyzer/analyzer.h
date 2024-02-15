#ifndef _TOOLS_IDE_QUERY_CC_EXTRACTOR_EXTRACTOR_H_
#define _TOOLS_IDE_QUERY_CC_EXTRACTOR_EXTRACTOR_H_

#include "ide_query.pb.h"

namespace tools::ide_query::cc_extractor {

// Scans the build graph and returns target names from the build graph to
// generate all the dependencies for the active files.
::ide_query::DepsResponse GetDeps(const ::ide_query::DepsRequest& req);

// Scans the sources and returns all the source files required for analyzing the
// active files.
::ide_query::IdeAnalysis GetBuildInputs(const ::ide_query::InputsRequest& req);

}  // namespace tools::ide_query::cc_extractor

#endif
