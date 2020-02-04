import svelte from 'rollup-plugin-svelte'
import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import { terser } from 'rollup-plugin-terser'
import html from 'rollup-plugin-bundle-html'

const mode = process.env.NODE_ENV 
const isDevelopment = mode === "development"

const plugins = [
    svelte({
        dev: isDevelopment,
        css: css => css.write("./dist/app.css", isDevelopment)
    }),
    html({
        template: "src/index.html",
        dest: "dist",
        filename: "index.html",
        absolute: true
    }),
    commonjs(),
    resolve(),
    !isDevelopment && terser(),
]

module.exports = [{
	input: "src/main.js",
	output: {
		file: "dist/app.js",
		sourcemap: isDevelopment,
		format: "iife"
	},
	plugins
},{
	input: "src/sw.js",
	output: {
		file: "dist/sw.js",
		sourcemap: isDevelopment,
		format: "iife"
	},
	plugins
}]
