import {sortBy} from 'lodash';
import React, {Component} from 'react';
import {graphql, compose} from 'react-apollo';
import gql from 'graphql-tag';
import {withRouter} from 'react-router';
import {Button, Form, Modal} from 'semantic-ui-react';

import Editor from '../../components/Editor';
import QueryVariables from '../../components/QueryVariables';
import routing from '../../helpers/routing';

class QuestionFormContainer extends Component {
  dataSourcesOptions = dataSources => {
    return sortBy(dataSources || [], ['name'], ['asc']).map(s => ({
      name: s.name,
      text: s.name,
      value: s.name,
    }));
  };

  handleInputChange = event => {
    const target = event.target;
    const name = target.name;
    const value = target.value;

    this.setState({
      [name]: value,
    });
  };

  handleSave = e => {
    e.preventDefault();
    const {currentQuery, currentDataSource, history, saveQuestion, variables} = this.props;

    saveQuestion({
      refetchQueries: ['questions'],
      variables: {
        title: this.state.title,
        description: this.state.description,
        query: currentQuery,
        dataSource: currentDataSource,
        variables: variables,
      },
    }).then(({data: {saveQuestion}}) => {
      const questionPath = routing.urlFor(saveQuestion, ['id', 'title']);
      history.push(`/questions/${questionPath}`);
    });
  };

  render() {
    const {
      saveEnabled,
      handleDataSourcesChange,
      handleRunClick,
      handleQueryChange,
      currentQuery,
      currentDataSource,
      querySuccess,
      handleVariableChange,
      variables,
      data: {loading, error, dataSources},
    } = this.props;

    if (error) {
      return <p>Error!</p>;
    }

    return (
      <Form>
        <Form.Group>
          <Form.Dropdown
            loading={loading}
            placeholder="Data Source"
            search
            selection
            onChange={handleDataSourcesChange}
            options={this.dataSourcesOptions(dataSources)}
          />
          <Button
            icon="play"
            primary
            content="Run"
            disabled={!(currentDataSource && currentQuery)}
            onClick={handleRunClick}
          />
          {saveEnabled && (
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
                  <Button primary icon="save" content="Save" onClick={this.handleSave} />
                  <Button content="Cancel" />
                </Modal.Actions>
              </Modal.Content>
            </Modal>
          )}
        </Form.Group>
        <Form.Group>
          <QueryVariables variables={variables} handleVariableChange={handleVariableChange} />
        </Form.Group>
        <Form.Group widths={16}>
          <Editor code={currentQuery} onChange={handleQueryChange} />
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
    saveQuestion(
      title: $title
      query: $query
      dataSource: $dataSource
      description: $description
      variables: $variables
    ) {
      id
      title
    }
  }
`;

const Query = gql`
  {
    dataSources {
      name
    }
  }
`;

const QuestionForm = withRouter(
  compose(graphql(SaveQuestion, {name: 'saveQuestion'}), graphql(Query))(QuestionFormContainer),
);

export default QuestionForm;
