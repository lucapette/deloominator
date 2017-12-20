import React, {Component} from 'react';
import {gql, graphql} from 'react-apollo';
import {withRouter} from 'react-router';
import {List, Loader} from 'semantic-ui-react';

import routing from '../../helpers/routing';

class QuestionListContainer extends Component {
  handleClick = (e, question) => {
    const questionPath = routing.urlFor(question, ['id', 'title']);
    this.props.history.push(`/questions/${questionPath}`);
  };

  render() {
    const {data: {loading, error, questions}} = this.props;

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

    return (
      <List selection verticalAlign="middle">
        {questions.map(question => (
          <List.Item key={question.id} onClick={e => this.handleClick(e, question)}>
            <List.Content>
              <List.Header>{question.title}</List.Header>
            </List.Content>
          </List.Item>
        ))}
      </List>
    );
  }
}

const Query = gql`
  {
    questions {
      id
      title
    }
  }
`;

const QuestionList = withRouter(graphql(Query)(QuestionListContainer));

export default QuestionList;
