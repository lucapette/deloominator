'use strict';

const glob = require('glob');
const prettier = require('prettier');
const fs = require('fs');

const mode = process.argv[2] || 'check';
const shouldWrite = mode === 'write';

const patterns = ['**/*.js'];
const ignore = ['**/node_modules/**', 'dist/**'];

let didWarn = false;
let didError = false;

const files = glob.sync(patterns.join(','), {ignore});

files.forEach(file => {
  try {
    const input = fs.readFileSync(file, 'utf8');

    prettier.resolveConfig(file).then(options => {
      if (shouldWrite) {
        const output = prettier.format(input, options);
        if (output !== input) {
          fs.writeFileSync(file, output, 'utf8');
        }
      } else {
        if (!prettier.check(input, options)) {
          if (!didWarn) {
            console.log('Please run prettier');
            didWarn = true;
          }
          console.log(file);
        }
      }
    });
  } catch (error) {
    didError = true;
    console.log('\n\n' + error.message);
    console.log(file);
  }
});

if (didWarn || didError) {
  process.exit(1);
}
