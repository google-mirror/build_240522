"""Event log tags generation rule"""
load("@bazel_skylib//lib:paths.bzl", "paths")

def _event_log_tags_impl(ctx):
    logtag_file = ctx.files.src[0]
    output_file = ctx.outputs.out
    ctx.actions.run(
        inputs = [logtag_file],
        outputs = [output_file],
        arguments =["-o",
                    output_file.path,
                    logtag_file.path],                    
        progress_message = "Generating Java logtag file from %s" % logtag_file.short_path,
        executable = ctx.executable.logtag_to_java_tool,
    )
    
event_log_tags = rule(
    implementation = _event_log_tags_impl,
    attrs = {
        "src": attr.label(allow_files = True, mandatory = True),
        "out": attr.output(mandatory = True),
        "logtag_to_java_tool": attr.label(
            executable = True,
            cfg = "exec",
            allow_files = True,
            default = Label("//build/make/tools:java-event-log-tags"),
        ),
        # "merge_tool": attr.label(
        #     executable = True,
        #     cfg = "exec",
        #     allow_files = True,
        #     default = Label("//build/make/tools:merge-event-log-tags"),
#        ),
    },
)
