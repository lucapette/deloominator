import React, {Component} from 'react';

import Flatpickr from 'react-flatpickr';
import 'flatpickr/dist/flatpickr.min.css';

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
    return <Flatpickr value={this.state.value} onChange={this.handleChange} />;
  }
}

export default Date;
