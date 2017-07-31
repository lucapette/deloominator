//@flow
import React, { Component } from "react";

import { Container, Table, Message, Loader, Segment, Divider } from "semantic-ui-react";

import { gql, graphql } from "react-apollo";

import RawResults from "./RawResults";

import Chart from "./Chart";

class QueryResultContainer extends Component {
  render() {
    const { data: { loading, error, query }, handleQuerySuccess } = this.props;

    handleQuerySuccess(false);

    if (loading) {
      return (
        <Container>
          <Loader active inline="centered">
            Loading
          </Loader>
        </Container>
      );
    }

    if (error) {
      return <p>Error!</p>;
    }

    if (query.__typename == "queryError") {
      return (
        <Message negative>
          <Message.Header>There is a problem with your query.</Message.Header>
          <p>
            {query.message}
          </p>
        </Message>
      );
    }

    handleQuerySuccess(true);

    if (query.chartName !== "UnknownChart") {
      return (
        <Container>
          <Segment padded>
            <Chart data={query} />
            <Divider horizontal>Raw data</Divider>
            <RawResults columns={query.columns} rows={query.rows} />
          </Segment>
        </Container>
      );
    }

    return (
      <Container>
        <Segment padded>
          <RawResults columns={query.columns} rows={query.rows} />
        </Segment>
      </Container>
    );
  }
}

const Query = gql`
  query Query($source: String!, $input: String!) {
    query(source: $source, input: $input) {
      ... on results {
        chartName
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
      ... on queryError {
        message
      }
    }
  }
`;

const QueryResult = graphql(Query, {
  options: ({ source, input }) => ({ variables: { source, input } }),
})(QueryResultContainer);

export default QueryResult;
