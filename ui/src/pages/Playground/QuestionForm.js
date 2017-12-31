import {sortBy} from 'lodash';
import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import {withRouter} from 'react-router';
import {Button, Form} from 'semantic-ui-react';

import Editor from '../../components/Editor';
import QueryVariables from '../../components/QueryVariables';
import routing from '../../helpers/routing';

class QuestionFormContainer extends Component {
  state = {
    title: 'Untitled visualization',
  };

  dataSourcesOptions = (dataSources: [string]) => {
    return sortBy(dataSources || [], ['name'], ['asc']).map(s => ({
      name: s.name,
      text: s.name,
      value: s.name,
    }));
  };

  handleTitleChange = e => {
    this.setState({title: e.target.value});
  };

  handleSave = e => {
    e.preventDefault();
    const {currentQuery, selectedDataSource, history, mutate, variables} = this.props;

    mutate({
      variables: {
        title: this.state.title,
        query: currentQuery,
        dataSource: selectedDataSource,
        variables: JSON.stringify(variables),
      },
    })
      .then(({data: {saveQuestion}}) => {
        const questionPath = routing.urlFor(saveQuestion, ['id', 'title']);
        history.push(`/questions/${questionPath}`);
      })
      .catch(({error}) => {
        console.log(error);
      });
  };

  render() {
    const {
      saveEnabled,
      handleDataSourcesChange,
      handleRunClick,
      handleQueryChange,
      dataSources,
      selectedDataSource,
      currentQuery,
      querySuccess,
      handleVariableChange,
      variables,
    } = this.props;

    return (
      <Form>
        <Form.Group>
          <Form.Dropdown
            placeholder="Data Source"
            search
            selection
            onChange={handleDataSourcesChange}
            options={this.dataSourcesOptions(dataSources)}
          />
        </Form.Group>
        <Form.Group>
          <Form.Input onChange={this.handleTitleChange} value={this.state.title} width={13} />
          <Button
            icon="play"
            primary
            content="Run"
            disabled={!(selectedDataSource && currentQuery)}
            onClick={handleRunClick}
          />
          {saveEnabled && (
            <Button icon="save" primary content="Save" disabled={!querySuccess} onClick={this.handleSave} />
          )}
        </Form.Group>
        <Form.Group>
          <QueryVariables variables={variables} handleVariableChange={handleVariableChange} />
        </Form.Group>
        <Form.Group>
          <Editor code={currentQuery} onChange={handleQueryChange} />
        </Form.Group>
      </Form>
    );
  }
}

const SaveQuestion = gql`
  mutation SaveQuestion($title: String!, $query: String!, $dataSource: String!, $variables: String) {
    saveQuestion(title: $title, query: $query, dataSource: $dataSource, variables: $variables) {
      id
      title
    }
  }
`;

const QuestionForm = withRouter(graphql(SaveQuestion)(QuestionFormContainer));

export default QuestionForm;
