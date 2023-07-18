#include <stdio.h>

#include "com_android_aconfig_test.h"

int
main() {
    printf("hello flag %d\n", com_android_aconfig_test_disabled_rw());
}
