not_system_module :=
ifeq (true,$(filter true, \
   $(LOCAL_PRODUCT_MODULE) $(LOCAL_PRODUCT_SERVICES_MODULE) \
   $(LOCAL_VENDOR_MODULE) $(LOCAL_PROPRIETARY_MODULE)))
  not_system_module := true
endif

$(if $(not_system_module), \
  $(if $(findstring $(dir $(DEFAULT_SYSTEM_DEV_CERTIFICATE)),$(LOCAL_CERTIFICATE)), \
  $(eval PACKAGES.$(LOCAL_MODULE).CERTIFICATE_VIOLATION :=true)))