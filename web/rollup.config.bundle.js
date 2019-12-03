import svelte from 'rollup-plugin-svelte';
import replace from '@rollup/plugin-replace';
import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';
import { terser } from 'rollup-plugin-terser';

const mode = process.env.NODE_ENV;

export default {
    plugins: [
        replace({
            'process.browser': true,
            'process.env.NODE_ENV': JSON.stringify(mode)
        }),
        svelte(),
        resolve(),
        commonjs(),
        //terser(),
    ]
}
