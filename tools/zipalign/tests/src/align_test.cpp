#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "ZipAlign.h"

#include <filesystem>
#include <stdio.h>
#include <string>

#include <android-base/file.h>

using namespace android;
using namespace base;

// This load the whole file to memory so be careful!
static bool sameContent(const std::string& path1, const std::string& path2) {
    FILE* f1 = fopen(path1.c_str(), "r");
    if (f1 == nullptr) {
        printf("Unable to open %s\n", path1.c_str());
    }

    FILE* f2 = fopen(path2.c_str(), "r");
    if (f2 == nullptr) {
        printf("Unable to open %s\n", path2.c_str());
    }

    fseek(f1, 0, SEEK_END); // seek to end of file
    int f1Size = ftell(f1);

    fseek(f2, 0, SEEK_END); // seek to end of file
    int f2Size = ftell(f2);

    if (f1Size != f2Size) {
        return false;
    }

    char* f1Content = (char*)malloc(f1Size);
    fread(f1Content, f1Size, 1, f1);
    fclose(f1);

    char* f2Content = (char*)malloc(f2Size);
    fread(f2Content, f1Size, 1, f2);
    fclose(f2);

    bool diffCount = strcmp(f1Content, f2Content);

    delete f1Content;
    delete f2Content;
    return diffCount == 0;
}

static std::string GetTestPath(const std::string& filename) {
  static std::string test_data_dir = android::base::GetExecutableDirectory() + "/tests/data/";
  return test_data_dir + filename;
}

static std::string GetTempPath(const std::string& filename) {
  std::filesystem::path temp_path = std::filesystem::path(testing::TempDir());
  temp_path += filename;
  return temp_path.string();
}

TEST(Align, Unaligned) {
  const std::string src = GetTestPath("unaligned.zip");
  const std::string dst = GetTempPath("unaligned_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst.c_str(), 4, true, false);
  ASSERT_EQ(0, verified);
}

TEST(Align, DoubleAligment) {
  const std::string src = GetTestPath("unaligned.zip");
  const std::string tmp = GetTempPath("da_aligned.zip");
  const std::string dst = GetTempPath("da_d_aligner.zip");

  int processed = process(src.c_str(), tmp.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(tmp.c_str(), 4, true, false);
  ASSERT_EQ(0, verified);

  // Align the result of the previous run. Essentially double aligning.
  processed = process(tmp.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  verified = verify(dst.c_str(), 4, true, false);
  ASSERT_EQ(0, verified);

  // Nothing should have changed between tmp and dst.
  std::string tmp_content;
  ASSERT_EQ(true, ReadFileToString(tmp, &tmp_content));

  std::string dst_content;
  ASSERT_EQ(true, ReadFileToString(dst, &dst_content));

  ASSERT_EQ(tmp_content, dst_content);
}

// Align a zip featuring a hole at the beginning. The
// hole in the archive is a delete entry in the Central
// Directory.
TEST(Align, Holes) {
  const std::string src = GetTestPath("holes.zip");
  const std::string dst = GetTempPath("holes_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst.c_str(), 4, false, true);
  ASSERT_EQ(0, verified);
}

// Align a zip where LFH order and CD entries differ.
TEST(Align, DifferenteOrders) {
  const std::string src = GetTestPath("diffOrders.zip");
  const std::string dst = GetTempPath("diffOrders_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);

  int verified = verify(dst.c_str(), 4, false, true);
  ASSERT_EQ(0, verified);
}

TEST(Align, DirectoryEntryDoNotRequireAlignment) {
  const std::string src = GetTestPath("archiveWithOneDirectoryEntry.zip");
  int verified = verify(src.c_str(), 4, false, true);
  ASSERT_EQ(0, verified);
}

TEST(Align, DirectoryEntry) {
  const std::string src = GetTestPath("archiveWithOneDirectoryEntry.zip");
  const std::string dst = GetTempPath("archiveWithOneDirectoryEntry_out.zip");

  int processed = process(src.c_str(), dst.c_str(), 4, true, false, 4096);
  ASSERT_EQ(0, processed);
  ASSERT_EQ(true, sameContent(src, dst));

  int verified = verify(dst.c_str(), 4, false, true);
  ASSERT_EQ(0, verified);
}
