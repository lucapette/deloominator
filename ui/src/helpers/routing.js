import {kebabCase} from 'lodash';

export function urlFor(obj, keys = ['id']) {
  const values = keys.map(k => obj[k]);
  return kebabCase(values.join('-'));
}

export function idFromSlug(slug) {
  const [id] = slug.split('-');
  return id;
}
