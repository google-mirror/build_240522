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
// General-purpose Zip archive access.  This class allows both reading and
// writing to Zip archives, including deletion of existing entries.
//
#ifndef __LIBS_ZIPFILE_H
#define __LIBS_ZIPFILE_H

#include <utils/Vector.h>
#include <utils/Errors.h>
#include <stdio.h>

#include "ZipEntry.h"

namespace android {

/*
 * Manipulate a Zip archive.
 *
 * Some changes will not be visible in the until until "flush" is called.
 *
 * The correct way to update a file archive is to make all changes to a
 * copy of the archive in a temporary file, and then unlink/rename over
 * the original after everything completes.  Because we're only interested
 * in using this for packaging, we don't worry about such things.  Crashing
 * after making changes and before flush() completes could leave us with
 * an unusable Zip archive.
 */
class ZipFile {
public:
    ZipFile(void)
      : mZipFp(NULL), mNeedCDRewrite(false)
      {}
    ~ZipFile(void) {
        flush();
        if (mZipFp != NULL)
            fclose(mZipFp);
        discardEntries();
    }

    /*
     * Open an archive.
     */
    status_t open(const char* zipFileName);

    /*
     * Set all timestamps to a static value
     */
    status_t removeTimestamps(void);

    /*
     * Flush changes.  If mNeedCDRewrite is set, this writes the central dir.
     */
    status_t flush(void);

private:
    /* these are private and not defined */
    ZipFile(const ZipFile& src);
    ZipFile& operator=(const ZipFile& src);

    class EndOfCentralDir {
    public:
        EndOfCentralDir(void) :
            mDiskNumber(0),
            mDiskWithCentralDir(0),
            mNumEntries(0),
            mTotalNumEntries(0),
            mCentralDirSize(0),
            mCentralDirOffset(0),
            mCommentLen(0),
            mComment(NULL)
            {}
        virtual ~EndOfCentralDir(void) {
            delete[] mComment;
        }

        status_t readBuf(const unsigned char* buf, int len);
        status_t write(FILE* fp);

        //unsigned long   mSignature;
        unsigned short  mDiskNumber;
        unsigned short  mDiskWithCentralDir;
        unsigned short  mNumEntries;
        unsigned short  mTotalNumEntries;
        unsigned long   mCentralDirSize;
        unsigned long   mCentralDirOffset;      // offset from first disk
        unsigned short  mCommentLen;
        unsigned char*  mComment;

        enum {
            kSignature      = 0x06054b50,
            kEOCDLen        = 22,       // EndOfCentralDir len, excl. comment

            kMaxCommentLen  = 65535,    // longest possible in ushort
            kMaxEOCDSearch  = kMaxCommentLen + EndOfCentralDir::kEOCDLen,

        };

        void dump(void) const;
    };


    /* read all entries in the central dir */
    status_t readCentralDir(void);

    /* clean up mEntries */
    void discardEntries(void);

    /*
     * We use stdio FILE*, which gives us buffering but makes dealing
     * with files >2GB awkward.  Until we support Zip64, we're fine.
     */
    FILE*           mZipFp;             // Zip file pointer

    /* one of these per file */
    EndOfCentralDir mEOCD;

    /* set this when we trash the central dir */
    bool            mNeedCDRewrite;

    /*
     * One ZipEntry per entry in the zip file.  I'm using pointers instead
     * of objects because it's easier than making operator= work for the
     * classes and sub-classes.
     */
    Vector<ZipEntry*>   mEntries;
};

}; // namespace android

#endif // __LIBS_ZIPFILE_H
