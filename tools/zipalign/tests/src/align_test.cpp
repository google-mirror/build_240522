#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <fstream>
#include <stdio.h>
#include <string>

#include <android-base/file.h>

using namespace android;

static std::string GetTestPath(const std::string& filename) {
  static std::string test_data_dir = android::base::GetExecutableDirectory() + "/tests/data/";
  return test_data_dir + filename;
}

TEST(Align, Unaligned) {
  const std::string src = GetTestPath("unaligned.zip");
  const std::string dst = GetTestPath("unaligned_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst.c_str(), 4, true, false);
  ASSERT_EQ(0, verified);
}

long GetFileSize(std::string filename) {
    struct stat stat_buf;
    int rc = stat(filename.c_str(), &stat_buf);
    return rc == 0 ? stat_buf.st_size : -1;
}

std::string GetFile(std::string filename) {
    long size = GetFileSize(filename);
    std::string buffer;
    buffer.reserve(size);

    std::ifstream t(filename);
    t.read(&buffer[0], size); 
    return buffer;
}

TEST(Align, DoubleAligment) {
  const std::string src = GetTestPath("unaligned.zip");
  const std::string tmp = GetTestPath("da_aligned.zip");
  const std::string dst = GetTestPath("da_d_aligner.zip");

  int processed = process(src.c_str(), tmp.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(tmp.c_str(), 4, true, false);
  ASSERT_EQ(0, verified);

  processed = process(tmp.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  verified = verify(dst.c_str(), 4, true, false);
  ASSERT_EQ(0, verified);

  // Nothing should have changed between tmp and dst.
  ASSERT_EQ(GetFileSize(tmp), GetFileSize(dst));
  ASSERT_EQ(GetFile(tmp), GetFile(dst));
}

// Align a zip featuring a hole at the beginning. The
// hole in the archive is a delete entry in the Central
// Directory.
TEST(Align, Holes) {
  const std::string src = GetTestPath("holes.zip");
  const std::string dst = GetTestPath("holes_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst.c_str(), 4, false, true);
  ASSERT_EQ(0, verified);
}

// Align a zip where LFH order and CD entries differ.
TEST(Align, DifferenteOrders) {
  const std::string src = GetTestPath("diffOrders.zip");
  const std::string dst = GetTestPath("diffOrders_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst.c_str(), 4, false, true);
  ASSERT_EQ(0, verified);
}
