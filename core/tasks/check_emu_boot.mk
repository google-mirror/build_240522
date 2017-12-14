EMULATOR_BOOT_TEST_SH := device/generic/goldfish/tools/emulator_boot_test.sh
emulator_current_dir := $(shell pwd)
emulator_out_dir := $(emulator_current_dir)/$(PRODUCT_OUT)
EMULATOR_BOOT_TEST_RESULT_FILE := $(DIST_DIR)/emulator-boot-test-result.txt
bootcompletestring := "emulator: INFO: boot completed"
bootfailstring := "emulator: ERROR: fail to boot after"
SUCCESS_FILE := $(DIST_DIR)/BOOT-SUCCESS.txt
FAIL_FILE := $(DIST_DIR)/BOOT-FAIL.txt
$(EMULATOR_BOOT_TEST_RESULT_FILE):
	@mkdir -p $(dir $@)
	$(hide) rm -f $@
	(export ANDROID_PRODUCT_OUT=$(emulator_out_dir);\
		export ANDROID_BUILD_TOP=$(emulator_current_dir);\
		$(EMULATOR_BOOT_TEST_SH) > $@)
.PHONY: check_emu_boot
$(EMULATOR_BOOT_TEST_RESULT_FILE): $(INSTALLED_QEMU_SYSTEMIMAGE) $(INSTALLED_QEMU_VENDORIMAGE)
check_emu_boot: $(EMULATOR_BOOT_TEST_RESULT_FILE)
	(if grep -q $(bootcompletestring) $(EMULATOR_BOOT_TEST_RESULT_FILE);\
	then echo boot_succeeded > $(SUCCESS_FILE); fi)
	(if grep -q $(bootfailstring) $(EMULATOR_BOOT_TEST_RESULT_FILE);\
	then echo boot_failed > $(FAIL_FILE); fi)
