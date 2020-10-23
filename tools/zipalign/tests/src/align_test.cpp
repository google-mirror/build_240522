#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <stdio.h>

using namespace android;

TEST(Align, AlignZipWithHoles) {
  const char* src = "tests/data/zip_with_holes.zip";
  const char* dst = "tests/data/test_align_with_holes.zip.tmp";
  int result = process(src, dst, 4, true, false, 4096);
  ASSERT_EQ(0, result);
}
