load("//tools/svelte/internal:get-files.bzl", "get_files")

def _bundle_dev(ctx):
  files = get_files(ctx)
      
  args = ctx.actions.args()
  args.add(ctx.bin_dir.path + "/" + ctx.file.entry_point.path)
  args.add(ctx.outputs.build.path)

  ctx.actions.run(
      executable = ctx.executable._rollup,
      inputs = files,
      outputs = [ctx.outputs.build, ctx.outputs.map],
      arguments = [args]
  )
 
bundle_dev = rule(
  implementation = _bundle_dev,
  attrs = {
    "deps": attr.label_list(),
    "entry_point": attr.label(allow_single_file = True),
    "_rollup": attr.label(
        executable = True,
        cfg="host",
        default = Label("//tools/svelte/internal:rollup")),
    },
    outputs = {"build": "%{name}.es6.js", "map": "%{name}.es6.js.map"}
)
