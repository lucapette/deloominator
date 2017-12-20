import {kebabCase} from 'lodash';

const routing = {
  urlFor: (obj, keys = ['id']) => {
    const values = keys.map(k => obj[k]);
    return kebabCase(values.join('-'));
  },
};

export default routing;
