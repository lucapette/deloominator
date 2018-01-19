//@flow
import React, {Component} from 'react';

import Flatpickr from 'react-flatpickr';
import 'flatpickr/dist/flatpickr.min.css';

import {capitalize} from 'lodash';

type Props = {
  name: string,
  handleVariableChange: Function,
  options: {},
  value: string,
};

type State = {
  value: string,
};

class Date extends Component<Props, State> {
  static defaultProps = {
    options: {
      dateFormat: 'Y-m-d',
    },
  };

  constructor(props: Props) {
    super(props);
    this.state = {
      value: props.value,
    };
  }

  handleChange = (selectedDates: Array<Date>, selected: string) => {
    const {name, handleVariableChange} = this.props;

    handleVariableChange(name, selected);
  };

  render() {
    const {options, name} = this.props;
    return (
      <div className="ui field labeled input">
        <div className="ui label">{capitalize(name.replace(/[{}]/g, ''))}</div>
        <Flatpickr options={options} value={this.state.value} onChange={this.handleChange} />
      </div>
    );
  }
}

export default Date;
