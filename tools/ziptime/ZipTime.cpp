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
#include "ZipFile.h"

#include <stdlib.h>
#include <stdio.h>

using namespace android;

/*
 * Show program usage.
 */
void usage(void)
{
    fprintf(stderr, "Zip timestamp utility\n");
    fprintf(stderr, "Copyright (C) 2015 The Android Open Source Project\n\n");
    fprintf(stderr,
        "Usage: ziptime file.zip\n" );
}

/*
 * Process a file.
 */
static int process(const char* fileName)
{
    ZipFile zip;

    if (zip.open(fileName) != NO_ERROR) {
        fprintf(stderr, "Unable to open '%s' as zip archive\n", fileName);
        return 1;
    }

    if (zip.removeTimestamps() != NO_ERROR) {
        fprintf(stderr, "Failed to set timestamps\n");
        return 1;
    }

    if (zip.flush() != NO_ERROR) {
        fprintf(stderr, "Failed to write zipfile\n");
        return 1;
    }

    return 0;
}

/*
 * Parse args.
 */
int main(int argc, char* const argv[])
{
    bool wantUsage = false;
    int result = 1;
    char* endp;

    if (argc < 2) {
        wantUsage = true;
        goto bail;
    }

    argc--;
    argv++;

    while (argc && argv[0][0] == '-') {
        const char* cp = argv[0] +1;

        while (*cp != '\0') {
            switch (*cp) {
            default:
                fprintf(stderr, "ERROR: unknown flag -%c\n", *cp);
                wantUsage = true;
                goto bail;
            }

            cp++;
        }

        argc--;
        argv++;
    }

    if (argc != 1) {
        wantUsage = true;
        goto bail;
    }

    result = process(argv[0]);

bail:
    if (wantUsage) {
        usage();
        result = 2;
    }

    return result;
}
