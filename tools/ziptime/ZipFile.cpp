/*
 * Copyright (C) 2006 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//
// Access to Zip archives.
//

#define LOG_TAG "zip"

#include <utils/Log.h>

#include "ZipFile.h"

#include <memory.h>
#include <sys/stat.h>
#include <errno.h>
#include <assert.h>

using namespace android;

/*
 * Some environments require the "b", some choke on it.
 */
#define FILE_OPEN_RO        "rb"
#define FILE_OPEN_RW        "r+b"
#define FILE_OPEN_RW_CREATE "w+b"

/* should live somewhere else? */
static status_t errnoToStatus(int err)
{
    if (err == ENOENT)
        return NAME_NOT_FOUND;
    else if (err == EACCES)
        return PERMISSION_DENIED;
    else
        return UNKNOWN_ERROR;
}

/*
 * Open a file and parse its guts.
 */
status_t ZipFile::open(const char* zipFileName)
{
    assert(mZipFp == NULL);     // no reopen

    /* open the file */
    mZipFp = fopen(zipFileName, FILE_OPEN_RW);
    if (mZipFp == NULL) {
        int err = errno;
        ALOGD("fopen failed: %d\n", err);
        return errnoToStatus(err);
    }

    /*
     * Load the central directory.  If that fails, then this probably
     * isn't a Zip archive.
     */
    return readCentralDir();
}

/*
 * Find the central directory and read the contents.
 *
 * The fun thing about ZIP archives is that they may or may not be
 * readable from start to end.  In some cases, notably for archives
 * that were written to stdout, the only length information is in the
 * central directory at the end of the file.
 *
 * Of course, the central directory can be followed by a variable-length
 * comment field, so we have to scan through it backwards.  The comment
 * is at most 64K, plus we have 18 bytes for the end-of-central-dir stuff
 * itself, plus apparently sometimes people throw random junk on the end
 * just for the fun of it.
 *
 * This is all a little wobbly.  If the wrong value ends up in the EOCD
 * area, we're hosed.  This appears to be the way that everbody handles
 * it though, so we're in pretty good company if this fails.
 */
