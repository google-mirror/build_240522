# Print a list of the modules that could be built

MODULE_INFO_JSON := $(PRODUCT_OUT)/module-info.json

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
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
			'"module_name": "$(ALL_MODULES.$(m).MODULE_NAME)", ' \
			'"test_config": [$(if $(ALL_MODULES.$(m).TEST_CONFIG),"$(ALL_MODULES.$(m).TEST_CONFIG)")], ' \
			'"dependencies": [$(foreach w,$(sort $(ALL_DEPS.$(m).ALL_DEPS)),"$(w)", )], ' \
			'"srcs": [$(foreach w,$(sort $(ALL_MODULES.$(m).SRCS)),"$(w)", )], ' \
			'"srcjars": [$(foreach w,$(sort $(ALL_MODULES.$(m).SRCJARS)),"$(w)", )], ' \
			'"classes_jar": [$(foreach w,$(sort $(ALL_MODULES.$(m).CLASSES_JAR)),"$(w)", )], ' \
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
			'},\n' \
	 ) | sed -e 's/, *\]/]/g' -e 's/, *\}/ }/g' -e '$$s/,$$//' >> $@
	$(hide) echo '}' >> $@


<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# If ONE_SHOT_MAKEFILE is set, our view of the world is smaller, so don't
# rewrite the file in that came.
ifndef ONE_SHOT_MAKEFILE
files: $(MODULE_INFO_JSON)
endif
=======
droidcore: $(MODULE_INFO_JSON)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

