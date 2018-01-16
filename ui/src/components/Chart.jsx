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

type Props = {
  name: string,
  rows: Array<Types.Row>,
  columns: Array<Types.Column>,
  onNewView: Object => void,
};

class Chart extends Component<Props> {
  render() {
    const {name, columns, rows, onNewView} = this.props;

    const Handler = CHARTS[name];

    const columnNames = map(columns, 'name');

    const values = rows.map(row => {
      const cells = map(row.cells, 'value');

      return zipObject(columnNames, cells);
    });

    return <Handler values={values} onNewView={onNewView} />;
  }
}

export default Chart;
