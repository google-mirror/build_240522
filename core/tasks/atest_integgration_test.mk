.PHONY: gen_atest_command

ATEST_INTEGRATION_RESULT := $(PRODUCT_OUT)/atest_cmd_result.txt

$(ATEST_INTEGRATION_RESULT): $(MODULE_INFO_JSON)
	@echo Generating $@
	/usr/local/google/home/yangbill/aosp/master/tools/asuite/atest/atest_integration_tests.py -cy
	$(hide) touch $@

gen_atest_command : $(ATEST_INTEGRATION_RESULT)
