# Print a list of the modules that could be built

MODULE_INFO_JSON := $(PRODUCT_OUT)/module-info.json
MODULE_INFO_BUILD := $(PRODUCT_OUT)/module-info/BUILD

$(MODULE_INFO_JSON):
	@echo Generating $@
	$(hide) echo -ne '{\n ' > $@
	$(hide) echo -ne $(foreach m, $(sort $(ALL_MODULES)), \
		' "$(m)": {' \
			'"class": [$(foreach w,$(sort $(ALL_MODULES.$(m).CLASS)),"$(w)", )], ' \
			'"path": [$(foreach w,$(sort $(ALL_MODULES.$(m).PATH)),"$(w)", )], ' \
			'"tags": [$(foreach w,$(sort $(ALL_MODULES.$(m).TAGS)),"$(w)", )], ' \
			'"installed": [$(foreach w,$(sort $(ALL_MODULES.$(m).INSTALLED)),"$(w)", )], ' \
			'"compatibility_suites": [$(foreach w,$(sort $(ALL_MODULES.$(m).COMPATIBILITY_SUITES)),"$(w)", )], ' \
			'"auto_test_config": [$(ALL_MODULES.$(m).auto_test_config)], ' \
			'"module_name": "$(ALL_MODULES.$(m).MODULE_NAME)", ' \
			'"test_config": [$(foreach w,$(strip $(ALL_MODULES.$(m).TEST_CONFIG) $(ALL_MODULES.$(m).EXTRA_TEST_CONFIGS)),"$(w)", )], ' \
			'"dependencies": [$(foreach w,$(sort $(ALL_DEPS.$(m).ALL_DEPS)),"$(w)", )], ' \
			'"srcs": [$(foreach w,$(sort $(ALL_MODULES.$(m).SRCS)),"$(w)", )], ' \
			'"srcjars": [$(foreach w,$(sort $(ALL_MODULES.$(m).SRCJARS)),"$(w)", )], ' \
			'"classes_jar": [$(foreach w,$(sort $(ALL_MODULES.$(m).CLASSES_JAR)),"$(w)", )], ' \
			'"test_mainline_modules": [$(foreach w,$(sort $(ALL_MODULES.$(m).TEST_MAINLINE_MODULES)),"$(w)", )], ' \
			'"is_unit_test": "$(ALL_MODULES.$(m).IS_UNIT_TEST)", ' \
			'},\n' \
	 ) | sed -e 's/, *\]/]/g' -e 's/, *\}/ }/g' -e '$$s/,$$//' >> $@
	$(hide) echo '}' >> $@

# RECORD_ALL_DEPS=true SOONG_COLLECT_JAVA_DEPS=true SOONG_COLLECT_CC_DEPS=true m out/target/product/generic/module-info/BUILD
$(MODULE_INFO_BUILD):
	@echo Generating $@
	$(hide) echo -ne 'load("//build/make/core/tasks:module-info.bzl", "module_info")\n\n' > $@
	$(hide) echo -ne $(foreach m, $(sort $(ALL_MODULES)), \
		'module_info(\n' \
			'    name = "$(m)",\n' \
			'    module_name = "$(ALL_MODULES.$(m).MODULE_NAME)",\n' \
			'    installed = [$(foreach w,$(sort $(ALL_MODULES.$(m).INSTALLED)),"$(w)", )],\n' \
			'    module_class = [$(foreach w,$(sort $(ALL_MODULES.$(m).CLASS)),"$(w)", )],\n' \
			'    module_path = [$(foreach w,$(sort $(ALL_MODULES.$(m).PATH)),"$(w)", )],\n' \
			'    dependencies = [$(foreach w,$(sort $(ALL_DEPS.$(m).ALL_DEPS)),\n        "$(w)",)    ],\n' \
			')\n\n' \
	) | sed -e 's/^\s//g' >> $@

droidcore-unbundled: $(MODULE_INFO_JSON)

$(call dist-for-goals, general-tests, $(MODULE_INFO_JSON))
$(call dist-for-goals, droidcore-unbundled, $(MODULE_INFO_JSON))
