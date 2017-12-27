import React, {Component} from 'react';

import DateTime from './variables/DateTime';
import Date from './variables/Date';

const VARIABLES = {
  '{timestamp}': DateTime,
  '{date}': Date,
};

class QueryVariables extends Component {
  render() {
    const {handleVariableChange, variables} = this.props;

    return (
      <div>
        {Object.keys(variables).map(name => {
          const Handler = VARIABLES[name];

          return <Handler key={name} name={name} value={variables[name]} handleVariableChange={handleVariableChange} />;
        })}
      </div>
    );
  }
}

export default QueryVariables;
