import {sortBy} from 'lodash';
import React, {Component} from 'react';
import {graphql, compose} from 'react-apollo';
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

  dataSourcesOptions = dataSources => {
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
    const {currentQuery, currentDataSource, history, saveQuestion, variables} = this.props;

    saveQuestion({
      variables: {
        title: this.state.title,
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
        </Form.Group>
        <Form.Group>
          <Form.Input onChange={this.handleTitleChange} value={this.state.title} width={13} />
          <Button
            icon="play"
            primary
            content="Run"
            disabled={!(currentQuery && currentQuery)}
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
  mutation SaveQuestion($title: String!, $query: String!, $dataSource: String!, $variables: [InputVariable]) {
    saveQuestion(title: $title, query: $query, dataSource: $dataSource, variables: $variables) {
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
