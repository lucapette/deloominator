import React, {Component} from 'react';

import Flatpickr from 'react-flatpickr';
import 'flatpickr/dist/flatpickr.min.css';
import {Label, Form} from 'semantic-ui-react';

class Date extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: props.value,
    };
  }

  handleChange = selectedDates => {
    const {name, handleVariableChange} = this.props;

    handleVariableChange(name, selectedDates[0]);
  };

  render() {
    const {options, name} = this.props;
    return (
      <div className="field">
        <label>{name}</label>
        <Flatpickr options={options} value={this.state.value} onChange={this.handleChange} />
      </div>
    );
  }
}

Date.defaultProps = {
  options: {
    dateFormat: 'Y-m-d',
  },
};

export default Date;
