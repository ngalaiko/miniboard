load("//tools/svelte/internal:bundle-dev.bzl", _bundle_dev = "bundle_dev")
load("//tools/svelte/internal:bundle-prod.bzl", _bundle_prod = "bundle_prod")
load("//tools/svelte/internal:svelte.bzl", _svelte = "svelte")

svelte = _svelte
bundle_prod = _bundle_prod
bundle_dev = _bundle_dev
