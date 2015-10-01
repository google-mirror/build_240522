#include <limits.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
  if (argc != 2) {
    printf("usage: %s <path>\n", argv[0]);
    exit(EXIT_FAILURE);
  }

  char *real = realpath(argv[1], NULL);
  if (real == nullptr) {
    perror("realpath");
  }

  printf("%s\n", real);
  return 0;
}

