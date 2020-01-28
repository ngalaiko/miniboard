import svelte from 'rollup-plugin-svelte';
import replace from '@rollup/plugin-replace';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import { terser } from 'rollup-plugin-terser';
import html from 'rollup-plugin-bundle-html';
import glob from 'glob';
import fs from 'fs';

const mode = process.env.NODE_ENV;

class ResolvePbJS {
    async resolveId(importee) {
        return importee.endsWith('_pb.js')
            ? importee
            : null
    }

    async load(importee) {
        if (!importee.endsWith('_pb.js')) {
            return null
        }

        if (importee.startsWith('/')) {
            return null
        }

        return await new Promise((resolve, reject) => {
            glob(`**/*/${importee}`, {}, (err, files) => err === null ? resolve(files) : reject(err))
        }).then(files => {
            if (files.length === 0) throw new Error(`can't find ${importee}'`)
            return new Promise((resolve, reject) => {
                fs.readFile(files[0], (err, content) => err === null ? resolve(content.toString('utf8')) : reject(err))
            })
        })
    }
}

export default {
    plugins: [
        replace({
            'process.browser': true,
            'process.env.NODE_ENV': JSON.stringify(mode)
        }),
        svelte({
            css: css => css.write(__dirname.split("miniboard/")[1] + "/rollup/app.css", false)
        }),
        html({
            template: "web/app/src/index.html",
            dest: __dirname.split("miniboard/")[1] + "/rollup",
            filename: "index.html",
            absolute: true
        }),
        resolve(),
        new ResolvePbJS(),
        commonjs(),
        terser(),
    ]
}
