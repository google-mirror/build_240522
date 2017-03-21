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

#include <stdbool.h>
#define _GNU_SOURCE /* for asprintf */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include <private/android_filesystem_config.h>

/*
 * This program expects android_device_dirs and android_device_files
 * to be defined in the supplied android_filesystem_config.h file in
 * the device/<vendor>/<product> $(TARGET_DEVICE_DIR). Then generates
 * the binary format used in the /system/etc/fs_config_dirs and
 * the /system/etc/fs_config_files to be used by the runtimes.
 */
#include "android_filesystem_config.h"

#ifdef NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS
  static const struct fs_path_config android_device_dirs[] = {
};
#endif

#ifdef NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_FILES
static const struct fs_path_config android_device_files[] = {
#ifdef NO_ANDROID_FILESYSTEM_CONFIG_DEVICE_DIRS
    { 0, AID_ROOT, AID_ROOT, 0, "system/etc/fs_config_dirs" },
    { 0, AID_ROOT, AID_ROOT, 0, "vendor/etc/fs_config_dirs" },
    { 0, AID_ROOT, AID_ROOT, 0, "oem/etc/fs_config_dirs" },
#endif
    { 0, AID_ROOT, AID_ROOT, 0, "system/etc/fs_config_files" },
    { 0, AID_ROOT, AID_ROOT, 0, "vendor/etc/fs_config_files" },
    { 0, AID_ROOT, AID_ROOT, 0, "oem/etc/fs_config_files" },
};
#endif

static void usage() {
  fprintf(stderr,
    "Generate binary content for fs_config_dirs (-D) and fs_config_files (-F)\n"
    "from device-specific android_filesystem_config.h override.\n"
    "Split (-S) the content targetted to system, vendor or oem partitions,\n"
    "adding a numerical priority prefix, and partition suffix added to the\n"
    "specified out-file.\n\n"
    "Usage: fs_config_generate -D|-F [[-S] -o output-file]\n");
}

#ifndef ARRAY_SIZE  /* popular macro */
#define ARRAY_SIZE(x) (sizeof(x) / sizeof((x)[0]))
#endif

int main(int argc, char** argv) {
  const struct fs_path_config* pc;
  const struct fs_path_config* end;
  bool dir = false, file = false, split = false;
  static const char* partition[] = {  /* In priority order < 10 items */
    "system", "vendor", "oem"
  };
  FILE* fp[ARRAY_SIZE(partition)];
  size_t plen[ARRAY_SIZE(partition)];
  size_t idx;
  int opt;

  for (idx = 0; idx < ARRAY_SIZE(fp); ++idx) {
    fp[idx] = stdout;
    plen[idx] = strlen(partition[idx]);
  }
  while((opt = getopt(argc, argv, "DFSho:")) != -1) {
    switch(opt) {
    case 'D':
      if (file) {
        fprintf(stderr, "Must specify only -D or -F\n");
        usage();
        exit(EXIT_FAILURE);
      }
      dir = true;
      break;
    case 'F':
      if (dir) {
        fprintf(stderr, "Must specify only -F or -D\n");
        usage();
        exit(EXIT_FAILURE);
      }
      file = true;
      break;
    case 'S':
      if (!split) {
        for (idx = 0; idx < ARRAY_SIZE(fp); ++idx) {
          if (fp[idx] != stdout) {
            fprintf(stderr,
                    "Must specify -S before specifying the output basename\n");
            usage();
            exit(EXIT_FAILURE);
          }
        }
      }
      split = true;
      break;
    case 'o':
      for (idx = 0; idx < ARRAY_SIZE(fp); ++idx) {
        if (fp[idx] != stdout) {
          fprintf(stderr, "Specify only one output %s\n",
                  split ? "basename" : "file");
          usage();
          exit(EXIT_FAILURE);
        }
      }
      if (!split) {
        fp[0] = fopen(optarg, "wb");
        if (fp[0] == NULL) {
          fprintf(stderr, "Can not open \"%s\"\n", optarg);
          exit(EXIT_FAILURE);
        }
        for (idx = 1; idx < ARRAY_SIZE(fp); ++idx) {
          fp[idx] = fp[0];
        }
        break;
      }
      for (idx = 0; idx < ARRAY_SIZE(fp); ++idx) {
        char* name = NULL;
        char* base = strrchr(optarg, '/');
        if (base) {
          ++base;
        } else {
          base = optarg;
        }
        /* Adjust format if idx can exceed a value of 9 */
        asprintf(&name, "%.*s%zu-%s.%s", (int)(base - optarg), optarg, idx, base, partition[idx]);
        if (!name) {
          fprintf(stderr, "Can not allocate for %.*s%zu-%s.%s\n", (int)(base - optarg), optarg, idx,
                  base, partition[idx]);
          exit(EXIT_FAILURE);
        }
        fp[idx] = fopen(name, "wb");
        if (!fp[idx]) {
          fprintf(stderr, "Can not open \"%s\"\n", name);
          free(name);
          exit(EXIT_FAILURE);
        }
        free(name);
      }
      break;
    case 'h':
      usage();
      exit(EXIT_SUCCESS);
    default:
      usage();
      exit(EXIT_FAILURE);
    }
  }

  if (!file && !dir) {
    fprintf(stderr, "Must specify either -F or -D\n");
    usage();
    exit(EXIT_FAILURE);
  }

  if (dir) {
    pc = android_device_dirs;
    end = &android_device_dirs[ARRAY_SIZE(android_device_dirs)];
  } else {
    pc = android_device_files;
    end = &android_device_files[ARRAY_SIZE(android_device_files)];
  }
  for(; (pc < end) && pc->prefix; pc++) {
    char buffer[512];
    ssize_t len = fs_config_generate(buffer, sizeof(buffer), pc);
    if (len < 0) {
      fprintf(stderr, "Entry too large\n");
      exit(EXIT_FAILURE);
    }
    for (idx = 1; idx < ARRAY_SIZE(fp); ++idx) {
      /* does the file name reference <prefix>/ ? */
      if (!strncmp(pc->prefix, partition[idx], plen[idx]) &&
          (pc->prefix[plen[idx]] == '/')) {
        break;
      }
      /* does the file name reference system/<prefix>/ ? */
      if (!strncmp(pc->prefix, partition[0], plen[0]) &&
          (pc->prefix[plen[0]] == '/') &&
          !strncmp(pc->prefix + plen[0] + 1, partition[idx], plen[idx]) &&
          (pc->prefix[plen[0] + 1 + plen[idx]] == '/')) {
        break;
      }
    }
    /* if not found, then default into first (system) priority group */
    if (idx >= ARRAY_SIZE(fp)) idx = 0;
    if (fwrite(buffer, 1, len, fp[idx]) != (size_t)len) {
      fprintf(stderr, "Write failure %s\n", partition[idx]);
      exit(EXIT_FAILURE);
    }
  }
  for (idx = 0; idx < ARRAY_SIZE(fp); ++idx) {
    size_t sub;
    for (sub = idx + 1; sub < ARRAY_SIZE(fp); ++sub) {
      if (fp[idx] == fp[sub]) fp[sub] = NULL;
    }
    if (fp[idx]) fclose(fp[idx]);
  }

  return 0;
}
