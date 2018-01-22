import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import {withRouter} from 'react-router';
import {Button, Form, Modal} from 'semantic-ui-react';

import Editor from '../../components/Editor';
import DataSources from '../../components/DataSources';
import QueryVariables from '../../components/QueryVariables';
import {urlFor} from '../../helpers/routing';

class SaveModal extends Component {
  handleInputChange = event => {
    const target = event.target;
    const name = target.name;
    const value = target.value;

    this.setState({
      [name]: value,
    });
  };

  handleSumbit = e => {
    e.preventDefault();
    const {handleSave} = this.props;

    const {title, description} = this.state;

    handleSave(title, description);
  };

  render() {
    const {querySuccess} = this.props;

    return (
      <Modal trigger={<Button icon="save" primary content="Save" disabled={!querySuccess} />}>
        <Modal.Header content="Save question" />
        <Modal.Content>
          <Form>
            <Form.Group widths="equal">
              <Form.Input
                fluid
                required
                name="title"
                label="Title"
                placeholder="Untitled question"
                onChange={this.handleInputChange}
              />
            </Form.Group>
            <Form.Group widths="equal">
              <Form.TextArea name="description" label="Description" onChange={this.handleInputChange} />
            </Form.Group>
          </Form>
          <Modal.Actions>
            <Button primary icon="save" content="Save" onClick={this.handleSumbit} />
            <Button content="Cancel" />
          </Modal.Actions>
        </Modal.Content>
      </Modal>
    );
  }
}

class QuestionFormContainer extends Component {
  handleSave = (title, description) => {
    const {query, dataSource, history, mutate, variables} = this.props;

    mutate({
      refetchQueries: ['questions'],
      variables: {
        title,
        description,
        query,
        dataSource,
        variables,
      },
    }).then(({data: {question}}) => {
      const questionPath = urlFor(question, ['id', 'title']);
      history.push(`/questions/${questionPath}`);
    });
  };

  render() {
    const {
      saveEnabled,
      handleDataSourcesChange,
      handleRunClick,
      handleQueryChange,
      query,
      dataSource,
      querySuccess,
      handleVariableChange,
      variables,
    } = this.props;

    return (
      <Form>
        <Form.Group>
          <DataSources handleDataSourcesChange={handleDataSourcesChange} />
          <Button icon="play" primary content="Run" disabled={!(dataSource && query)} onClick={handleRunClick} />
          {saveEnabled && (
            <SaveModal
              querySuccess={querySuccess}
              handleInputChange={this.handleInputChange}
              handleSave={this.handleSave}
            />
          )}
        </Form.Group>
        <Form.Group>
          <QueryVariables variables={variables} handleVariableChange={handleVariableChange} />
        </Form.Group>
        <Form.Group widths={16}>
          <Editor code={query} onChange={handleQueryChange} />
        </Form.Group>
      </Form>
    );
  }
}

const SaveQuestion = gql`
  mutation SaveQuestion(
    $title: String!
    $query: String!
    $dataSource: String!
    $description: String
    $variables: [InputVariable]
  ) {
    question(title: $title, query: $query, dataSource: $dataSource, description: $description, variables: $variables) {
      id
      title
    }
  }
`;

const QuestionForm = withRouter(graphql(SaveQuestion)(QuestionFormContainer));

export default QuestionForm;
