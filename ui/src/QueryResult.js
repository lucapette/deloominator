//@flow
import React, { Component } from 'react'

import { Container, Table } from 'semantic-ui-react'

import { gql, graphql } from 'react-apollo';

class QueryResultContainer extends Component {
  render() {
    const { data: { loading, error, query } } = this.props;

    if (loading) {
      return <p>Loading...</p>;
    }

    if (error) {
      return <p>Error!</p>;
    }

    if (query.__typename == 'queryError') {
      return <p>{query.message}</p>;
    }

    return (
      <Container>
        <Table striped>
          <Table.Header>
            <Table.Row>
              {query.columns.map((c) => (<Table.HeaderCell key={c.name}>{c.name}</Table.HeaderCell>) )}
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {query.rows.map((r, i) =>(<Table.Row key={i}>{r.cells.map((c, j) => (<Table.Cell key={j}>{c.value}</Table.Cell>))}</Table.Row>))}
          </Table.Body>
        </Table>
      </Container>
    );
  }
}

const Query = gql`
query Query($source: String!, $input: String!) {
  query(source: $source, input: $input) {
    ... on rawResults {
        total
        columns {
          name
          type
        }

        rows {
          cells {
            value
          }
        }
    }

    ... on queryError { message }
  }
}`;

const QueryResult = graphql(Query, {
  options: ({source, input}) => ({ variables: {source, input} }),
})(QueryResultContainer);

export default QueryResult;
