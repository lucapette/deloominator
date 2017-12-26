import React, {Component} from 'react';

import Date from './variables/Date';

const VARIABLES = {
  '{timestamp}': Date,
  '{yesterday}': Date,
};

class QueryVariables extends Component {
  render() {
    const {handleVariableChange} = this.props;

    return (
      <div>
        {Object.keys(this.props.variables).map(name => {
          const Handler = VARIABLES[name];

          return (
            <Handler
              key={name}
              name={name}
              value={this.props.variables[name]}
              handleVariableChange={handleVariableChange}
            />
          );
        })}
      </div>
    );
  }
}

export default QueryVariables;
