const pickBy = require('lodash/pickBy');

module.exports = {
  DELOOMINATOR_PORT: 3000,
  ...pickBy(process.env, (v, k) => k.startsWith('DELOOMINATOR_')),
};
