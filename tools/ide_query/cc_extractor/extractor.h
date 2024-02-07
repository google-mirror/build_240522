#ifndef _TOOLS_IDE_QUERY_CC_EXTRACTOR_EXTRACTOR_H_
#define _TOOLS_IDE_QUERY_CC_EXTRACTOR_EXTRACTOR_H_

#include "extractor.pb.h"

namespace tools::ide_query::cc_extractor {

// Scans the build graph and returns target names from the build graph to
// generate all the dependencies for the active files.
cider::services::build::companion::DepsResponse GetDeps(
    const cider::services::build::companion::DepsRequest& req);

// Scans the sources and returns all the source files required for analyzing the
// active files.
cider::services::build::companion::IdeAnalysis GetBuildInputs(
    const cider::services::build::companion::InputsRequest& req);

}  // namespace tools::ide_query::cc_extractor

#endif
