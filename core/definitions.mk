#
# Copyright (C) 2008 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

##
## Common build system definitions.  Mostly standard
## commands for building various types of targets, which
## are used by others to construct the final targets.
##

# These are variables we use to collect overall lists
# of things being processed.

# Full paths to all of the documentation
ALL_DOCS:=

# The short names of all of the targets in the system.
# For each element of ALL_MODULES, two other variables
# are defined:
#   $(ALL_MODULES.$(target)).BUILT
#   $(ALL_MODULES.$(target)).INSTALLED
# The BUILT variable contains LOCAL_BUILT_MODULE for that
# target, and the INSTALLED variable contains the LOCAL_INSTALLED_MODULE.
# Some targets may have multiple files listed in the BUILT and INSTALLED
# sub-variables.
ALL_MODULES:=

# Full paths to targets that should be added to the "make droid"
# set of installed targets.
ALL_DEFAULT_INSTALLED_MODULES:=

# The list of tags that have been defined by
# LOCAL_MODULE_TAGS.  Each word in this variable maps
# to a corresponding ALL_MODULE_TAGS.<tagname> variable
# that contains all of the INSTALLED_MODULEs with that tag.
ALL_MODULE_TAGS:=

# Similar to ALL_MODULE_TAGS, but contains the short names
# of all targets for a particular tag.  The top-level variable
# won't have the list of tags;  ust ALL_MODULE_TAGS to get
# the list of all known tags.  (This means that this variable
# will always be empty; it's just here as a placeholder for
# its sub-variables.)
ALL_MODULE_NAME_TAGS:=

# Full paths to all prebuilt files that will be copied
# (used to make the dependency on acp)
ALL_PREBUILT:=

# Full path to all files that are made by some tool
ALL_GENERATED_SOURCES:=

# Full path to all asm, C, C++, lex and yacc generated C files.
# These all have an order-only dependency on the copied headers
ALL_C_CPP_ETC_OBJECTS:=

# The list of dynamic binaries that haven't been stripped/compressed/prelinked.
ALL_ORIGINAL_DYNAMIC_BINARIES:=

# These files go into the SDK
ALL_SDK_FILES:=

# Files for dalvik.  This is often build without building the rest of the OS.
INTERNAL_DALVIK_MODULES:=

# All findbugs xml files
ALL_FINDBUGS_FILES:=

###########################################################
## Debugging; prints a variable list to stdout
###########################################################

# $(1): variable name list, not variable values
define print-vars
$(foreach var,$(1), \
  $(info $(var):) \
  $(foreach word,$($(var)), \
    $(info $(space)$(space)$(word)) \
   ) \
 )
endef

###########################################################
## Evaluates to true if the string contains the word true,
## and empty otherwise
## $(1): a var to test
###########################################################

define true-or-empty
$(filter true, $(1))
endef


###########################################################
## Retrieve the directory of the current makefile
###########################################################

# Figure out where we are.
define my-dir
$(strip \
  $(eval md_file_ := $$(lastword $$(MAKEFILE_LIST))) \
  $(if $(filter $(CLEAR_VARS),$(md_file_)), \
    $(error LOCAL_PATH must be set before including $$(CLEAR_VARS)) \
   , \
    $(patsubst %/,%,$(dir $(md_file_))) \
   ) \
 )
endef

###########################################################
## Retrieve a list of all makefiles immediately below some directory
###########################################################

