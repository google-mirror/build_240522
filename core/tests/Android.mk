ifeq ($(TEST_BUILD_SYSTEM),true)

# Can't use all-subdir-makefiles due to sanity check in $(call my-dir)
include $(call all-makefiles-under,build/core/tests)

endif
