import React, {Component} from 'react';

import Flatpickr from 'react-flatpickr';
import 'flatpickr/dist/flatpickr.min.css';
import Date from './Date';

class DateTime extends Date {}

DateTime.defaultProps = Object.assign({}, Date.defaultProps, {options: {enableTime: true}});

export default DateTime;
