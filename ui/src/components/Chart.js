//@flow
import {map, zipObject} from 'lodash';
import React, {Component} from 'react';

import SimpleBar from './charts/SimpleBar';
import SimpleLine from './charts/SimpleLine';
import GroupedBar from './charts/GroupedBar';
import MultiLine from './charts/MultiLine';

import * as Types from '../types';

const CHARTS = {
  SimpleBar: SimpleBar,
  SimpleLine: SimpleLine,
  GroupedBar: GroupedBar,
  MultiLine: MultiLine,
};

class Chart extends Component {
  props: {
    name: string,
    rows: Array<Types.Row>,
    columns: Array<Types.Column>,
    onNewView: Object => void,
  };

  render() {
    const {name, columns, rows, onNewView} = this.props;

    const Handler = CHARTS[name];

    const columnNames = map(columns, 'name');

    const values = rows.map(row => {
      const cells = map(row.cells, 'value');

      return zipObject(columnNames, cells);
    });

    return (
      <div>
        <Handler values={values} onNewView={onNewView} />
      </div>
    );
  }
}

export default Chart;
