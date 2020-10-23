#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <stdio.h>
#include <string>

#include <android-base/file.h>

using namespace android;

static std::string test_data_dir = android::base::GetExecutableDirectory() + "/tests/data";

TEST(Align, Unaligned) {
  const std::string src = test_data_dir + "unaligned.zip";
  const std::string dst = test_data_dir + "/unaligned_out.zip";
  int result = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, result);
}
