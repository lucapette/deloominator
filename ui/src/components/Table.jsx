//@flow
import {Table as FixedTable, Column, Cell} from 'fixed-data-table-2';
import React, {Component, Fragment} from 'react';

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
      <Fragment>
        <span>
          <b>
            <i>{rows.length}</i>
          </b>{' '}
          results found.
        </span>
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
      </Fragment>
    );
  }
}

export default Table;
