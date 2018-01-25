import React, {Component, Fragment} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import {connect} from 'react-redux';

import {Message, Loader, Grid, Divider} from 'semantic-ui-react';

import Chart from './Chart';
import Table from './Table';
import * as actions from '../actions/queryEditor';

class QueryResultContainer extends Component {
  componentWillUpdate(nextProps) {
    const {data: {loading, error, query}, setQueryWasSuccessful} = nextProps;

    setQueryWasSuccessful(!(loading || error) && !(query != null && query.__typename == 'queryError'));
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

    return (
      <Fragment>
        {query.chartName !== 'UnknownChart' && (
          <Chart name={query.chartName} columns={query.columns} rows={query.rows} onNewView={onNewView} />
        )}
        <Divider hidden />
        <Grid.Row>
          <Table columns={query.columns} rows={query.rows} />
        </Grid.Row>
      </Fragment>
    );
  }
}

const Query = gql`
  query Query($dataSource: String!, $query: String!, $variables: [InputVariable]) {
    query(dataSource: $dataSource, query: $query, variables: $variables) {
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

const QueryResult = connect(null, {setQueryWasSuccessful: actions.setQueryWasSuccessful})(
  graphql(Query, {
    options: ({dataSource, query, variables}) => ({variables: {dataSource, query, variables: variables}}),
  })(QueryResultContainer),
);

export default QueryResult;
