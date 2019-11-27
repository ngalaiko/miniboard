const rollup = require("rollup");
const path = require("path");
const input = process.argv[2];
const file = process.argv[3];

class ResolveInternal {
  resolveId(importee) {
    let parts = importee.split("/");

    if (parts[0] != "svelte") {
      return;
    }

    let subpackage = parts.length == 2 ? parts[1] : "internal";

    return path.join(
      __dirname,
      "../../../..",
      `svelte_deps/node_modules/svelte/${subpackage}.mjs`
    );
  }
}

async function build() {
  const inputOptions = { input, plugins: [new ResolveInternal()] };
  const outputOptions = {
    file,
    format: "iife",
    name: "svelte_app",
    sourcemap: true
  };
  const bundle = await rollup.rollup(inputOptions);

  await bundle.write(outputOptions);
}

build();
