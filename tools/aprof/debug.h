/*
 * Copyright (C) 2011 The Android Open Source Project
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 2.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software Foundation,
 * Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 */

#ifndef DEBUG_H
#define DEBUG_H

#include <stdlib.h>
#include <stdio.h>
#include <common.h>

#ifdef DEBUG

    #define FAILIF(cond, msg...) do {                        \
	if (unlikely(cond)) {                                \
        fprintf(stderr, "%s(%d): ", __FILE__, __LINE__); \
		fprintf(stderr, ##msg);                          \
		exit(1);                                         \
	}                                                    \
} while(0)

/* Debug enabled */
    #define ASSERT(x) do {                                \
	if (unlikely(!(x))) {                             \
		fprintf(stderr,                               \
				"ASSERTION FAILURE %s:%d: [%s]\n",    \
				__FILE__, __LINE__, #x);              \
		exit(1);                                      \
	}                                                 \
} while(0)

#else

    #define FAILIF(cond, msg...) do { \
	if (unlikely(cond)) {         \
		fprintf(stderr, ##msg);   \
		exit(1);                  \
	}                             \
} while(0)

/* No debug */
    #define ASSERT(x)   do { } while(0)

#endif/* DEBUG */

#define FAILIF_LIBELF(cond, function) \
    FAILIF(cond, "%s(): %s\n", #function, elf_errmsg(elf_errno()));

static inline void *MALLOC(unsigned int size) {
    void *m = malloc(size);
    FAILIF(NULL == m, "malloc(%d) failed!\n", size);
    return m;
}

static inline void *CALLOC(unsigned int num_entries, unsigned int entry_size) {
    void *m = calloc(num_entries, entry_size);
    FAILIF(NULL == m, "calloc(%d, %d) failed!\n", num_entries, entry_size);
    return m;
}

static inline void *REALLOC(void *ptr, unsigned int size) {
    void *m = realloc(ptr, size);
    FAILIF(NULL == m, "realloc(%p, %d) failed!\n", ptr, size);
    return m;
}

static inline void FREE(void *ptr) {
    free(ptr);
}

static inline void FREEIF(void *ptr) {
    if (ptr) FREE(ptr);
}

#define PRINT(x...)  do {                             \
    extern int quiet_flag;                            \
    if(likely(!quiet_flag))                           \
        fprintf(stdout, ##x);                         \
} while(0)

#define ERROR PRINT

#define INFO(x...)  do {                              \
    extern int verbose;                          \
    if(unlikely(verbose))                        \
        fprintf(stdout, ##x);                         \
} while(0)

#define WARING(x...)                                  \
    fprintf(stderr, ##x);

/* Prints a hex and ASCII dump of the selected buffer to the selected stream. */
int dump_hex_buffer(FILE *s, void *b, size_t l, size_t elsize);

#endif/*DEBUG_H*/