define all-makefiles-under
$(wildcard $(1)/*/Android.mk)
endef

###########################################################
## Look under a directory for makefiles that don't have parent
## makefiles.
###########################################################

# $(1): directory to search under
# Ignores $(1)/Android.mk
define first-makefiles-under
$(shell build/tools/findleaves.py --prune=out --prune=.repo --prune=.git \
        --mindepth=2 $(1) Android.mk)
endef

###########################################################
## Retrieve a list of all makefiles immediately below your directory
###########################################################

define all-subdir-makefiles
$(call all-makefiles-under,$(call my-dir))
endef

###########################################################
## Look in the named list of directories for makefiles,
## relative to the current directory.
###########################################################

# $(1): List of directories to look for under this directory
define all-named-subdir-makefiles
$(wildcard $(addsuffix /Android.mk, $(addprefix $(call my-dir)/,$(1))))
endef

###########################################################
## Find all of the java files under the named directories.
## Meant to be used like:
##    SRC_FILES := $(call all-java-files-under,src tests)
###########################################################

define all-java-files-under
$(patsubst ./%,%, \
  $(shell cd $(LOCAL_PATH) ; \
          find $(1) -name "*.java" -and -not -name ".*") \
 )
endef

###########################################################
## Find all of the java files from here.  Meant to be used like:
##    SRC_FILES := $(call all-subdir-java-files)
###########################################################

define all-subdir-java-files
$(call all-java-files-under,.)
endef

###########################################################
## Find all of the c files under the named directories.
## Meant to be used like:
##    SRC_FILES := $(call all-c-files-under,src tests)
###########################################################

define all-c-files-under
$(patsubst ./%,%, \
  $(shell cd $(LOCAL_PATH) ; \
          find $(1) -name "*.c" -and -not -name ".*") \
 )
endef

###########################################################
## Find all of the c files from here.  Meant to be used like:
##    SRC_FILES := $(call all-subdir-c-files)
###########################################################

define all-subdir-c-files
$(call all-c-files-under,.)
endef

###########################################################
## Find all files named "I*.aidl" under the named directories,
## which must be relative to $(LOCAL_PATH).  The returned list
## is relative to $(LOCAL_PATH).
###########################################################

define all-Iaidl-files-under
$(patsubst ./%,%, \
  $(shell cd $(LOCAL_PATH) ; \
          find $(1) -name "I*.aidl" -and -not -name ".*") \
 )
endef

###########################################################
## Find all of the "I*.aidl" files under $(LOCAL_PATH).
###########################################################

define all-subdir-Iaidl-files
$(call all-Iaidl-files-under,.)
endef

###########################################################
## Find all of the logtags files under the named directories.
## Meant to be used like:
##    SRC_FILES := $(call all-logtags-files-under,src)
###########################################################

define all-logtags-files-under
$(patsubst ./%,%, \
  $(shell cd $(LOCAL_PATH) ; \
          find $(1) -name "*.logtags" -and -not -name ".*") \
  )
endef

###########################################################
## Find all of the .proto files under the named directories.
## Meant to be used like:
##    SRC_FILES := $(call all-proto-files-under,src)
###########################################################

define all-proto-files-under
$(patsubst ./%,%, \
  $(shell cd $(LOCAL_PATH) ; \
          find $(1) -name "*.proto" -and -not -name ".*") \
  )
endef

###########################################################
## Find all of the html files under the named directories.
## Meant to be used like:
##    SRC_FILES := $(call all-html-files-under,src tests)
###########################################################

define all-html-files-under
$(patsubst ./%,%, \
  $(shell cd $(LOCAL_PATH) ; \
          find $(1) -name "*.html" -and -not -name ".*") \
 )
endef

###########################################################
## Find all of the html files from here.  Meant to be used like:
##    SRC_FILES := $(call all-subdir-html-files)
###########################################################

define all-subdir-html-files
$(call all-html-files-under,.)
endef

###########################################################
## Find all of the files matching pattern
##    SRC_FILES := $(call find-subdir-files, <pattern>)
###########################################################

define find-subdir-files
$(patsubst ./%,%,$(shell cd $(LOCAL_PATH) ; find $(1)))
endef

###########################################################
# find the files in the subdirectory $1 of LOCAL_DIR
# matching pattern $2, filtering out files $3
# e.g.
#     SRC_FILES += $(call find-subdir-subdir-files, \
#                         css, *.cpp, DontWantThis.cpp)
###########################################################

define find-subdir-subdir-files
$(filter-out $(patsubst %,$(1)/%,$(3)),$(patsubst ./%,%,$(shell cd \
            $(LOCAL_PATH) ; find $(1) -maxdepth 1 -name $(2))))
endef

###########################################################
## Find all of the files matching pattern
##    SRC_FILES := $(call all-subdir-java-files)
###########################################################

define find-subdir-assets
$(if $(1),$(patsubst ./%,%, \
	$(shell if [ -d $(1) ] ; then cd $(1) ; find ./ -type f -and -not -type l ; fi)), \
	$(warning Empty argument supplied to find-subdir-assets) \
)
endef

###########################################################
## Find various file types in a list of directories relative to $(LOCAL_PATH)
###########################################################

define find-other-java-files
	$(call find-subdir-files,$(1) -name "*.java" -and -not -name ".*")
endef

define find-other-html-files
	$(call find-subdir-files,$(1) -name "*.html" -and -not -name ".*")
endef

###########################################################
## Scan through each directory of $(1) looking for files
## that match $(2) using $(wildcard).  Useful for seeing if
## a given directory or one of its parents contains
## a particular file.  Returns the first match found,
## starting furthest from the root.
###########################################################

define find-parent-file
$(strip \
  $(eval _fpf := $(wildcard $(strip $(1))/$(strip $(2)))) \
  $(if $(_fpf),$(_fpf), \
       $(if $(filter-out ./ .,$(1)), \
             $(call find-parent-file,$(patsubst %/,%,$(dir $(1))),$(2)) \
        ) \
   ) \
)
endef

###########################################################
## Function we can evaluate to introduce a dynamic dependency
###########################################################

define add-dependency
$(1): $(2)
endef

###########################################################
## Set up the dependencies for a prebuilt target
##  $(call add-prebuilt-file, srcfile, [targetclass])
###########################################################

define add-prebuilt-file
    $(eval $(include-prebuilt))
endef

define include-prebuilt
    include $$(CLEAR_VARS)
    LOCAL_SRC_FILES := $(1)
    LOCAL_BUILT_MODULE_STEM := $(1)
    LOCAL_MODULE_SUFFIX := $$(suffix $(1))
    LOCAL_MODULE := $$(basename $(1))
    LOCAL_MODULE_CLASS := $(2)
    include $$(BUILD_PREBUILT)
endef

###########################################################
## do multiple prebuilts
##  $(call target class, files ...)
###########################################################

define add-prebuilt-files
    $(foreach f,$(2),$(call add-prebuilt-file,$f,$(1)))
endef



###########################################################
## The intermediates directory.  Where object files go for
## a given target.  We could technically get away without
## the "_intermediates" suffix on the directory, but it's
## nice to be able to grep for that string to find out if
## anyone's abusing the system.
###########################################################

# $(1): target class, like "APPS"
# $(2): target name, like "NotePad"
# $(3): if non-empty, this is a HOST target.
# $(4): if non-empty, force the intermediates to be COMMON
define intermediates-dir-for
$(strip \
    $(eval _idfClass := $(strip $(1))) \
    $(if $(_idfClass),, \
        $(error $(LOCAL_PATH): Class not defined in call to intermediates-dir-for)) \
    $(eval _idfName := $(strip $(2))) \
    $(if $(_idfName),, \
        $(error $(LOCAL_PATH): Name not defined in call to intermediates-dir-for)) \
    $(eval _idfPrefix := $(if $(strip $(3)),HOST,TARGET)) \
    $(if $(filter $(_idfPrefix)-$(_idfClass),$(COMMON_MODULE_CLASSES))$(4), \
        $(eval _idfIntBase := $($(_idfPrefix)_OUT_COMMON_INTERMEDIATES)) \
      , \
        $(eval _idfIntBase := $($(_idfPrefix)_OUT_INTERMEDIATES)) \
     ) \
    $(_idfIntBase)/$(_idfClass)/$(_idfName)_intermediates \
)
endef

# Uses LOCAL_MODULE_CLASS, LOCAL_MODULE, and LOCAL_IS_HOST_MODULE
# to determine the intermediates directory.
#
# $(1): if non-empty, force the intermediates to be COMMON
define local-intermediates-dir
$(strip \
    $(if $(strip $(LOCAL_MODULE_CLASS)),, \
        $(error $(LOCAL_PATH): LOCAL_MODULE_CLASS not defined before call to local-intermediates-dir)) \
    $(if $(strip $(LOCAL_MODULE)),, \
        $(error $(LOCAL_PATH): LOCAL_MODULE not defined before call to local-intermediates-dir)) \
    $(call intermediates-dir-for,$(LOCAL_MODULE_CLASS),$(LOCAL_MODULE),$(LOCAL_IS_HOST_MODULE),$(1)) \
)
endef

###########################################################
## Convert "path/to/libXXX.so" to "-lXXX".
## Any "path/to/libXXX.a" elements pass through unchanged.
###########################################################

define normalize-libraries
$(foreach so,$(filter %.so,$(1)),-l$(patsubst lib%.so,%,$(notdir $(so))))\
$(filter-out %.so,$(1))
endef

# TODO: change users to call the common version.
define normalize-host-libraries
$(call normalize-libraries,$(1))
endef

define normalize-target-libraries
$(call normalize-libraries,$(1))
endef

###########################################################
## Convert a list of short module names (e.g., "framework", "Browser")
## into the list of files that are built for those modules.
## NOTE: this won't return reliable results until after all
## sub-makefiles have been included.
## $(1): target list
###########################################################

define module-built-files
$(foreach module,$(1),$(ALL_MODULES.$(module).BUILT))
endef

###########################################################
## Convert a list of short modules names (e.g., "framework", "Browser")
## into the list of files that are installed for those modules.
## NOTE: this won't return reliable results until after all
## sub-makefiles have been included.
## $(1): target list
###########################################################

define module-installed-files
$(foreach module,$(1),$(ALL_MODULES.$(module).INSTALLED))
endef

###########################################################
## Convert a list of short modules names (e.g., "framework", "Browser")
## into the list of files that should be used when linking
## against that module as a public API.
## TODO: Allow this for more than JAVA_LIBRARIES modules
## NOTE: this won't return reliable results until after all
## sub-makefiles have been included.
## $(1): target list
###########################################################

define module-stubs-files
$(foreach module,$(1),$(ALL_MODULES.$(module).STUBS))
endef

###########################################################
## Evaluates to the timestamp file for a doc module, which
## is the dependency that should be used.
## $(1): doc module
###########################################################

define doc-timestamp-for
$(OUT_DOCS)/$(strip $(1))-timestamp
endef


###########################################################
## Convert "core ext framework" to "out/.../javalib.jar ..."
## $(1): library list
## $(2): Non-empty if IS_HOST_MODULE
###########################################################

# $(1): library name
# $(2): Non-empty if IS_HOST_MODULE
define _java-lib-dir
$(call intermediates-dir-for, \
	JAVA_LIBRARIES,$(1),$(2),COMMON)
endef

# $(1): library name
# $(2): Non-empty if IS_HOST_MODULE
define _java-lib-full-classes.jar
$(call _java-lib-dir,$(1),$(2))/classes$(COMMON_JAVA_PACKAGE_SUFFIX)
endef

# $(1): library name list
# $(2): Non-empty if IS_HOST_MODULE
define java-lib-files
$(foreach lib,$(1),$(call _java-lib-full-classes.jar,$(lib),$(2)))
endef

# $(1): library name
# $(2): Non-empty if IS_HOST_MODULE
define _java-lib-full-dep
$(call _java-lib-dir,$(1),$(2))/javalib$(COMMON_JAVA_PACKAGE_SUFFIX)
endef

# $(1): library name list
# $(2): Non-empty if IS_HOST_MODULE
define java-lib-deps
$(foreach lib,$(1),$(call _java-lib-full-dep,$(lib),$(2)))
endef

###########################################################
## Convert "a b c" into "a:b:c"
###########################################################

empty :=
space := $(empty) $(empty)

define normalize-path-list
$(subst $(space),:,$(strip $(1)))
endef

###########################################################
## Read the word out of a colon-separated list of words.
## This has the same behavior as the built-in function
## $(word n,str).
##
## The individual words may not contain spaces.
##
## $(1): 1 based index
## $(2): value of the form a:b:c...
###########################################################

define word-colon
$(word $(1),$(subst :,$(space),$(2)))
endef

###########################################################
## Convert "a=b c= d e = f" into "a=b c=d e=f"
##
## $(1): list to collapse
## $(2): if set, separator word; usually "=", ":", or ":="
##       Defaults to "=" if not set.
###########################################################

define collapse-pairs
$(eval _cpSEP := $(strip $(if $(2),$(2),=)))\
$(subst $(space)$(_cpSEP)$(space),$(_cpSEP),$(strip \
    $(subst $(_cpSEP), $(_cpSEP) ,$(1))))
endef

###########################################################
## MODULE_TAG set operations
###########################################################

# Given a list of tags, return the targets that specify
# any of those tags.
# $(1): tag list
define modules-for-tag-list
$(sort $(foreach tag,$(1),$(ALL_MODULE_TAGS.$(tag))))
endef

# Same as modules-for-tag-list, but operates on
# ALL_MODULE_NAME_TAGS.
# $(1): tag list
define module-names-for-tag-list
$(sort $(foreach tag,$(1),$(ALL_MODULE_NAME_TAGS.$(tag))))
endef

# Given an accept and reject list, find the matching
# set of targets.  If a target has multiple tags and
# any of them are rejected, the target is rejected.
# Reject overrides accept.
# $(1): list of tags to accept
# $(2): list of tags to reject
#TODO(dbort): do $(if $(strip $(1)),$(1),$(ALL_MODULE_TAGS))
#TODO(jbq): as of 20100106 nobody uses the second parameter
define get-tagged-modules
$(filter-out \
	$(call modules-for-tag-list,$(2)), \
	    $(call modules-for-tag-list,$(1)))
endef

###########################################################
## Append a leaf to a base path.  Properly deals with
## base paths ending in /.
##
## $(1): base path
## $(2): leaf path
###########################################################

define append-path
$(subst //,/,$(1)/$(2))
endef


###########################################################
## Package filtering
###########################################################

# Given a list of installed modules (short or long names)
# return a list of the packages (yes, .apk packages, not
# modules in general) that are overridden by this list and,
# therefore, should not be installed.
# $(1): mixed list of installed modules
# TODO: This is fragile; find a reliable way to get this information.
define _get-package-overrides
 $(eval ### Discard any words containing slashes, unless they end in .apk, \
        ### in which case trim off the directory component and the suffix. \
        ### If there are no slashes, keep the entire word.)
 $(eval _gpo_names := $(subst /,@@@ @@@,$(1)))
 $(eval _gpo_names := \
     $(filter %.apk,$(_gpo_names)) \
     $(filter-out %@@@ @@@%,$(_gpo_names)))
 $(eval _gpo_names := $(patsubst %.apk,%,$(_gpo_names)))
 $(eval _gpo_names := $(patsubst @@@%,%,$(_gpo_names)))

 $(eval ### Remove any remaining words that contain dots.)
 $(eval _gpo_names := $(subst .,@@@ @@@,$(_gpo_names)))
 $(eval _gpo_names := $(filter-out %@@@ @@@%,$(_gpo_names)))

 $(eval ### Now we have a list of any words that could possibly refer to \
        ### packages, although there may be words that do not.  Only \
        ### real packages will be present under PACKAGES.*, though.)
 $(foreach _gpo_name,$(_gpo_names),$(PACKAGES.$(_gpo_name).OVERRIDES))
endef

define get-package-overrides
$(strip $(sort $(call _get-package-overrides,$(1))))
endef

###########################################################
## Output the command lines, or not
###########################################################

ifeq ($(strip $(SHOW_COMMANDS)),)
define pretty
@echo $1
endef
hide := @
else
define pretty
endef
hide :=
endif

###########################################################
## Dump the variables that are associated with targets
###########################################################

define dump-module-variables
@echo all_dependencies=$^
@echo PRIVATE_YACCFLAGS=$(PRIVATE_YACCFLAGS);
@echo PRIVATE_CFLAGS=$(PRIVATE_CFLAGS);
@echo PRIVATE_CPPFLAGS=$(PRIVATE_CPPFLAGS);
@echo PRIVATE_DEBUG_CFLAGS=$(PRIVATE_DEBUG_CFLAGS);
@echo PRIVATE_C_INCLUDES=$(PRIVATE_C_INCLUDES);
@echo PRIVATE_LDFLAGS=$(PRIVATE_LDFLAGS);
@echo PRIVATE_LDLIBS=$(PRIVATE_LDLIBS);
@echo PRIVATE_ARFLAGS=$(PRIVATE_ARFLAGS);
@echo PRIVATE_AAPT_FLAGS=$(PRIVATE_AAPT_FLAGS);
@echo PRIVATE_DX_FLAGS=$(PRIVATE_DX_FLAGS);
@echo PRIVATE_JAVACFLAGS=$(PRIVATE_JAVACFLAGS);
@echo PRIVATE_JAVA_LIBRARIES=$(PRIVATE_JAVA_LIBRARIES);
@echo PRIVATE_ALL_SHARED_LIBRARIES=$(PRIVATE_ALL_SHARED_LIBRARIES);
@echo PRIVATE_ALL_STATIC_LIBRARIES=$(PRIVATE_ALL_STATIC_LIBRARIES);
@echo PRIVATE_ALL_WHOLE_STATIC_LIBRARIES=$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES);
@echo PRIVATE_ALL_OBJECTS=$(PRIVATE_ALL_OBJECTS);
endef

###########################################################
## Commands for using sed to replace given variable values
###########################################################

define transform-variables
@mkdir -p $(dir $@)
@echo "Sed: $(if $(PRIVATE_MODULE),$(PRIVATE_MODULE),$@) <= $<"
$(hide) sed $(foreach var,$(REPLACE_VARS),-e "s/{{$(var)}}/$(subst /,\/,$(PWD)/$($(var)))/g") $< >$@
$(hide) if [ "$(suffix $@)" = ".sh" ]; then chmod a+rx $@; fi
endef


###########################################################
## Commands for munging the dependency files GCC generates
###########################################################

define transform-d-to-p
@cp $(@:%.o=%.d) $(@:%.o=%.P); \
	sed -e 's/#.*//' -e 's/^[^:]*: *//' -e 's/ *\\$$//' \
		-e '/^$$/ d' -e 's/$$/ :/' < $(@:%.o=%.d) >> $(@:%.o=%.P); \
	rm -f $(@:%.o=%.d)
endef

###########################################################
## Commands for running lex
###########################################################

define transform-l-to-cpp
@mkdir -p $(dir $@)
@echo "Lex: $(PRIVATE_MODULE) <= $<"
$(hide) $(LEX) -o$@ $<
endef

###########################################################
## Commands for running yacc
##
## Because the extension of c++ files can change, the
## extension must be specified in $1.
## E.g, "$(call transform-y-to-cpp,.cpp)"
###########################################################

define transform-y-to-cpp
@mkdir -p $(dir $@)
@echo "Yacc: $(PRIVATE_MODULE) <= $<"
$(YACC) $(PRIVATE_YACCFLAGS) -o $@ $<
touch $(@:$1=$(YACC_HEADER_SUFFIX))
echo '#ifndef '$(@F:$1=_h) > $(@:$1=.h)
echo '#define '$(@F:$1=_h) >> $(@:$1=.h)
cat $(@:$1=$(YACC_HEADER_SUFFIX)) >> $(@:$1=.h)
echo '#endif' >> $(@:$1=.h)
rm -f $(@:$1=$(YACC_HEADER_SUFFIX))
endef


###########################################################
## Commands for running aidl
###########################################################

define transform-aidl-to-java
@mkdir -p $(dir $@)
@echo "Aidl: $(PRIVATE_MODULE) <= $<"
$(hide) $(AIDL) -d$(patsubst %.java,%.P,$@) $(PRIVATE_AIDL_FLAGS) $< $@
endef
#$(AIDL) $(PRIVATE_AIDL_FLAGS) $< - | indent -nut -br -npcs -l1000 > $@


###########################################################
## Commands for running java-event-log-tags.py
###########################################################

define transform-logtags-to-java
@mkdir -p $(dir $@)
@echo "logtags: $@ <= $<"
$(hide) $(JAVATAGS) -o $@ $^
endef


###########################################################
## Commands for running protoc to compile .proto into .java
###########################################################

define transform-proto-to-java
@mkdir -p $(dir $@)
@echo "Protoc: $@ <= $(PRIVATE_PROTO_SRC_FILES)"
@rm -rf $(PRIVATE_PROTO_JAVA_OUTPUT_DIR)
@mkdir -p $(PRIVATE_PROTO_JAVA_OUTPUT_DIR)
$(hide) $(PROTOC) \
	$(addprefix --proto_path=, $(PRIVATE_PROTO_INCLUDES)) \
	$(PRIVATE_PROTO_JAVA_OUTPUT_OPTION)=$(PRIVATE_PROTO_JAVA_OUTPUT_DIR) \
	$(PRIVATE_PROTOC_FLAGS) \
	$(PRIVATE_PROTO_SRC_FILES)
$(hide) touch $@
endef

######################################################################
## Commands for running protoc to compile .proto into .pb.cc and .pb.h
######################################################################
define transform-proto-to-cc
@mkdir -p $(dir $@)
@echo "Protoc: $@ <= $<"
$(hide) $(PROTOC) \
	$(addprefix --proto_path=, $(PRIVATE_PROTO_INCLUDES)) \
	$(PRIVATE_PROTOC_FLAGS) \
	--cpp_out=$(PRIVATE_PROTO_CC_OUTPUT_DIR) $<
endef


###########################################################
## Commands for running gcc to compile a C++ file
###########################################################

define transform-cpp-to-o
@mkdir -p $(dir $@)
@echo "target $(PRIVATE_ARM_MODE) C++: $(PRIVATE_MODULE) <= $<"
$(hide) $(PRIVATE_CXX) \
	$(foreach incdir, \
	    $(PRIVATE_C_INCLUDES) \
	    $(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
		$(PRIVATE_TARGET_PROJECT_INCLUDES) \
		$(PRIVATE_TARGET_C_INCLUDES) \
	     ) \
	  , \
	    -I $(incdir) \
	 ) \
	-c \
	$(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
	    $(PRIVATE_TARGET_GLOBAL_CFLAGS) \
	    $(PRIVATE_TARGET_GLOBAL_CPPFLAGS) \
	    $(PRIVATE_ARM_CFLAGS) \
	 ) \
	-fno-rtti \
	$(PRIVATE_CFLAGS) \
	$(PRIVATE_CPPFLAGS) \
	$(PRIVATE_DEBUG_CFLAGS) \
	-MD -o $@ $<
$(hide) $(transform-d-to-p)
endef


###########################################################
## Commands for running gcc to compile a C file
###########################################################

# $(1): extra flags
define transform-c-or-s-to-o-no-deps
@mkdir -p $(dir $@)
$(hide) $(PRIVATE_CC) \
	$(foreach incdir, \
	    $(PRIVATE_C_INCLUDES) \
	    $(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
		$(PRIVATE_TARGET_PROJECT_INCLUDES) \
		$(PRIVATE_TARGET_C_INCLUDES) \
	     ) \
	  , \
	    -I $(incdir) \
	 ) \
	-c \
	$(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
	    $(PRIVATE_TARGET_GLOBAL_CFLAGS) \
	    $(PRIVATE_ARM_CFLAGS) \
	 ) \
	$(PRIVATE_CFLAGS) \
	$(1) \
	$(PRIVATE_DEBUG_CFLAGS) \
	-MD -o $@ $<
endef

define transform-c-to-o-no-deps
@echo "target $(PRIVATE_ARM_MODE) C: $(PRIVATE_MODULE) <= $<"
$(call transform-c-or-s-to-o-no-deps, )
endef

define transform-s-to-o-no-deps
@echo "target asm: $(PRIVATE_MODULE) <= $<"
$(call transform-c-or-s-to-o-no-deps, $(PRIVATE_ASFLAGS))
endef

define transform-c-to-o
$(transform-c-to-o-no-deps)
$(hide) $(transform-d-to-p)
endef

define transform-s-to-o
$(transform-s-to-o-no-deps)
$(hide) $(transform-d-to-p)
endef

###########################################################
## Commands for running gcc to compile an Objective-C file
## This should never happen for target builds but this
## will error at build time.
###########################################################

define transform-m-to-o-no-deps
@echo "target ObjC: $(PRIVATE_MODULE) <= $<"
$(call transform-c-or-s-to-o-no-deps)
endef

define transform-m-to-o
$(transform-m-to-o-no-deps)
$(hide) $(transform-d-to-p)
endef

###########################################################
## Commands for running gcc to compile a host C++ file
###########################################################

define transform-host-cpp-to-o
@mkdir -p $(dir $@)
@echo "host C++: $(PRIVATE_MODULE) <= $<"
$(hide) $(PRIVATE_CXX) \
	$(foreach incdir, \
	    $(PRIVATE_C_INCLUDES) \
	    $(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
		$(HOST_PROJECT_INCLUDES) \
		$(HOST_C_INCLUDES) \
	     ) \
	  , \
	    -I $(incdir) \
	 ) \
	-c \
	$(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
	    $(HOST_GLOBAL_CFLAGS) \
	    $(HOST_GLOBAL_CPPFLAGS) \
	 ) \
	$(PRIVATE_CFLAGS) \
	$(PRIVATE_CPPFLAGS) \
	$(PRIVATE_DEBUG_CFLAGS) \
	-MD -o $@ $<
$(transform-d-to-p)
endef


###########################################################
## Commands for running gcc to compile a host C file
###########################################################

# $(1): extra flags
define transform-host-c-or-s-to-o-no-deps
@mkdir -p $(dir $@)
$(hide) $(PRIVATE_CC) \
	$(foreach incdir, \
	    $(PRIVATE_C_INCLUDES) \
	    $(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
		$(HOST_PROJECT_INCLUDES) \
		$(HOST_C_INCLUDES) \
	     ) \
	  , \
	    -I $(incdir) \
	 ) \
	-c \
	$(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
	    $(HOST_GLOBAL_CFLAGS) \
	 ) \
	$(PRIVATE_CFLAGS) \
	$(1) \
	$(PRIVATE_DEBUG_CFLAGS) \
	-MD -o $@ $<
endef

define transform-host-c-to-o-no-deps
@echo "host C: $(PRIVATE_MODULE) <= $<"
$(call transform-host-c-or-s-to-o-no-deps, )
endef

define transform-host-s-to-o-no-deps
@echo "host asm: $(PRIVATE_MODULE) <= $<"
$(call transform-host-c-or-s-to-o-no-deps, $(PRIVATE_ASFLAGS))
endef

define transform-host-c-to-o
$(transform-host-c-to-o-no-deps)
$(transform-d-to-p)
endef

define transform-host-s-to-o
$(transform-host-s-to-o-no-deps)
$(transform-d-to-p)
endef

###########################################################
## Commands for running gcc to compile a host Objective-C file
###########################################################

define transform-host-m-to-o-no-deps
@echo "host ObjC: $(PRIVATE_MODULE) <= $<"
$(call transform-host-c-or-s-to-o-no-deps)
endef

define tranform-host-m-to-o
$(transform-host-m-to-o-no-deps)
$(transform-d-to-p)
endef

###########################################################
## Commands for running ar
###########################################################

define extract-and-include-target-whole-static-libs
$(foreach lib,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES), \
	$(hide) echo "preparing StaticLib: $(PRIVATE_MODULE) [including $(lib)]"; \
	ldir=$(PRIVATE_INTERMEDIATES_DIR)/WHOLE/$(basename $(notdir $(lib)))_objs;\
	rm -rf $$ldir; \
	mkdir -p $$ldir; \
	filelist=; \
	for f in `$(TARGET_AR) t $(lib)`; do \
	    $(TARGET_AR) p $(lib) $$f > $$ldir/$$f; \
	    filelist="$$filelist $$ldir/$$f"; \
	done ; \
	$(TARGET_AR) $(TARGET_GLOBAL_ARFLAGS) $(PRIVATE_ARFLAGS) $@ $$filelist;\
)
endef

