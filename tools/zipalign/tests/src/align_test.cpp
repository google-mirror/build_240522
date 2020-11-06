#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <stdio.h>

using namespace android;

TEST(Align, Unaligned) {
  const char* src = "tests/data/unaligned.zip";
  const char* dst = "tests/data/unaligned_out.zip";

  int processed = process(src, dst, 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst, 4, true, false);
  ASSERT_EQ(0, verified);
}
