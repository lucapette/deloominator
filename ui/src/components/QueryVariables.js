//@flow
import React, {Component} from 'react';

import * as types from '../types';

import DateTime from './variables/DateTime';
import Date from './variables/Date';

const VARIABLES = {
  timestamp: DateTime,
  date: Date,
};

type Props = {
  handleVariableChange: Function,
  variables: Array<types.Variable>,
};

class QueryVariables extends Component<Props> {
  render() {
    const {handleVariableChange, variables} = this.props;

    return variables.filter(v => v.isControllable).map(({name, value}) => {
      const Handler = VARIABLES[name];

      return <Handler key={name} name={name} value={value} handleVariableChange={handleVariableChange} />;
    });
  }
}

export default QueryVariables;
