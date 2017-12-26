'use strict';

const glob = require('glob');
const prettier = require('prettier');
const fs = require('fs');

const mode = process.argv[2] || 'check';
const shouldWrite = mode === 'write';

const defaultOptions = {
  bracketSpacing: false,
  singleQuote: true,
  jsxBracketSameLine: true,
  trailingComma: 'all',
  printWidth: 120,
};
const config = {
  default: {
    patterns: ['**/*.js'],
    ignore: ['**/node_modules/**', 'dist/**'],
    options: {
      trailingComma: 'es5',
    },
  },
};

let didWarn = false;
let didError = false;

Object.keys(config).forEach(key => {
  const patterns = config[key].patterns;
  const options = config[key].options;
  const ignore = config[key].ignore;

  const globPattern = patterns.length > 1 ? `{${patterns.join(',')}}` : `${patterns.join(',')}`;
  const files = glob.sync(globPattern, {ignore});

  if (!files.length) {
    return;
  }

  const args = Object.assign({}, defaultOptions, options);
  files.forEach(file => {
    try {
      const input = fs.readFileSync(file, 'utf8');
      if (shouldWrite) {
        const output = prettier.format(input, args);
        if (output !== input) {
          fs.writeFileSync(file, output, 'utf8');
        }
      } else {
        if (!prettier.check(input, args)) {
          if (!didWarn) {
            console.log('Please run prettier');
            didWarn = true;
          }
          console.log(file);
        }
      }
    } catch (error) {
      didError = true;
      console.log('\n\n' + error.message);
      console.log(file);
    }
  });
});

if (didWarn || didError) {
  process.exit(1);
}
