import replace from '@rollup/plugin-replace';
import { terser } from "rollup-plugin-terser";

const mode = process.env.NODE_ENV;

let version = '<unknown>';
if (bazel_stamp_file) {
  const versionTag = require('fs')
                         .readFileSync(bazel_stamp_file, {encoding: 'utf-8'})
                         .split('\n')
                         .find(s => s.startsWith('VERSION'));
  if (versionTag) {
    version = versionTag.split(' ')[1].trim();
  }
}

export default {
    plugins: [
        replace({
            '__VERSION__': `${version}`,
            'process.browser': true,
            'process.env.NODE_ENV': JSON.stringify(mode)
        }),
        terser(),
    ]
}
