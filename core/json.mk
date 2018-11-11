# Converts a list to a JSON list.
# $1: List separator.
# $2: List.
_json_list = [$(if $(2),"$(subst $(1),"$(comma)",$(2))")]

# Converts a space-separated list to a JSON list.
json_list = $(call _json_list,$(space),$(1))

# Converts a comma-separated list to a JSON list.
csv_to_json_list = $(call _json_list,$(comma),$(1))

4space :=$= $(space)$(space)$(space)$(space)

# 1: Key name
# 2: Value
add_json_val = $(eval _contents := $$(_contents)$$(_indent)"$$(strip $$(1))": $$(strip $$(2))$$(comma)$$(newline))
add_json_str = $(call add_json_val,$(1),"$(strip $(2))")
add_json_list = $(call add_json_val,$(1),$(call json_list,$(patsubst %,%,$(2))))
add_json_csv = $(call add_json_val,$(1),$(call csv_to_json_list,$(strip $(2))))
add_json_bool = $(call add_json_val,$(1),$(if $(strip $(2)),true,false))
add_json_raw = $(eval _contents := $(_contents)    $(1))

json_increase_indent = $(eval _indent := $$(_indent)$$(4space))
json_decrease_indent = $(eval _indent := $$(subst _,$$(space),$$(patsubst %____,%,$$(subst $$(space),_,$$(_indent)))))

add_json_map = $(eval _contents := $$(_contents)$$(_indent)"$$(strip $$(1))": {$$(newline))$(json_increase_indent)
end_json_map = $(json_decrease_indent)$(eval _contents := $$(_contents)$$(if $$(filter %$$(comma),$$(lastword $$(_contents))),__SV_END)$$(_indent)},$$(newline))

invert_bool = $(if $(strip $(1)),,true)

json_start = $(eval _contents := {$$(newline))$(eval _indent := $$(4space))

json_end = $(eval _contents := $$(subst $$(comma)$$(newline)__SV_END,$$(newline),$$(_contents)__SV_END}$$(newline)))

json_contents = $(_contents)