load("//build_rules/svelte/internal:bundle-prod.bzl", _bundle_prod = "bundle_prod")
load("//build_rules/svelte/internal:svelte.bzl", _svelte = "svelte")

svelte = _svelte
bundle_prod = _bundle_prod
