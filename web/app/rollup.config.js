import svelte from 'rollup-plugin-svelte';
import replace from '@rollup/plugin-replace';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import { terser } from 'rollup-plugin-terser';
import glob from 'glob';
import fs from 'fs';

const mode = process.env.NODE_ENV;

let apiUrl = 'http://localhost:8080'

if (bazel_stamp_file) {
    const versionTag = require('fs')
                        .readFileSync(bazel_stamp_file, {encoding: 'utf-8'})
                        .split('\n')
                        .find(s => s.startsWith('API_URL'));
    if (versionTag) {
        apiUrl = versionTag.split(' ')[1].trim()
    }
}

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

console.log(`api url: ${apiUrl}`)

export default {
    plugins: [
        replace({
            '__API_URL__': `${apiUrl}`,
            'process.browser': true,
            'process.env.NODE_ENV': JSON.stringify(mode)
        }),
        svelte(),
        resolve(),
        new ResolvePbJS(),
        commonjs(),
        terser(),
    ]
}
