/*
 * Copyright (C) 2015 The Android Open Source Project
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

/*
 * Zip tool to remove dynamic timestamps
 */
#include <libazip/ZipFile.h>

#include <stdlib.h>
#include <stdio.h>

using namespace android;

/* Jan 01 2008 */
#define STATIC_DATE (28 << 9 | 1 << 5 | 1)
#define STATIC_TIME 0

static void usage(void)
{
    fprintf(stderr, "Zip timestamp utility\n");
    fprintf(stderr, "Copyright (C) 2015 The Android Open Source Project\n\n");
    fprintf(stderr, "Usage: ziptime file.zip\n");
}

int main(int argc, char* const argv[])
{
    if (argc != 2) {
        usage();
        return 2;
    }

    ZipFile zip;
    status_t status;
    status = zip.open(argv[1], ZipFile::kOpenReadWrite);
    if (status != NO_ERROR) {
        fprintf(stderr, "Unable to open zip archive '%s': %d\n",
                argv[1], (int) status);
        return 1;
    }

    status = zip.resetTimestamps(STATIC_TIME, STATIC_DATE);
    if (status != NO_ERROR) {
        fprintf(stderr, "Unable to zero timestamps in '%s': %d\n",
                argv[1], (int) status);
        return 1;
    }

    status = zip.flush();
    if (status != NO_ERROR) {
        fprintf(stderr, "unable to finalize zip '%s': %d\n",
                argv[1], (int) status);
        return 1;
    }

    return 0;
}
