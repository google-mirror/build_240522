#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

int main(int argc, char **argv) {
    if ((argc != 8)
        || strcmp(argv[1], "-a") || strcmp(argv[3], "--align-file-size")
        || access(argv[4], F_OK) || access(argv[5], F_OK)
        || access(argv[6], F_OK)) {
        fprintf(stderr, "Bad arguments\n");
        return 1;
    }
    int n = strlen(argv[6]) + strlen(argv[7]) + 50;
    char *buf = (char *)malloc(n);
    snprintf(buf, n, "cp %s %s", argv[6], argv[7]);
    return system(buf);
}
