import React, {Component, Fragment} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';

import {Message, Loader, Grid, Divider} from 'semantic-ui-react';

import Chart from './Chart';
import Table from './Table';

class QueryResultContainer extends Component {
  componentWillUpdate(nextProps) {
    const {data: {loading, error, query}, handleQuerySuccess} = nextProps;
    if (handleQuerySuccess) {
      handleQuerySuccess(!(loading || error) && !(query != null && query.__typename == 'queryError'));
    }
  }

  render() {
    const {data: {loading, error, query}, onNewView} = this.props;

    if (loading) {
      return (
        <Loader active inline="centered">
          Loading
        </Loader>
      );
    }

    if (error) {
      return <p>Error!</p>;
    }

    if (query.__typename == 'queryError') {
      return (
        <Message negative>
          <Message.Header>There is a problem with your query.</Message.Header>
          <p>{query.message}</p>
        </Message>
      );
    }

    if (query.chartName !== 'UnknownChart') {
      return (
        <Fragment>
          <Grid.Row>
            <Chart name={query.chartName} columns={query.columns} rows={query.rows} onNewView={onNewView} />
          </Grid.Row>
          <Divider hidden />
          <Grid.Row>
            <Table columns={query.columns} rows={query.rows} />
          </Grid.Row>
        </Fragment>
      );
    }

    return (
      <Grid.Row>
        <Table columns={query.columns} rows={query.rows} />
      </Grid.Row>
    );
  }
}

const Query = gql`
  query Query($source: String!, $query: String!, $variables: [InputVariable]) {
    query(source: $source, query: $query, variables: $variables) {
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
        variables {
          name
          value
          isControllable
        }
      }
      ... on queryError {
        message
      }
    }
  }
`;

const QueryResult = graphql(Query, {
  options: ({source, query, variables}) => ({variables: {source, query, variables: variables}}),
})(QueryResultContainer);

export default QueryResult;
