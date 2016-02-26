# Rules to generate a zip file that included test modules
# based on the configuration.
#

.PHONY: brillo_tests

my_modules := $(shell sed -n -e 's/\(^[^#].*\),.*/\1/p' device/generic/brillo/tests.txt)

my_package_name := brillo_tests

include $(BUILD_SYSTEM)/tasks/tools/package-modules.mk

.PHONY: brillo_tests
brillo_tests : $(my_package_zip)

name := $(TARGET_PRODUCT)-brillo-tests-$(FILE_NAME_TAG)
$(call dist-for-goals, brillo_tests, $(my_package_zip):$(name).zip)

