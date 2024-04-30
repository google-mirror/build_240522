# APS causes errors
# LOCAL_MODULE_CLASS := FAKE
# build/make/core/base_rules.mk:266: error: cts/hostsidetests/devicepolicy: unhandled install path "HOST_OUT_APPS for CtsDevicePolicyManagerTestCases_LockSettings_NoFlakes"

# LOCAL_MODULE_CLASS := ETC

include $(BUILD_SYSTEM)/base_rules.mk

# tryue LOCAL_MODULE_STE, LOCAL_BUILT_MODULE_STEM, and LOCAL_INSTALLED_MODULE_STEM
# Base wants LOCAL_BUILT_MODULE for something like this: out/target/product/vsoc_x86_64/obj/APPS/FrameworksServicesTests_intermediates/package.apk
# It also gets added to the device.$$ME_all_targets
