#################################################
## Perform configuration steps for SANITIZE_HOST.
#################################################
ifeq ($(SANITIZE_HOST),true)
ifneq ($(strip $(LOCAL_CLANG)),false)
ifneq ($(strip $(LOCAL_ADDRESS_SANITIZER)),false)
    LOCAL_ADDRESS_SANITIZER := true
endif
endif
endif
