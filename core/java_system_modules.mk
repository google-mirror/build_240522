###########################################################
## Builds a Java 9 system image, suitable for javac --system
##
###########################################################

LOCAL_MODULE_CLASS := JAVA_LIBRARIES
LOCAL_BUILT_MODULE_STEM := system-modules
#include $(BUILD_SYSTEM)/config.mk

JAVA_BASE_PACKAGES := $(strip \
    java.io \
    java.lang \
    java.lang.annotation \
    java.lang.invoke \
    java.lang.module \
    java.lang.ref \
    java.lang.reflect \
    java.math \
    java.net \
    java.net.spi \
    java.nio \
    java.nio.channels \
    java.nio.channels.spi \
    java.nio.charset \
    java.nio.charset.spi \
    java.nio.file \
    java.nio.file.attribute \
    java.nio.file.spi \
    java.security \
    java.security.acl \
    java.security.cert \
    java.security.interfaces \
    java.security.spec \
    java.text \
    java.text.spi \
    java.time \
    java.time.chrono \
    java.time.format \
    java.time.temporal \
    java.time.zone \
    java.util \
    java.util.concurrent \
    java.util.concurrent.atomic \
    java.util.concurrent.locks \
    java.util.function \
    java.util.jar \
    java.util.regex \
    java.util.spi \
    java.util.stream \
    java.util.zip \
    javax.crypto \
    javax.crypto.interfaces \
    javax.crypto.spec \
    javax.net \
    javax.net.ssl \
    javax.security.auth \
    javax.security.auth.callback \
    javax.security.auth.login \
    javax.security.auth.spi \
    javax.security.auth.x500 \
    javax.security.cert \
    com.sun.security.ntlm \
    jdk.internal.jimage \
    jdk.internal.jimage.decompressor \
    jdk.internal.loader \
    jdk.internal.jmod \
    jdk.internal.logger \
    jdk.internal.org.objectweb.asm \
    jdk.internal.org.objectweb.asm.tree \
    jdk.internal.org.objectweb.asm.util \
    jdk.internal.org.objectweb.asm.commons \
    jdk.internal.org.objectweb.asm.signature \
    jdk.internal.math \
    jdk.internal.misc \
    jdk.internal.module \
    jdk.internal.perf \
    jdk.internal.ref \
    jdk.internal.reflect \
    jdk.internal.vm \
    jdk.internal.vm.annotation \
    jdk.internal.util.jar \
    sun.net \
    sun.net.ext \
    sun.net.dns \
    sun.net.util \
    sun.net.www \
    sun.net.www. protocol.http \
    sun.nio.ch \
    sun.nio.cs \
    sun.nio.fs \
    sun.reflect.annotation \
    sun.reflect.generics.reflectiveObjects \
    sun.reflect.misc \
    sun.security.action \
    sun.security.internal.interfaces \
    sun.security.internal.spec \
    sun.security.jca \
    sun.security.pkcs \
    sun.security.provider \
    sun.security.provider.certpath \
    sun.security.rsa \
    sun.security.ssl \
    sun.security.timestamp \
    sun.security.tools \
    sun.security.util \
    sun.security.x509 \
    sun.security.validator \
    sun.text.resources \
    sun.util.cldr \
    sun.util.locale.provider \
    sun.util.logging \
    sun.util.resources \
)

LOCAL_CLASSES := $(call local-intermediates-dir,COMMON)/classes
LOCAL_CLASSPATH_ENTRIES := $(call java-lib-header-files,$(LOCAL_JAVA_LIBRARIES))
LOCAL_MODULE_DIR := $(call local-intermediates-dir,COMMON)/module
LOCAL_MODULE_INFO := $(LOCAL_MODULE_DIR)/module-info
LOCAL_MODULE_NAME := java.base
LOCAL_JMOD_DIR := $(call local-intermediates-dir,COMMON)/jmod
LOCAL_JMOD_FILE := $(LOCAL_JMOD_DIR)/$(LOCAL_MODULE_NAME).jmod
LOCAL_SYSTEM_MODULES := $(call local-intermediates-dir,COMMON)/system-modules

$(LOCAL_MODULE_INFO).class: PRIVATE_MODULE_DIR := $(LOCAL_MODULE_DIR)
$(LOCAL_MODULE_INFO).class: PRIVATE_MODULE_INFO := $(LOCAL_MODULE_INFO)
$(LOCAL_MODULE_INFO).class: PRIVATE_MODULE_NAME := $(LOCAL_MODULE_NAME)
$(LOCAL_MODULE_INFO).class: PRIVATE_CLASSPATH_ENTRIES := $(LOCAL_CLASSPATH_ENTRIES)
$(LOCAL_MODULE_INFO).class: $(LOCAL_CLASSPATH_ENTRIES)
	rm -rf $(PRIVATE_MODULE_DIR)
	mkdir -p $(PRIVATE_MODULE_DIR)
	echo "module $(PRIVATE_MODULE_NAME) {" > $(PRIVATE_MODULE_INFO).java