# Explicitly delete the archive first so that ar doesn't
# try to add to an existing archive.
define transform-o-to-static-lib
@mkdir -p $(dir $@)
@rm -f $@
$(extract-and-include-target-whole-static-libs)
@echo "target StaticLib: $(PRIVATE_MODULE) ($@)"
$(hide) echo $(filter %.o, $^) | \
    xargs $(TARGET_AR) $(TARGET_GLOBAL_ARFLAGS) $(PRIVATE_ARFLAGS) $@
endef

###########################################################
## Commands for running host ar
###########################################################

define extract-and-include-host-whole-static-libs
$(foreach lib,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES), \
	@echo "preparing StaticLib: $(PRIVATE_MODULE) [including $(lib)]"; \
	ldir=$(PRIVATE_INTERMEDIATES_DIR)/WHOLE/$(basename $(notdir $(lib)))_objs;\
	rm -rf $$ldir; \
	mkdir -p $$ldir; \
	filelist=; \
	for f in `$(HOST_AR) t $(lib) | grep '\.o$$'`; do \
	    $(HOST_AR) p $(lib) $$f > $$ldir/$$f; \
	    filelist="$$filelist $$ldir/$$f"; \
	done ; \
	$(HOST_AR) $(HOST_GLOBAL_ARFLAGS) $(PRIVATE_ARFLAGS) $@ $$filelist;\
)
endef

