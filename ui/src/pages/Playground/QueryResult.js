//@flow
import React, { Component } from "react";
import { gql, graphql } from "react-apollo";
import { Container, Table, Message, Loader, Segment, Divider } from "semantic-ui-react";

import Chart from "./Chart";
import RawResults from "./RawResults";

class QueryResultContainer extends Component {
  componentWillUpdate(nextProps, nextState) {
    const { data: { loading, error, query }, handleQuerySuccess } = nextProps;
    handleQuerySuccess(!(loading || error) && !(query != null && query.__typename == "queryError"));
  }

  render() {
    const { data: { loading, error, query }, handleQuerySuccess } = this.props;

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
