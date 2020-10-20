#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <stdio.h>

using namespace android;

TEST(Align, Unaligned) {
  const char* src = "tests/data/unaligned.zip";
  const char* dst = "tests/data/unaligned_out.zip";
  int result = process(src, dst, 4, true, false, 4096);
  ASSERT_EQ(0, result);
}

// Align a zip featuring a hole at the beginning. The
// hole in the archive is a delete entry in the Central
// Directory.
TEST(Align, Holes) {
  const char* src = "tests/data/holes.zip";
  const char* dst = "tests/data/holes_out.zip";

  int processed = process(src, dst, 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst, 4, false, true);
  ASSERT_EQ(0, verified);
}

// Align a zip where LFH order and CD entries differ.
TEST(Align, DifferenteOrders) {
  const char* src = "tests/data/diffOrders.zip";
  const char* dst = "tests/data/diffOrders_out.zip";

  int processed = process(src, dst, 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst, 4, false, true);
  ASSERT_EQ(0, verified);
}
