def _module_info_impl(ctx):
    # Do nothing for now
    pass

module_info = rule(
    implementation = _module_info_impl,
    attrs = {
        "installed": attr.label_list(),
        "module_name": attr.string(),
        "module_class": attr.string_list(),
        "module_path": attr.string_list(),
        "dependencies": attr.label_list(),
    }
)
