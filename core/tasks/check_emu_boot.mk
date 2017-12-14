PRIVATE_EMULATOR_BOOT_TEST_SH := device/generic/goldfish/tools/emulator_boot_test.sh
PRIVATE_EMULATOR_BOOT_TEST_RESULT_FILE := $(DIST_DIR)/emulator-boot-test-result.txt
PRIVATE_BOOT_COMPLETE_STRING := "emulator: INFO: boot completed"
PRIVATE_BOOT_FAIL_STRING := "emulator: ERROR: fail to boot after"
PRIVATE_SUCCESS_FILE := $(DIST_DIR)/BOOT-SUCCESS.txt
PRIVATE_FAIL_FILE := $(DIST_DIR)/BOOT-FAIL.txt
.PHONY: check_emu_boot
check_emu_boot : $(PRIVATE_EMULATOR_BOOT_TEST_SH) droid
	@mkdir -p $(dir $(PRIVATE_EMULATOR_BOOT_TEST_RESULT_FILE))
	$(hide) rm -f $(PRIVATE_EMULATOR_BOOT_TEST_RESULT_FILE)
	$(hide) rm -f $(PRIVATE_SUCCESS_FILE)
	$(hide) rm -f $(PRIVATE_FAIL_FILE)
	(export ANDROID_PRODUCT_OUT=$$(cd $(PRODUCT_OUT);pwd);\
		export ANDROID_BUILD_TOP=$$(pwd);\
		$(PRIVATE_EMULATOR_BOOT_TEST_SH) > $(PRIVATE_EMULATOR_BOOT_TEST_RESULT_FILE))
	(if grep -q $(PRIVATE_BOOT_COMPLETE_STRING) $(PRIVATE_EMULATOR_BOOT_TEST_RESULT_FILE);\
	then echo boot_succeeded > $(PRIVATE_SUCCESS_FILE); fi)
	(if grep -q $(PRIVATE_BOO_TFAIL_STRING) $(PRIVATE_EMULATOR_BOOT_TEST_RESULT_FILE);\
	then echo boot_failed > $(PRIVATE_FAIL_FILE); fi)
