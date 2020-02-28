import svelte from 'rollup-plugin-svelte'
import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import { terser } from 'rollup-plugin-terser'
import html from 'rollup-plugin-bundle-html'
import typescript from 'rollup-plugin-typescript2'
import typescriptCompiler from 'typescript'
import sveltePreprocessor from 'svelte-preprocess'

const mode = process.env.NODE_ENV 
const isDevelopment = mode === "development"

const name = isDevelopment ? 'app' : 'app-[hash]'

const appPlugins = [
    svelte({
        dev: isDevelopment,
        preprocess: sveltePreprocessor({}),
        css: css => css.write(`./dist/${name}.css`, isDevelopment)
    }),
    html({
        template: "src/index.html",
        dest: "dist",
        filename: "index.html",
        absolute: true
    }),
    typescript({ typescript: typescriptCompiler }),
    commonjs(),
    resolve(),
    !isDevelopment && terser(),
]

const swPlugins = [
    typescript({ typescript: typescriptCompiler }),
    commonjs(),
    !isDevelopment && terser(),
]

module.exports = [{
	input: "src/index.ts",
	output: {
		file: `dist/${name}.js`,
		sourcemap: isDevelopment,
		format: "iife"
	},
    plugins: appPlugins
},{
	input: "src/sw.ts",
	output: {
		file: "dist/sw.js",
		sourcemap: isDevelopment,
		format: "iife"
	},
    plugins: swPlugins
}]
