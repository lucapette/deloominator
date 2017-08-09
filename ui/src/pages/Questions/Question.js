//@flow
import React, { Component } from "react";
import { gql, graphql } from "react-apollo";
import DocumentTitle from "react-document-title";
import { Container, Table, Message, Loader, Segment, Divider } from "semantic-ui-react";

class QuestionContainer extends Component {
  render() {
    const { data: { loading, error, question } } = this.props;

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

    return <DocumentTitle title={question.title} />;
  }
}

const Query = gql`
  query Question($id: ID!) {
    question(id: $id) {
      id
      title
      query
    }
  }
`;

const Question = graphql(Query, {
  options: ({ id }) => ({ variables: { id } }),
})(QuestionContainer);

export default Question;