status_t ZipFile::readCentralDir(void)
{
    status_t result = NO_ERROR;
    unsigned char* buf = NULL;
    off_t fileLength, seekStart;
    long readAmount;
    int i;

    fseek(mZipFp, 0, SEEK_END);
    fileLength = ftell(mZipFp);
    rewind(mZipFp);

    /* too small to be a ZIP archive? */
    if (fileLength < EndOfCentralDir::kEOCDLen) {
        ALOGD("Length is %ld -- too small\n", (long)fileLength);
        result = INVALID_OPERATION;
        goto bail;
    }

    buf = new unsigned char[EndOfCentralDir::kMaxEOCDSearch];
    if (buf == NULL) {
        ALOGD("Failure allocating %d bytes for EOCD search",
             EndOfCentralDir::kMaxEOCDSearch);
        result = NO_MEMORY;
        goto bail;
    }

    if (fileLength > EndOfCentralDir::kMaxEOCDSearch) {
        seekStart = fileLength - EndOfCentralDir::kMaxEOCDSearch;
        readAmount = EndOfCentralDir::kMaxEOCDSearch;
    } else {
        seekStart = 0;
        readAmount = (long) fileLength;
    }
    if (fseek(mZipFp, seekStart, SEEK_SET) != 0) {
        ALOGD("Failure seeking to end of zip at %ld", (long) seekStart);
        result = UNKNOWN_ERROR;
        goto bail;
    }

    /* read the last part of the file into the buffer */
    if (fread(buf, 1, readAmount, mZipFp) != (size_t) readAmount) {
        ALOGD("short file? wanted %ld\n", readAmount);
        result = UNKNOWN_ERROR;
        goto bail;
    }

    /* find the end-of-central-dir magic */
    for (i = readAmount - 4; i >= 0; i--) {
        if (buf[i] == 0x50 &&
            ZipEntry::getLongLE(&buf[i]) == EndOfCentralDir::kSignature)
        {
            ALOGV("+++ Found EOCD at buf+%d\n", i);
            break;
        }
    }
    if (i < 0) {
        ALOGD("EOCD not found, not Zip\n");
        result = INVALID_OPERATION;
        goto bail;
    }

    /* extract eocd values */
    result = mEOCD.readBuf(buf + i, readAmount - i);
    if (result != NO_ERROR) {
        ALOGD("Failure reading %ld bytes of EOCD values", readAmount - i);
        goto bail;
    }
    //mEOCD.dump();

    if (mEOCD.mDiskNumber != 0 || mEOCD.mDiskWithCentralDir != 0 ||
        mEOCD.mNumEntries != mEOCD.mTotalNumEntries)
    {
        ALOGD("Archive spanning not supported\n");
        result = INVALID_OPERATION;
        goto bail;
    }

    /*
     * So far so good.  "mCentralDirSize" is the size in bytes of the
     * central directory, so we can just seek back that far to find it.
     * We can also seek forward mCentralDirOffset bytes from the
     * start of the file.
     *
     * We're not guaranteed to have the rest of the central dir in the
     * buffer, nor are we guaranteed that the central dir will have any
     * sort of convenient size.  We need to skip to the start of it and
     * read the header, then the other goodies.
     *
     * The only thing we really need right now is the file comment, which
     * we're hoping to preserve.
     */
    if (fseek(mZipFp, mEOCD.mCentralDirOffset, SEEK_SET) != 0) {
        ALOGD("Failure seeking to central dir offset %ld\n",
             mEOCD.mCentralDirOffset);
        result = UNKNOWN_ERROR;
        goto bail;
    }

    /*
     * Loop through and read the central dir entries.
     */
    ALOGV("Scanning %d entries...\n", mEOCD.mTotalNumEntries);
    int entry;
    for (entry = 0; entry < mEOCD.mTotalNumEntries; entry++) {
        ZipEntry* pEntry = new ZipEntry;

        result = pEntry->initFromCDE(mZipFp);
        if (result != NO_ERROR) {
            ALOGD("initFromCDE failed\n");
            delete pEntry;
            goto bail;
        }

        mEntries.add(pEntry);
    }


    /*
     * If all went well, we should now be back at the EOCD.
     */
    {
        unsigned char checkBuf[4];
        if (fread(checkBuf, 1, 4, mZipFp) != 4) {
            ALOGD("EOCD check read failed\n");
            result = INVALID_OPERATION;
            goto bail;
        }
        if (ZipEntry::getLongLE(checkBuf) != EndOfCentralDir::kSignature) {
            ALOGD("EOCD read check failed\n");
            result = UNKNOWN_ERROR;
            goto bail;
        }
        ALOGV("+++ EOCD read check passed\n");
    }

bail:
    delete[] buf;
    return result;
}

/*
 * Empty the mEntries vector.
 */
void ZipFile::discardEntries(void)
{
    int count = mEntries.size();

    while (--count >= 0)
        delete mEntries[count];

    mEntries.clear();
}

/*
 * Set all timestamps to static values, write out the LFHs
 */
status_t ZipFile::removeTimestamps(void)
{
    for (auto entry : mEntries) {
        entry->removeTimestamps();

        if (fseek(mZipFp, entry->getLFHOffset(), SEEK_SET) != 0) {
            return UNKNOWN_ERROR;
        }
        entry->mLFH.write(mZipFp);
    }

    mNeedCDRewrite = true;

    return NO_ERROR;
}

/*
 * Flush any pending writes.
 *
 * In particular, this will crunch out deleted entries, and write the
 * Central Directory and EOCD if we have stomped on them.
 */
status_t ZipFile::flush(void)
{
    status_t result = NO_ERROR;
    long eocdPosn;
    int i, count;

    if (!mNeedCDRewrite)
        return NO_ERROR;

    assert(mZipFp != NULL);

    if (fseek(mZipFp, mEOCD.mCentralDirOffset, SEEK_SET) != 0)
        return UNKNOWN_ERROR;

    count = mEntries.size();
    for (i = 0; i < count; i++) {
        ZipEntry* pEntry = mEntries[i];
        pEntry->mCDE.write(mZipFp);
    }

    eocdPosn = ftell(mZipFp);
    mEOCD.mCentralDirSize = eocdPosn - mEOCD.mCentralDirOffset;

    mEOCD.write(mZipFp);

    /* should we clear the "newly added" flag in all entries now? */

    mNeedCDRewrite = false;
    return NO_ERROR;
}


