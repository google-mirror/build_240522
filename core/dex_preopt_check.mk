# Checks that some critical dexpreopt output files are installed.
# $(1): the full list of modules going to be installed
define dexpreopt-install-check
	$(foreach artifact,$(DEXPREOPT_SYSTEMSERVER_ARTIFACTS), \
		$(if $(filter $(artifact),$(1)),, \
			$(error Missing compilation artifact $(artifact). Dexpreopting is not working for some system server jars) \
		) \
	)
endef
