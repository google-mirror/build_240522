#include "builtin_headers.h"

#include <cstddef>

#include "clang_builtin_headers_resources.inc"

const struct FileToc *builtin_headers_create() { return kPackedFiles; }
size_t builtin_headers_size() {
  return sizeof(kPackedFiles) / sizeof(FileToc) - 1;
}