# Explicitly delete the archive first so that ar doesn't
# try to add to an existing archive.
define transform-host-o-to-static-lib
@mkdir -p $(dir $@)
@rm -f $@
$(extract-and-include-host-whole-static-libs)
@echo "host StaticLib: $(PRIVATE_MODULE) ($@)"
echo $(filter %.o, $^) | \
	xargs $(HOST_AR) $(HOST_GLOBAL_ARFLAGS) $(PRIVATE_ARFLAGS) $@
endef


###########################################################
## Commands for running gcc to link a shared library or package
###########################################################

# ld just seems to be so finicky with command order that we allow
# it to be overriden en-masse see combo/linux-arm.make for an example.
ifneq ($(HOST_CUSTOM_LD_COMMAND),true)
define transform-host-o-to-shared-lib-inner
$(PRIVATE_CXX) \
	-Wl,-rpath-link=$(TARGET_OUT_INTERMEDIATE_LIBRARIES) \
	-Wl,-rpath,\$$ORIGIN/../lib \
	-shared -Wl,-soname,$(notdir $@) \
	$(PRIVATE_LDFLAGS) \
	$(HOST_GLOBAL_LD_DIRS) \
	$(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
	   $(HOST_GLOBAL_LDFLAGS) \
	) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-host-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(call normalize-host-libraries,$(PRIVATE_ALL_STATIC_LIBRARIES)) \
	$(call normalize-host-libraries,$(PRIVATE_ALL_SHARED_LIBRARIES)) \
	-o $@ \
	$(PRIVATE_LDLIBS)
endef
endif

define transform-host-o-to-shared-lib
@mkdir -p $(dir $@)
@echo "host SharedLib: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-host-o-to-shared-lib-inner)
endef

define transform-host-o-to-package
@mkdir -p $(dir $@)
@echo "host Package: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-host-o-to-shared-lib-inner)
endef


