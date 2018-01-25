import {urlFor, idFromSlug} from '../../helpers/routing';

describe('urlFor', () => {
  test('it returns the id by default', () => {
    const object = {id: 42};
    expect(urlFor(object)).toBe('42');
  });

  test('it returns a URL for the specified keys', () => {
    const object = {id: 42, name: 'Grace', surname: 'Hopper'};
    expect(urlFor(object, ['id', 'name', 'surname'])).toBe('42-grace-hopper');
  });

  test('it returns a kebabCases URL for the specified keys', () => {
    const object = {id: 42, title: 'The answer is forty two'};
    expect(urlFor(object, ['id', 'title'])).toBe('42-the-answer-is-forty-two');
  });
});

describe('idFromSlug', () => {
  test('it returns the id', () => {
    expect(idFromSlug('42-is-the-answer')).toBe('42');
  });
});
