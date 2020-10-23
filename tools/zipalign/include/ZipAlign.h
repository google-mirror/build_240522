/*
 * Copyright (C) 2020 The Android Open Source Project
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

#ifndef ZIPALIGN_H
#define ZIPALIGN_H

namespace android {

/*
 * Process a file.  We open the input and output files, failing if the
 * output file exists and "force" wasn't specified.
 */
int process(const char* inFileName, const char* outFileName, int alignment, bool force, bool zopfli, bool pageAlignSharedLibs) ;

/*
 * Verify the alignment of a zip archive.
 */
int verify(const char* fileName, int alignment, bool verbose, bool pageAlignSharedLibs);

} // namespace android

#endif // ZIPALIGN_H
