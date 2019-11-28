
SvelteFiles = provider("transitive_sources")

def get_transitive_srcs(srcs, deps):
    
  return depset(
      srcs,
      transitive = [dep[SvelteFiles].transitive_sources for dep in deps],
  )

def _svelte(ctx):
  args = ctx.actions.args()
  args.add(ctx.file.entry_point.path)
  args.add(ctx.outputs.build.path)

  ctx.actions.run(
      mnemonic = "Svelte",
      executable = ctx.executable._svelte,
      outputs = [ctx.outputs.build],
      inputs = [ctx.file.entry_point],
      arguments = [args]
  )

  trans_srcs = get_transitive_srcs(ctx.files.srcs + [ctx.outputs.build], ctx.attr.deps)
  
  return [
      SvelteFiles(transitive_sources = trans_srcs),
      DefaultInfo(files = trans_srcs),
  ]

svelte = rule(
  implementation = _svelte,
  attrs = {
  "entry_point": attr.label(allow_single_file = True),
  "deps": attr.label_list(),
  "srcs": attr.label_list(allow_files = True),
  "_svelte": attr.label(
        default=Label("//build_rules/svelte/internal:svelte"),
        executable=True,
        cfg="host"),
  },
  outputs = {
      "build": "%{name}.svelte.js"
  }
)
