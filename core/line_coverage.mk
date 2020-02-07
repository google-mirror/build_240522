# -----------------------------------------------------------------
# Make target for line coverage. This target generates a zip file
# called `line_coverage.zip` that contains a large set of zip files
# one for each fuzz target/critical component. Each zip file contains
# a set of profile files (*.gcno) that we will use to generate line
# coverage reports. Furthermore, target compiles all fuzz targets
# with line coverage instrumentation enabled and packs them into
# another zip file called `line_coverage_fuzz_targets.zip`.
# For more details: go/android-line-coverage-doc
#
# To run the make target set the coverage related envvars first:
# 	NATIVE_LINE_COVERAGE=true NATIVE_COVERAGE=true \
#	COVERAGE_PATHS=* make haiku-line-coverage
# -----------------------------------------------------------------

# TODO(b/148306195): Due this issue some fuzz targets cannot be built with
# line coverage instrumentation. For now we just blacklist them.
blacklisted_fuzz_targets := libneuralnetworks_fuzzer

fuzz_targets := $(ALL_FUZZ_TARGETS)
fuzz_targets := $(filter-out $(blacklisted_fuzz_targets),$(fuzz_targets))

# Critical components are defined in: go/haiku-most-wanted
# Please note that adding/Removing critical components is very rare.
critical_components := \
	libaudioprocessing \
	libbinder \
	libbluetooth_gd \
	libbrillo \
	lib-bt-packets \
	libbt-stack \
	libcameraservice \
	libcurl \
	libffi \
	libhardware \
	libhevcdec \
	libhevcenc \
	libinputflinger \
	libmpeg2dec \
	libopus \
	libosi \
	libpdx \
	libselinux \
	libstagefright \
	libunwind \
	libvixl \
	libvold \
	libyuv

# We want the line coverage for all fuzz targets + critical components.
all_components := $(critical_components) $(fuzz_targets)

# This make target generates a zip file with all profile files.
.PHONY: haiku-line-coverage
haiku-line-coverage: $(SOONG_ZIP) $(all_components)
	$(SOONG_ZIP) -o $(PRODUCT_OUT)/line_coverage.zip -D $(PRODUCT_OUT)/coverage
	$(SOONG_ZIP) -o $(PRODUCT_OUT)/line_coverage_fuzz_targets.zip -D $(PRODUCT_OUT)/data/fuzz
$(call dist-for-goals,haiku-line-coverage,$(PRODUCT_OUT)/line_coverage.zip $(PRODUCT_OUT)/line_coverage_fuzz_targets.zip)
