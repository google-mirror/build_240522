#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <stdio.h>
#include <string>

using namespace android;

#include <android-base/file.h>


static std::string test_data_dir = android::base::GetExecutableDirectory() + "/tests/data";

TEST(Align, Unaligned) {
  const std::string src = test_data_dir + "unaligned.zip";
  const std::string dst = test_data_dir + "/unaligned_out.zip";
  int result = process(src.c_str(), dsti.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, result);
}
