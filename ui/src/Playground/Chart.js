import React, { Component } from 'react';

import SimpleBar from '../Charts/SimpleBar';
import SimpleLine from '../Charts/SimpleLine';

import _ from 'lodash';

const CHARTS = {
  'SimpleBar': SimpleBar,
  'SimpleLine': SimpleLine
};

export default class Chart extends Component {
  render() {
    const {chartName, columns, rows} = this.props.data;

    const Handler = CHARTS[chartName];

    const columnNames = _.map(columns, 'name');

    const values = rows.map(row => {
      const cells = _.map(row.cells, 'value');

      return _.zipObject(columnNames, cells);
    });

    return <Handler values={values}/>;
  }
}

