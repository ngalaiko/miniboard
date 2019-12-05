import svelte from 'rollup-plugin-svelte';
import replace from '@rollup/plugin-replace';
import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';
import { terser } from 'rollup-plugin-terser';

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
        commonjs(),
        terser(),
    ]
}
