//@flow
import React, {Component} from 'react';

import Flatpickr from 'react-flatpickr';
import 'flatpickr/dist/flatpickr.min.css';
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