###########################################################
## Commands for running gcc to link a shared library or package
###########################################################

#echo >$@.vers "{"; \
#echo >>$@.vers " global:"; \
#$(BUILD_SYSTEM)/filter_symbols.sh $(TARGET_NM) "  " ";" $(filter %.o,$^) | sort -u >>$@.vers; \
#echo >>$@.vers " local:"; \
#echo >>$@.vers "  *;"; \
#echo >>$@.vers "};"; \

#	-Wl,--version-script=$@.vers \

# ld just seems to be so finicky with command order that we allow
# it to be overriden en-masse see combo/linux-arm.make for an example.
ifneq ($(TARGET_CUSTOM_LD_COMMAND),true)
define transform-o-to-shared-lib-inner
$(PRIVATE_CXX) \
	$(PRIVATE_TARGET_GLOBAL_LDFLAGS) \
	-Wl,-rpath-link=$(TARGET_OUT_INTERMEDIATE_LIBRARIES) \
	-Wl,-rpath,\$$ORIGIN/../lib \
	-shared -Wl,-soname,$(notdir $@) \
	$(PRIVATE_LDFLAGS) \
	$(PRIVATE_TARGET_GLOBAL_LD_DIRS) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-host-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(call normalize-target-libraries,$(PRIVATE_ALL_STATIC_LIBRARIES)) \
	$(call normalize-target-libraries,$(PRIVATE_ALL_SHARED_LIBRARIES)) \
	-o $@ \
	$(PRIVATE_LDLIBS)
endef
endif

define transform-o-to-shared-lib
@mkdir -p $(dir $@)
@echo "target SharedLib: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-o-to-shared-lib-inner)
endef

define transform-o-to-package
@mkdir -p $(dir $@)
@echo "target Package: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-o-to-shared-lib-inner)
endef


###########################################################
## Commands for filtering a target executable or library
###########################################################

define transform-to-stripped
@mkdir -p $(dir $@)
@echo "target Strip: $(PRIVATE_MODULE) ($@)"
$(hide) $(TARGET_STRIP_COMMAND)
endef

define transform-to-prelinked
@mkdir -p $(dir $@)
@echo "target Prelink: $(PRIVATE_MODULE) ($@)"
$(hide) $(APRIORI) \
		--prelinkmap $(TARGET_PRELINKER_MAP) \
		--locals-only \
		--quiet \
		$< \
		--output $@
endef


###########################################################
## Commands for running gcc to link an executable
###########################################################

ifneq ($(TARGET_CUSTOM_LD_COMMAND),true)
define transform-o-to-executable-inner
$(PRIVATE_CXX) \
	$(TARGET_GLOBAL_LDFLAGS) \
	-Wl,-rpath-link=$(TARGET_OUT_INTERMEDIATE_LIBRARIES) \
	$(TARGET_GLOBAL_LD_DIRS) \
	-Wl,-rpath-link=$(TARGET_OUT_INTERMEDIATE_LIBRARIES) \
	-Wl,-rpath,\$$ORIGIN/../lib \
	$(PRIVATE_LDFLAGS) \
	$(TARGET_GLOBAL_LD_DIRS) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-target-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(call normalize-target-libraries,$(PRIVATE_ALL_STATIC_LIBRARIES)) \
	$(call normalize-target-libraries,$(PRIVATE_ALL_SHARED_LIBRARIES)) \
	-o $@ \
	$(PRIVATE_LDLIBS)
endef
endif

define transform-o-to-executable
@mkdir -p $(dir $@)
@echo "target Executable: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-o-to-executable-inner)
endef


###########################################################
## Commands for running gcc to link a statically linked
## executable.  In practice, we only use this on arm, so
## the other platforms don't have the
## transform-o-to-static-executable defined
###########################################################

ifneq ($(TARGET_CUSTOM_LD_COMMAND),true)
define transform-o-to-static-executable-inner
endef
endif

define transform-o-to-static-executable
@mkdir -p $(dir $@)
@echo "target StaticExecutable: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-o-to-static-executable-inner)
endef


###########################################################
## Commands for running gcc to link a host executable
###########################################################

ifneq ($(HOST_CUSTOM_LD_COMMAND),true)
define transform-host-o-to-executable-inner
$(PRIVATE_CXX) \
	-Wl,-rpath-link=$(TARGET_OUT_INTERMEDIATE_LIBRARIES) \
	-Wl,-rpath,\$$ORIGIN/../lib \
	$(HOST_GLOBAL_LD_DIRS) \
	$(PRIVATE_LDFLAGS) \
	$(if $(PRIVATE_NO_DEFAULT_COMPILER_FLAGS),, \
		$(HOST_GLOBAL_LDFLAGS) \
	) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-host-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(call normalize-host-libraries,$(PRIVATE_ALL_STATIC_LIBRARIES)) \
	$(call normalize-host-libraries,$(PRIVATE_ALL_SHARED_LIBRARIES)) \
	-o $@ \
	$(PRIVATE_LDLIBS)
endef
endif

define transform-host-o-to-executable
@mkdir -p $(dir $@)
@echo "host Executable: $(PRIVATE_MODULE) ($@)"
$(hide) $(transform-host-o-to-executable-inner)
endef


###########################################################
## Commands for running javac to make .class files
###########################################################

#@echo "Source intermediates dir: $(PRIVATE_SOURCE_INTERMEDIATES_DIR)"
#@echo "Source intermediates: $$(find $(PRIVATE_SOURCE_INTERMEDIATES_DIR) -name '*.java')"

# TODO: Right now we generate the asset resources twice, first as part
# of generating the Java classes, then at the end when packaging the final
# assets.  This should be changed to do one of two things: (1) Don't generate
# any resource files the first time, only create classes during that stage;
# or (2) Don't use the -c flag with the second stage, instead taking the
# resource files from the first stage as additional input.  My original intent
# was to use approach (2), but this requires a little more work in the tool.
# Maybe we should just use approach (1).

