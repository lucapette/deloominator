//@flow
import {Table as FixedTable, Column, Cell} from 'fixed-data-table-2';
import React, {Component} from 'react';
import {Statistic} from 'semantic-ui-react';

import 'fixed-data-table-2/dist/fixed-data-table.min.css';

import * as Types from '../types';

type Props = {
  rows: Array<Types.Row>,
  columns: Array<Types.Column>,
};

class Table extends Component<Props> {
  render() {
    const {rows, columns} = this.props;

    const tableWidth = columns.length > 11 ? 1127 : columns.length * 100;

    return (
      <div>
        <Statistic size="mini" value={rows.length} label="rows" horizontal />

        <FixedTable rowHeight={50} rowsCount={rows.length} maxHeight={600} width={tableWidth} headerHeight={50}>
          {columns.map((column, i) => (
            <Column
              key={i}
              header={<Cell>{column.name}</Cell>}
              cell={props => <Cell {...props}>{rows[props.rowIndex].cells[i].value}</Cell>}
              width={100}
            />
          ))}
        </FixedTable>
      </div>
    );
  }
}

export default Table;
