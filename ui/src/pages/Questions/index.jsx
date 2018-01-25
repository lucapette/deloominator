import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import {withRouter} from 'react-router';
import DocumentTitle from 'react-document-title';
import {Loader, Card, Icon} from 'semantic-ui-react';

import distanceInWordsToNow from 'date-fns/distance_in_words_to_now';
import parse from 'date-fns/parse';

import {urlFor} from '../../helpers/routing';

const ExtraContent = ({question: {results}}) => {
  const iconName = results.chartName == 'UnknownChart' ? 'table' : 'area chart';

  return (
    <Card.Content extra>
      {results.__typename == 'queryError' ? (
        <p>
          <Icon name="broken chain" /> No results
        </p>
      ) : (
        <p>
          <Icon name={iconName} /> {results.rows.length}
        </p>
      )}
    </Card.Content>
  );
};

class QuestionListContainer extends Component {
  handleClick = question => {
    const questionPath = urlFor(question, ['id', 'title']);
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
      <DocumentTitle title="Q&A">
        <Card.Group>
          {questions.map(question => (
            <Card key={question.id} onClick={() => this.handleClick(question)}>
              <Card.Content>
                <Card.Header>{question.title}</Card.Header>
                <Card.Meta>{question.dataSource}</Card.Meta>
                <Card.Description>
                  Created <span className="date">{distanceInWordsToNow(parse(question.createdAt))}</span> ago
                </Card.Description>
                <ExtraContent question={question} />
              </Card.Content>
            </Card>
          ))}
        </Card.Group>
      </DocumentTitle>
    );
  }
}

const Query = gql`
  {
    questions {
      id
      title
      dataSource
      createdAt
      updatedAt
      results {
        ... on results {
          chartName
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
  }
`;

const Questions = withRouter(graphql(Query)(QuestionListContainer));

export default Questions;
