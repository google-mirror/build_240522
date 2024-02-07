#ifndef _TOOLS_IDE_QUERY_CC_EXTRACTOR_BUILTIN_HEADERS_H_
#define _TOOLS_IDE_QUERY_CC_EXTRACTOR_BUILTIN_HEADERS_H_

#include <cstddef>

struct FileToc {
  const char *name;
  const char *data;
};

const struct FileToc *builtin_headers_create();
size_t builtin_headers_size();

#endif
