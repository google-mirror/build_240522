
ifeq (true,$(non_system_module))
  ifneq (,$(filter $(dir $(DEFAULT_SYSTEM_DEV_CERTIFICATE))%,$(LOCAL_CERTIFICATE)))
    CERTIFICATE_VIOLATION_MODULES += $(LOCAL_MODULE)
    ifeq (true,$(PRODUCT_ENFORCE_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT))
      $(if $(filter $(LOCAL_MODULE),$(PRODUCT_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT_ALLOW_LIST)),,\
        $(call pretty-error,The module in product partition cannot be signed with certificate in system.))
    endif
  endif
endif

$(KATI_obsolete_var PRODUCT_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT_WHITELIST,Use PRODUCT_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT_ALLOW_LIST.)
