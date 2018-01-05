//@flow

import Date from './Date';

class DateTime extends Date {
  static defaultProps = {
    ...Date.defaultProps,
    ...{
      options: {
        dateFormat: 'Z',
        enableTime: true,
      },
    },
  };
}

export default DateTime;