# This rule creates the R.java and Manifest.java files, both of which
# are PRODUCT-neutral.  Don't pass PRODUCT_AAPT_CONFIG to this invocation.
define create-resource-java-files
@mkdir -p $(PRIVATE_SOURCE_INTERMEDIATES_DIR)
@mkdir -p $(dir $(PRIVATE_RESOURCE_PUBLICS_OUTPUT))
$(hide) $(AAPT) package $(PRIVATE_AAPT_FLAGS) -m \
    $(eval # PRODUCT_AAPT_CONFIG is intentionally missing-- see comment.) \
    $(addprefix -J , $(PRIVATE_SOURCE_INTERMEDIATES_DIR)) \
    $(addprefix -M , $(PRIVATE_ANDROID_MANIFEST)) \
    $(addprefix -P , $(PRIVATE_RESOURCE_PUBLICS_OUTPUT)) \
    $(addprefix -S , $(PRIVATE_RESOURCE_DIR)) \
    $(addprefix -A , $(PRIVATE_ASSET_DIR)) \
    $(addprefix -I , $(PRIVATE_AAPT_INCLUDES)) \
    $(addprefix -G , $(PRIVATE_PROGUARD_OPTIONS_FILE)) \
    $(addprefix --min-sdk-version , $(DEFAULT_APP_TARGET_SDK)) \
    $(addprefix --target-sdk-version , $(DEFAULT_APP_TARGET_SDK)) \
    $(if $(filter --version-code,$(PRIVATE_AAPT_FLAGS)),,$(addprefix --version-code , $(PLATFORM_SDK_VERSION))) \
    $(if $(filter --version-name,$(PRIVATE_AAPT_FLAGS)),,$(addprefix --version-name , $(PLATFORM_VERSION))) \
    $(addprefix --rename-manifest-package , $(PRIVATE_MANIFEST_PACKAGE_NAME)) \
    $(addprefix --rename-instrumentation-target-package , $(PRIVATE_MANIFEST_INSTRUMENTATION_FOR))
endef

ifeq ($(HOST_OS),windows)
xlint_unchecked :=
else
xlint_unchecked := -Xlint:unchecked
endif

# emit-line, <word list>, <output file>
define emit-line
   $(if $(1),echo -n '$(strip $(1)) ' >> $(2))
endef

# dump-words-to-file, <word list>, <output file>
define dump-words-to-file
        @rm -f $(2)
        @$(call emit-line,$(wordlist 1,200,$(1)),$(2))
        @$(call emit-line,$(wordlist 201,400,$(1)),$(2))
        @$(call emit-line,$(wordlist 401,600,$(1)),$(2))
        @$(call emit-line,$(wordlist 601,800,$(1)),$(2))
        @$(call emit-line,$(wordlist 801,1000,$(1)),$(2))
        @$(call emit-line,$(wordlist 1001,1200,$(1)),$(2))
        @$(call emit-line,$(wordlist 1201,1400,$(1)),$(2))
        @$(call emit-line,$(wordlist 1401,1600,$(1)),$(2))
        @$(call emit-line,$(wordlist 1601,1800,$(1)),$(2))
        @$(call emit-line,$(wordlist 1801,2000,$(1)),$(2))
        @$(call emit-line,$(wordlist 2001,2200,$(1)),$(2))
        @$(call emit-line,$(wordlist 2201,2400,$(1)),$(2))
        @$(call emit-line,$(wordlist 2401,2600,$(1)),$(2))
        @$(call emit-line,$(wordlist 2601,2800,$(1)),$(2))
        @$(call emit-line,$(wordlist 2801,3000,$(1)),$(2))
        @$(call emit-line,$(wordlist 3001,3200,$(1)),$(2))
        @$(call emit-line,$(wordlist 3201,3400,$(1)),$(2))
        @$(call emit-line,$(wordlist 3401,3600,$(1)),$(2))
        @$(call emit-line,$(wordlist 3601,3800,$(1)),$(2))
        @$(call emit-line,$(wordlist 3801,4000,$(1)),$(2))
        @$(if $(wordlist 4001,4002,$(1)),$(error Too many words ($(words $(1)))))
endef

# For a list of jar files, unzip them to a specified directory,
# but make sure that no META-INF files come along for the ride.
#
# $(1): files to unzip
# $(2): destination directory
define unzip-jar-files
  $(hide) for f in $(1); \
  do \
    if [ ! -f $$f ]; then \
      echo Missing file $$f; \
      exit 1; \
    fi; \
    unzip -qo $$f -d $(2); \
    (cd $(2) && rm -rf META-INF); \
  done
endef

# below we write the list of java files to java-source-list to avoid argument
# list length problems with Cygwin we filter out duplicate java file names
# because eclipse's compiler doesn't like them.
define transform-java-to-classes.jar
@echo "target Java: $(PRIVATE_MODULE) ($(PRIVATE_CLASS_INTERMEDIATES_DIR))"
$(hide) rm -f $@
$(hide) rm -rf $(PRIVATE_CLASS_INTERMEDIATES_DIR)
$(hide) mkdir -p $(PRIVATE_CLASS_INTERMEDIATES_DIR)
$(call unzip-jar-files,$(PRIVATE_STATIC_JAVA_LIBRARIES), \
    $(PRIVATE_CLASS_INTERMEDIATES_DIR))
$(call dump-words-to-file,$(PRIVATE_JAVA_SOURCES),$(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list)
$(hide) if [ -d "$(PRIVATE_SOURCE_INTERMEDIATES_DIR)" ]; then \
	    find $(PRIVATE_SOURCE_INTERMEDIATES_DIR) -name '*.java' >> $(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list; \
fi
$(hide) tr ' ' '\n' < $(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list \
    | sort -u > $(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list-uniq
$(hide) $(TARGET_JAVAC) -encoding ascii $(PRIVATE_BOOTCLASSPATH) \
    $(addprefix -classpath ,$(strip \
        $(call normalize-path-list,$(PRIVATE_ALL_JAVA_LIBRARIES)))) \
    $(PRIVATE_JAVACFLAGS) $(strip $(PRIVATE_JAVAC_DEBUG_FLAGS)) \
	$(if $(findstring true,$(LOCAL_WARNINGS_ENABLE)),$(xlint_unchecked),) \
    -extdirs "" -d $(PRIVATE_CLASS_INTERMEDIATES_DIR) \
    \@$(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list-uniq \
    || ( rm -rf $(PRIVATE_CLASS_INTERMEDIATES_DIR) ; exit 41 )
$(hide) rm -f $(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list
$(hide) rm -f $(dir $(PRIVATE_CLASS_INTERMEDIATES_DIR))/java-source-list-uniq
$(hide) mkdir -p $(dir $@)
$(hide) jar $(if $(strip $(PRIVATE_JAR_MANIFEST)),-cfm,-cf) \
    $@ $(PRIVATE_JAR_MANIFEST) -C $(PRIVATE_CLASS_INTERMEDIATES_DIR) .
$(hide) rm -rf $(PRIVATE_CLASS_INTERMEDIATES_DIR)
endef

define transform-classes.jar-to-emma
$(hide) java -classpath $(EMMA_JAR) emma instr -outmode fullcopy -outfile \
    $(PRIVATE_EMMA_COVERAGE_FILE) -ip $< -d $(PRIVATE_EMMA_INTERMEDIATES_DIR) \
    $(addprefix -ix , $(PRIVATE_EMMA_COVERAGE_FILTER))
endef

#TODO: use a smaller -Xmx value for most libraries;
#      only core.jar and framework.jar need a heap this big.
# Avoid the memory arguments on Windows, dx fails to load for some reason with them.
define transform-classes.jar-to-dex
@echo "target Dex: $(PRIVATE_MODULE)"
@mkdir -p $(dir $@)
$(hide) $(DX) \
    $(if $(findstring windows,$(HOST_OS)),,-JXms16M -JXmx1536M) \
    --dex --output=$@ \
    $(if $(NO_OPTIMIZE_DX), \
        --no-optimize) \
    $(if $(GENERATE_DEX_DEBUG), \
	    --debug --verbose \
	    --dump-to=$(@:.dex=.lst) \
	    --dump-width=1000) \
    $(PRIVATE_DX_FLAGS) \
    $<
endef

# Create a mostly-empty .jar file that we'll add to later.
# The MacOS jar tool doesn't like creating empty jar files,
# so we need to give it something.
define create-empty-package
@mkdir -p $(dir $@)
$(hide) touch $(dir $@)/dummy
$(hide) (cd $(dir $@) && jar cf $(notdir $@) dummy)
$(hide) zip -qd $@ dummy
$(hide) rm $(dir $@)/dummy
endef

#TODO: we kinda want to build different asset packages for
#      different configurations, then combine them later (or something).
#      Per-locale, etc.
#      A list of dynamic and static parameters;  build layers for
#      dynamic params that lay over the static ones.
#TODO: update the manifest to point to the package file
#Note that the version numbers are given to aapt as simple default
#values; applications can override these by explicitly stating
#them in their manifest.
define add-assets-to-package
$(hide) $(AAPT) package -u $(PRIVATE_AAPT_FLAGS) \
    $(addprefix -c , $(PRODUCT_AAPT_CONFIG)) \
    $(addprefix -M , $(PRIVATE_ANDROID_MANIFEST)) \
    $(addprefix -S , $(PRIVATE_RESOURCE_DIR)) \
    $(addprefix -A , $(PRIVATE_ASSET_DIR)) \
    $(addprefix -I , $(PRIVATE_AAPT_INCLUDES)) \
    $(addprefix --min-sdk-version , $(DEFAULT_APP_TARGET_SDK)) \
    $(addprefix --target-sdk-version , $(DEFAULT_APP_TARGET_SDK)) \
    $(addprefix --product , $(TARGET_AAPT_CHARACTERISTICS)) \
    $(if $(filter --version-code,$(PRIVATE_AAPT_FLAGS)),,$(addprefix --version-code , $(PLATFORM_SDK_VERSION))) \
    $(if $(filter --version-name,$(PRIVATE_AAPT_FLAGS)),,$(addprefix --version-name , $(PLATFORM_VERSION))) \
    $(addprefix --rename-manifest-package , $(PRIVATE_MANIFEST_PACKAGE_NAME)) \
    $(addprefix --rename-instrumentation-target-package , $(PRIVATE_MANIFEST_INSTRUMENTATION_FOR)) \
    -F $@
endef

define add-jni-shared-libs-to-package
$(hide) rm -rf $(dir $@)lib
$(hide) mkdir -p $(dir $@)lib/$(TARGET_CPU_ABI)
$(hide) cp $(PRIVATE_JNI_SHARED_LIBRARIES) $(dir $@)lib/$(TARGET_CPU_ABI)
$(hide) (cd $(dir $@) && zip -r $(notdir $@) lib)
$(hide) rm -rf $(dir $@)lib
endef

#TODO: update the manifest to point to the dex file
define add-dex-to-package
$(if $(filter classes.dex,$(notdir $(PRIVATE_DEX_FILE))),\
$(hide) $(AAPT) add -k $@ $(PRIVATE_DEX_FILE),\
$(eval _adtp_classes.dex := $(dir $(PRIVATE_DEX_FILE))/classes.dex)\
$(hide) cp $(PRIVATE_DEX_FILE) $(_adtp_classes.dex) && \
$(AAPT) add -k $@ $(_adtp_classes.dex) && \
rm -f $(_adtp_classes.dex))
endef

define add-java-resources-to-package
$(hide) jar uf $@ $(PRIVATE_EXTRA_JAR_ARGS)
endef

# Sign a package using the specified key/cert.
#
define sign-package
$(hide) mv $@ $@.unsigned
$(hide) java -jar $(SIGNAPK_JAR) \
	$(PRIVATE_CERTIFICATE) $(PRIVATE_PRIVATE_KEY) $@.unsigned $@.signed
$(hide) mv $@.signed $@
endef

# Align STORED entries of a package on 4-byte boundaries
# to make them easier to mmap.
#
define align-package
$(hide) mv $@ $@.unaligned
$(hide) $(ZIPALIGN) -f 4 $@.unaligned $@.aligned
$(hide) mv $@.aligned $@
endef

define install-dex-debug
$(hide) if [ -f "$(PRIVATE_INTERMEDIATES_DIR)/classes.dex" ]; then \
	    mkdir -p $(TOP)/dalvik/DEBUG-FILES; \
	    $(ACP) $(PRIVATE_INTERMEDIATES_DIR)/classes.dex \
		$(TOP)/dalvik/DEBUG-FILES/$(PRIVATE_MODULE).dex; \
	fi
$(hide) if [ -f "$(PRIVATE_INTERMEDIATES_DIR)/classes.lst" ]; then \
	    mkdir -p $(TOP)/dalvik/DEBUG-FILES; \
	    $(ACP) $(PRIVATE_INTERMEDIATES_DIR)/classes.lst \
		$(TOP)/dalvik/DEBUG-FILES/$(PRIVATE_MODULE).lst; \
	fi
endef

# TODO(joeo): If we can ever upgrade to post 3.81 make and get the
# new prebuilt rules to work, we should change this to copy the
# resources to the out directory and then copy the resources.

# Note: not using aapt tool for this because we aren't making
# an android package for the host.
define transform-host-java-to-package
@echo "host Java: $(PRIVATE_MODULE) ($(PRIVATE_CLASS_INTERMEDIATES_DIR))"
@rm -f $@
@rm -rf $(PRIVATE_CLASS_INTERMEDIATES_DIR)
@mkdir -p $(dir $@)
@mkdir -p $(PRIVATE_CLASS_INTERMEDIATES_DIR)
$(call unzip-jar-files,$(PRIVATE_STATIC_JAVA_LIBRARIES), \
    $(PRIVATE_CLASS_INTERMEDIATES_DIR))
$(call dump-words-to-file,$(sort\
	$(PRIVATE_JAVA_SOURCES)),\
	$(PRIVATE_INTERMEDIATES_DIR)/java-source-list-uniq)
$(hide) $(HOST_JAVAC) -encoding ascii -g \
	$(PRIVATE_JAVACFLAGS) $(xlint_unchecked) \
	$(addprefix -classpath ,$(strip \
		$(call normalize-path-list,$(PRIVATE_ALL_JAVA_LIBRARIES)))) \
	-extdirs "" -d $(PRIVATE_CLASS_INTERMEDIATES_DIR)\
        \@$(PRIVATE_INTERMEDIATES_DIR)/java-source-list-uniq || \
	( rm -rf $(PRIVATE_CLASS_INTERMEDIATES_DIR) ; exit 41 )
$(hide) rm -f $(PRIVATE_INTERMEDIATES_DIR)/java-source-list
$(hide) rm -f $(PRIVATE_INTERMEDIATES_DIR)/java-source-list-uniq
$(hide) jar $(if $(strip $(PRIVATE_JAR_MANIFEST)),-cfm,-cf) \
    $@ $(PRIVATE_JAR_MANIFEST) $(PRIVATE_EXTRA_JAR_ARGS) \
    -C $(PRIVATE_CLASS_INTERMEDIATES_DIR) .
endef

###########################################################
## Obfuscate a jar file
###########################################################

# PRIVATE_KEEP_FILE is a file containing a list of classes
# PRIVATE_INTERMEDIATES_DIR is a directory we can use for temporary files
# The module using this must depend on
#        $(HOST_OUT_JAVA_LIBRARIES)/proguard-4.0.1.jar
define obfuscate-jar
@echo "Obfuscate jar: $(notdir $@) ($@)"
@mkdir -p $(dir $@)
@rm -f $@
@mkdir -p $(PRIVATE_INTERMEDIATES_DIR)
$(hide) sed -e 's/^/-keep class /' < $(PRIVATE_KEEP_FILE) > \
		$(PRIVATE_INTERMEDIATES_DIR)/keep.pro
$(hide) java -Xmx512M -jar $(HOST_OUT_JAVA_LIBRARIES)/proguard-4.0.1.jar \
		-injars $< \
		-outjars $@ \
		-target 1.5 \
		-dontnote -dontwarn \
		-printmapping $(PRIVATE_INTERMEDIATES_DIR)/out.map \
		-forceprocessing \
		-renamesourcefileattribute SourceFile \
		-keepattributes Exceptions,InnerClasses,Signature,Deprecated,SourceFile,LineNumberTable,*Annotation*,EnclosingMethod \
		-repackageclasses \
		-keepclassmembers "class * { public protected *; }" \
		@$(PRIVATE_INTERMEDIATES_DIR)/keep.pro
endef

###########################################################
## Commands for copying files
###########################################################

# Define a rule to copy a header.  Used via $(eval) by copy_headers.make.
# $(1): source header
# $(2): destination header
define copy-one-header
$(2): $(1)
	@echo "Header: $$@"
	$$(copy-file-to-new-target-with-cp)
endef

# Define a rule to copy a file.  For use via $(eval).
# $(1): source file
# $(2): destination file
define copy-one-file
$(2): $(1) | $(ACP)
	@echo "Copy: $$@"
	$$(copy-file-to-target)
endef

# The -t option to acp and the -p option to cp is
# required for OSX.  OSX has a ridiculous restriction
# where it's an error for a .a file's modification time
# to disagree with an internal timestamp, and this
# macro is used to install .a files (among other things).

# Copy a single file from one place to another,
# preserving permissions and overwriting any existing
# file.
define copy-file-to-target
@mkdir -p $(dir $@)
$(hide) $(ACP) -fpt $< $@
endef

# The same as copy-file-to-target, but use the local
# cp command instead of acp.
define copy-file-to-target-with-cp
@mkdir -p $(dir $@)
$(hide) cp -fp $< $@
endef

# The same as copy-file-to-target, but use the zipalign tool to do so.
define copy-file-to-target-with-zipalign
@mkdir -p $(dir $@)
$(hide) $(ZIPALIGN) -f 4 $< $@
endef

# The same as copy-file-to-target, but strip out "# comment"-style
# comments (for config files and such).
define copy-file-to-target-strip-comments
@mkdir -p $(dir $@)
$(hide) sed -e 's/#.*$$//' -e 's/[ \t]*$$//' -e '/^$$/d' < $< > $@
endef

# The same as copy-file-to-target, but don't preserve
# the old modification time.
define copy-file-to-new-target
@mkdir -p $(dir $@)
$(hide) $(ACP) -fp $< $@
endef

# The same as copy-file-to-new-target, but use the local
# cp command instead of acp.
define copy-file-to-new-target-with-cp
@mkdir -p $(dir $@)
$(hide) cp -f $< $@
endef

# Copy a prebuilt file to a target location.
define transform-prebuilt-to-target
@echo "$(if $(PRIVATE_IS_HOST_MODULE),host,target) Prebuilt: $(PRIVATE_MODULE) ($@)"
$(copy-file-to-target)
endef

# Copy a prebuilt file to a target location, using zipalign on it.
define transform-prebuilt-to-target-with-zipalign
@echo "$(if $(PRIVATE_IS_HOST_MODULE),host,target) Prebuilt APK: $(PRIVATE_MODULE) ($@)"
$(copy-file-to-target-with-zipalign)
endef

# Copy a prebuilt file to a target location, stripping "# comment" comments.
define transform-prebuilt-to-target-strip-comments
@echo "$(if $(PRIVATE_IS_HOST_MODULE),host,target) Prebuilt: $(PRIVATE_MODULE) ($@)"
$(copy-file-to-target-strip-comments)
endef

###########################################################
## On some platforms (MacOS), after copying a static
## library, ranlib must be run to update an internal
## timestamp!?!?!
###########################################################

ifeq ($(HOST_RUN_RANLIB_AFTER_COPYING),true)
define transform-host-ranlib-copy-hack
    $(hide) ranlib $@ || true
endef
else
define transform-host-ranlib-copy-hack
true
endef
endif

ifeq ($(TARGET_RUN_RANLIB_AFTER_COPYING),true)
define transform-ranlib-copy-hack
    $(hide) ranlib $@
endef
else
define transform-ranlib-copy-hack
true
endef
endif


###########################################################
## Commands to call Proguard
###########################################################

# Command to copy the file with acp, if proguard is disabled.
define proguard-disabled-commands
@echo Copying: $@
$(hide) $(ACP) $< $@
endef

# Command to call Proguard
# $(1): extra flags for instrumentation.
define proguard-enabled-commands
@echo Proguard: $@
$(hide) $(PROGUARD) -injars $< -outjars $@ $(PRIVATE_PROGUARD_FLAGS) $(1)
endef

# Figure out the proguard dictionary file of the module that is instrumentationed for.
define get-instrumentation-proguard-flags
$(if $(PRIVATE_INSTRUMENTATION_FOR),$(if $(ALL_MODULES.$(PRIVATE_INSTRUMENTATION_FOR).PROGUARD_ENABLED),-applymapping $(call intermediates-dir-for,APPS,$(PRIVATE_INSTRUMENTATION_FOR),,COMMON)/proguard_dictionary))
endef

define transform-jar-to-proguard
$(eval _instrumentation_proguard_flags:=$(call get-instrumentation-proguard-flags))
$(eval _enable_proguard:=$(PRIVATE_PROGUARD_ENABLED)$(_instrumentation_proguard_flags))
$(if $(_enable_proguard),$(call proguard-enabled-commands,$(_instrumentation_proguard_flags)),$(call proguard-disabled-commands))
$(eval _instrumentation_proguard_flags:=)
$(eval _enable_proguard:=)
endef


###########################################################
## Stuff source generated from one-off tools
###########################################################

define transform-generated-source
@echo "target Generated: $(PRIVATE_MODULE) <= $<"
@mkdir -p $(dir $@)
$(hide) $(PRIVATE_CUSTOM_TOOL)
endef


###########################################################
## Assertions about attributes of the target
###########################################################

# $(1): The file to check
ifndef get-file-size
$(error HOST_OS must define get-file-size)
endif

# Convert a partition data size (eg, as reported in /proc/mtd) to the
# size of the image used to flash that partition (which includes a
# spare area for each page).
# $(1): the partition data size
define image-size-from-data-size
$(shell echo $$(($(1) / $(BOARD_NAND_PAGE_SIZE) * \
  ($(BOARD_NAND_PAGE_SIZE)+$(BOARD_NAND_SPARE_SIZE)))))
endef

# $(1): The file(s) to check (often $@)
# $(2): The maximum total image size, in decimal bytes
# $(3): the type of filesystem "yaffs" or "raw"
#
# If $(2) is empty, evaluates to "true"
#
# Reserve bad blocks.  Make sure that MAX(1% of partition size, 2 blocks)
# is left over after the image has been flashed.  Round the 1% up to the
# next whole flash block size.
define assert-max-file-size
$(if $(2), \
  size=$$(for i in $(1); do $(call get-file-size,$$i); echo +; done; echo 0); \
  total=$$(( $$( echo "$$size" ) )); \
  printname=$$(echo -n "$(1)" | tr " " +); \
  echo "$$printname total size is $$total"; \
  img_blocksize=$(call image-size-from-data-size,$(BOARD_FLASH_BLOCK_SIZE)); \
  if [ "$(3)" == "yaffs" ]; then \
    reservedblocks=8; \
  else \
    reservedblocks=0; \
  fi; \
  twoblocks=$$((img_blocksize * 2)); \
  onepct=$$((((($(2) / 100) - 1) / img_blocksize + 1) * img_blocksize)); \
  reserve=$$(((twoblocks > onepct ? twoblocks : onepct) + \
               reservedblocks * img_blocksize)); \
  maxsize=$$(($(2) - reserve)); \
  if [ "$$total" -gt "$$maxsize" ]; then \
    echo "error: $$printname too large ($$total > [$(2) - $$reserve])"; \
    false; \
  elif [ "$$total" -gt $$((maxsize - 32768)) ]; then \
    echo "WARNING: $$printname approaching size limit ($$total now; limit $$maxsize)"; \
  fi \
 , \
  true \
 )
endef

# Like assert-max-file-size, but the second argument is a partition
# size, which we'll convert to a max image size before checking it
# against the files.
#
# $(1): The file(s) to check (often $@)
# $(2): The partition size.
define assert-max-image-size
$(if $(2), \
  $(call assert-max-file-size,$(1),$(call image-size-from-data-size,$(2))), \
  true)
endef


###########################################################
## Define device-specific radio files
###########################################################

# Copy a radio image file to the output location, and add it to
# INSTALLED_RADIOIMAGE_TARGET.
# $(1): filename
define add-radio-file
  $(eval $(call add-radio-file-internal,$(1),$(notdir $(1))))
endef
define add-radio-file-internal
INSTALLED_RADIOIMAGE_TARGET += $$(PRODUCT_OUT)/$(2)
ALL_PREBUILT += $$(PRODUCT_OUT)/$(2)
$$(PRODUCT_OUT)/$(2) : $$(LOCAL_PATH)/$(1) | $$(ACP)
	$$(transform-prebuilt-to-target)
endef


###########################################################
# Override the package defined in $(1), setting the
# variables listed below differently.
#
#  $(1): The makefile to override (relative to the source
#        tree root)
#  $(2): Old LOCAL_PACKAGE_NAME value.
#  $(3): New LOCAL_PACKAGE_NAME value.
#  $(4): New LOCAL_MANIFEST_PACKAGE_NAME value.
#  $(5): New LOCAL_CERTIFICATE value.
#  $(6): New LOCAL_INSTRUMENTATION_FOR value.
#  $(7): New LOCAL_MANIFEST_INSTRUMENTATION_FOR value.
#
# Note that LOCAL_PACKAGE_OVERRIDES is NOT cleared in
# clear_vars.mk.
###########################################################
define inherit-package
  $(eval $(call inherit-package-internal,$(1),$(2),$(3),$(4),$(5),$(6),$(7)))
endef

define inherit-package-internal
  LOCAL_PACKAGE_OVERRIDES \
      := $(strip $(1))||$(strip $(2))||$(strip $(3))||$(strip $(4))||&&$(strip $(5))||&&$(strip $(6))||&&$(strip $(7)) $(LOCAL_PACKAGE_OVERRIDES)
  include $(1)
  LOCAL_PACKAGE_OVERRIDES \
      := $(wordlist 1,$(words $(LOCAL_PACKAGE_OVERRIDES)), $(LOCAL_PACKAGE_OVERRIDES))
endef

# To be used with inherit-package above
# Evalutes to true if the package was overridden
define set-inherited-package-variables
$(strip $(call set-inherited-package-variables-internal))
endef

define keep-or-override
$(eval $(1) := $(if $(2),$(2),$($(1))))
endef

define set-inherited-package-variables-internal
  $(eval _o := $(subst ||, ,$(lastword $(LOCAL_PACKAGE_OVERRIDES))))
  $(eval _n := $(subst ||, ,$(firstword $(LOCAL_PACKAGE_OVERRIDES))))
  $(if $(filter $(word 2,$(_n)),$(LOCAL_PACKAGE_NAME)), \
    $(eval LOCAL_PACKAGE_NAME := $(word 3,$(_o))) \
    $(eval LOCAL_MANIFEST_PACKAGE_NAME := $(word 4,$(_o))) \
    $(call keep-or-override,LOCAL_CERTIFICATE,$(patsubst &&%,%,$(word 5,$(_o)))) \
    $(call keep-or-override,LOCAL_INSTRUMENTATION_FOR,$(patsubst &&%,%,$(word 6,$(_o)))) \
    $(call keep-or-override,LOCAL_MANIFEST_INSTRUMENTATION_FOR,$(patsubst &&%,%,$(word 7,$(_o)))) \
    $(eval LOCAL_OVERRIDES_PACKAGES := $(sort $(LOCAL_OVERRIDES_PACKAGES) $(word 2,$(_o)))) \
    true \
  ,)
endef


###########################################################
## Other includes
###########################################################

# -----------------------------------------------------------------
# Rules and functions to help copy important files to DIST_DIR
# when requested.
include $(BUILD_SYSTEM)/distdir.mk

# -----------------------------------------------------------------
# The modules allowed to use a user tag
include $(BUILD_SYSTEM)/user_tags.mk

# broken:
#	$(foreach file,$^,$(if $(findstring,.a,$(suffix $file)),-l$(file),$(file)))

###########################################################
## Misc notes
###########################################################

#DEPDIR = .deps
#df = $(DEPDIR)/$(*F)

#SRCS = foo.c bar.c ...

#%.o : %.c
#	@$(MAKEDEPEND); \
#	  cp $(df).d $(df).P; \
#	  sed -e 's/#.*//' -e 's/^[^:]*: *//' -e 's/ *\\$$//' \
#	      -e '/^$$/ d' -e 's/$$/ :/' < $(df).d >> $(df).P; \
#	  rm -f $(df).d
#	$(COMPILE.c) -o $@ $<

#-include $(SRCS:%.c=$(DEPDIR)/%.P)


#%.o : %.c
#	$(COMPILE.c) -MD -o $@ $<
#	@cp $*.d $*.P; \
#	  sed -e 's/#.*//' -e 's/^[^:]*: *//' -e 's/ *\\$$//' \
#	      -e '/^$$/ d' -e 's/$$/ :/' < $*.d >> $*.P; \
#	  rm -f $*.d
