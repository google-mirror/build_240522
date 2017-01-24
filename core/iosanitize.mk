define find_in_iosanitize_never_projects
$(subst $(space),, \
 $(foreach project,$(IOSANITIZE_NEVER_PROJECTS), \
  $(if $(filter $(project)%,$(1)),$(project)) \
 ) \
)
endef
