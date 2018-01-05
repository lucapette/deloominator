import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import DocumentTitle from 'react-document-title';
import {Container, Message, Loader, Grid, Header} from 'semantic-ui-react';

import QueryResult from '../../components/QueryResult';

class QuestionContainer extends Component {
  render() {
    const {data: {loading, error, question}} = this.props;

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

    return (
      <DocumentTitle title={question.title}>
        <Container>
          <Header as="h1">{question.title}</Header>
          <Grid.Row>
            <Grid.Column>
              <QueryResult source={question.dataSource} query={question.query} variables={question.variables} />
            </Grid.Column>
          </Grid.Row>
        </Container>
      </DocumentTitle>
    );
  }
}

const Query = gql`
  query Question($id: ID!) {
    question(id: $id) {
      id
      title
      query
      dataSource
    }
  }
`;

const Question = graphql(Query, {
  options: ({id}) => ({variables: {id}}),
})(QuestionContainer);

export default Question;
