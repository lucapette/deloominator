//@flow
import React, {Component} from 'react';

import DateTime from './variables/DateTime';
import Date from './variables/Date';

const VARIABLES = {
  '{timestamp}': DateTime,
  '{date}': Date,
};

type Props = {
  handleVariableChange: Function,
  variables: {},
};

class QueryVariables extends Component<Props> {
  render() {
    const {handleVariableChange, variables} = this.props;

    const controllableNames = Object.keys(VARIABLES);

    return Object.keys(variables)
      .filter(name => controllableNames.includes(name))
      .map(name => {
        const Handler = VARIABLES[name];

        return <Handler key={name} name={name} value={variables[name]} handleVariableChange={handleVariableChange} />;
      });
  }
}

export default QueryVariables;
