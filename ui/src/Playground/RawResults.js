//@flow
import React, { Component } from 'react';

import { Table } from 'semantic-ui-react';

export default class RawResults extends Component {
  render() {
    return (
      <Table striped>
        <Table.Header>
          <Table.Row>
            {this.props.columns.map((c) => (<Table.HeaderCell key={c.name}>{c.name}</Table.HeaderCell>) )}
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {this.props.rows.map((r, i) =>(<Table.Row key={i}>{r.cells.map((c, j) => (<Table.Cell key={j}>{c.value}</Table.Cell>))}</Table.Row>))}
        </Table.Body>
      </Table>
    );
  }
}