ifneq ($(LOCAL_MODULE),core-all-system-modules)
	build/tools/list-jar-packages.sh $(PRIVATE_CLASSPATH_ENTRIES) \
	      | grep -xF $(addprefix -e ,$(JAVA_BASE_PACKAGES)) \
	      | sed 's/^/    exports /' \
	      | sed 's/$$/;/' \
          >> $(PRIVATE_MODULE_INFO).java
else
	build/tools/list-jar-packages.sh $(PRIVATE_CLASSPATH_ENTRIES) \
	      | sed 's/^/    exports /' \
	      | sed 's/$$/;/' \
          >> $(PRIVATE_MODULE_INFO).java
endif
	echo "}" >> $(PRIVATE_MODULE_INFO).java
	$(JAVAC) --patch-module=java.base=$(call normalize-path-list,$(PRIVATE_CLASSPATH_ENTRIES)) $(PRIVATE_MODULE_INFO).java

$(LOCAL_CLASSES): PRIVATE_MODULE_INFO := $(LOCAL_MODULE_INFO)
$(LOCAL_CLASSES): PRIVATE_CLASSES := $(LOCAL_CLASSES)
$(LOCAL_CLASSES): PRIVATE_CLASSPATH_ENTRIES := $(LOCAL_CLASSPATH_ENTRIES)
$(LOCAL_CLASSES): $(LOCAL_MODULE_INFO).class $(LOCAL_CLASSPATH_ENTRIES)
	rm -rf $(PRIVATE_CLASSES)
	mkdir -p $(PRIVATE_CLASSES)
	$(call unzip-jar-files,$(PRIVATE_CLASSPATH_ENTRIES),$(PRIVATE_CLASSES))
ifneq ($(LOCAL_MODULE),core-all-system-modules)
	cd $(PRIVATE_CLASSES) ; find . -type f | $(ANDROID_BUILD_TOP)/build/tools/filter-classes-by-packages.py exclude $(JAVA_BASE_PACKAGES) | xargs rm -f ; find . -type d -empty -delete
endif
	cp $(PRIVATE_MODULE_INFO).class $(PRIVATE_CLASSES)/

$(LOCAL_JMOD_FILE): PRIVATE_JMOD_DIR := $(LOCAL_JMOD_DIR)
$(LOCAL_JMOD_FILE): PRIVATE_JMOD_FILE := $(LOCAL_JMOD_FILE)
$(LOCAL_JMOD_FILE): PRIVATE_CLASSES := $(LOCAL_CLASSES)
$(LOCAL_JMOD_FILE): $(LOCAL_CLASSPATH_ENTRIES) $(LOCAL_CLASSES)
	rm -rf $(PRIVATE_JMOD_DIR)
	mkdir -p $(PRIVATE_JMOD_DIR)
	$(JMOD) create \
	    --module-version 9 \
	    --target-platform $(HOST_OS)-$(HOST_ARCH) \
	    --class-path $(PRIVATE_CLASSES) \
	    --module-path $(PRIVATE_JMOD_DIR) $(PRIVATE_JMOD_FILE)

$(LOCAL_SYSTEM_MODULES): PRIVATE_SYSTEM_MODULES := $(LOCAL_SYSTEM_MODULES)
$(LOCAL_SYSTEM_MODULES): PRIVATE_JMOD_DIR := $(LOCAL_JMOD_DIR)
$(LOCAL_SYSTEM_MODULES): PRIVATE_MODULE_NAME := $(LOCAL_MODULE_NAME)
$(LOCAL_SYSTEM_MODULES): $(LOCAL_JMOD_FILE)
	rm -rf $(PRIVATE_SYSTEM_MODULES)
	$(JLINK) \
	    --module-path $(PRIVATE_JMOD_DIR) \
	    --add-modules $(PRIVATE_MODULE_NAME) \
	    --output $(PRIVATE_SYSTEM_MODULES)
	mkdir -p $(PRIVATE_SYSTEM_MODULES)/lib
	cp $(ANDROID_JAVA_HOME)/lib/jrt-fs.jar $(PRIVATE_SYSTEM_MODULES)/lib/

$(LOCAL_MODULE): $(LOCAL_SYSTEM_MODULES)

LOCAL_CLASSES:=
LOCAL_CLASSPATH_ENTRIES:=
LOCAL_MODULE_DIR:=
LOCAL_MODULE_INFO:=
LOCAL_MODULE_NAME:=
LOCAL_JMOD_DIR:=
LOCAL_JMOD_FILE:=
LOCAL_SYSTEM_MODULES:=