/*
 * ===========================================================================
 *      ZipFile::EndOfCentralDir
 * ===========================================================================
 */

/*
 * Read the end-of-central-dir fields.
 *
 * "buf" should be positioned at the EOCD signature, and should contain
 * the entire EOCD area including the comment.
 */
status_t ZipFile::EndOfCentralDir::readBuf(const unsigned char* buf, int len)
{
    /* don't allow re-use */
    assert(mComment == NULL);

    if (len < kEOCDLen) {
        /* looks like ZIP file got truncated */
        ALOGD(" Zip EOCD: expected >= %d bytes, found %d\n",
            kEOCDLen, len);
        return INVALID_OPERATION;
    }

    /* this should probably be an assert() */
    if (ZipEntry::getLongLE(&buf[0x00]) != kSignature)
        return UNKNOWN_ERROR;

    mDiskNumber = ZipEntry::getShortLE(&buf[0x04]);
    mDiskWithCentralDir = ZipEntry::getShortLE(&buf[0x06]);
    mNumEntries = ZipEntry::getShortLE(&buf[0x08]);
    mTotalNumEntries = ZipEntry::getShortLE(&buf[0x0a]);
    mCentralDirSize = ZipEntry::getLongLE(&buf[0x0c]);
    mCentralDirOffset = ZipEntry::getLongLE(&buf[0x10]);
    mCommentLen = ZipEntry::getShortLE(&buf[0x14]);

    // TODO: validate mCentralDirOffset

    if (mCommentLen > 0) {
        if (kEOCDLen + mCommentLen > len) {
            ALOGD("EOCD(%d) + comment(%d) exceeds len (%d)\n",
                kEOCDLen, mCommentLen, len);
            return UNKNOWN_ERROR;
        }
        mComment = new unsigned char[mCommentLen];
        memcpy(mComment, buf + kEOCDLen, mCommentLen);
    }

    return NO_ERROR;
}

/*
 * Write an end-of-central-directory section.
 */
status_t ZipFile::EndOfCentralDir::write(FILE* fp)
{
    unsigned char buf[kEOCDLen];

    ZipEntry::putLongLE(&buf[0x00], kSignature);
    ZipEntry::putShortLE(&buf[0x04], mDiskNumber);
    ZipEntry::putShortLE(&buf[0x06], mDiskWithCentralDir);
    ZipEntry::putShortLE(&buf[0x08], mNumEntries);
    ZipEntry::putShortLE(&buf[0x0a], mTotalNumEntries);
    ZipEntry::putLongLE(&buf[0x0c], mCentralDirSize);
    ZipEntry::putLongLE(&buf[0x10], mCentralDirOffset);
    ZipEntry::putShortLE(&buf[0x14], mCommentLen);

    if (fwrite(buf, 1, kEOCDLen, fp) != kEOCDLen)
        return UNKNOWN_ERROR;
    if (mCommentLen > 0) {
        assert(mComment != NULL);
        if (fwrite(mComment, mCommentLen, 1, fp) != mCommentLen)
            return UNKNOWN_ERROR;
    }

    return NO_ERROR;
}

/*
 * Dump the contents of an EndOfCentralDir object.
 */
void ZipFile::EndOfCentralDir::dump(void) const
{
    ALOGD(" EndOfCentralDir contents:\n");
    ALOGD("  diskNum=%u diskWCD=%u numEnt=%u totalNumEnt=%u\n",
        mDiskNumber, mDiskWithCentralDir, mNumEntries, mTotalNumEntries);
    ALOGD("  centDirSize=%lu centDirOff=%lu commentLen=%u\n",
        mCentralDirSize, mCentralDirOffset, mCommentLen);
}

