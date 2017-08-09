//@flow
import { sortBy, kebabCase } from "lodash";
import React, { Component } from "react";
import { gql, graphql } from "react-apollo";
import { Button, Form } from "semantic-ui-react";
import { withRouter } from "react-router";

class QuestionFormContainer extends Component {
  state = {
    title: "",
  };

  dataSourcesOptions = (dataSources: [string]) => {
    return sortBy(dataSources || [], ["name"], ["asc"]).map(s => ({
      name: s.name,
      text: s.name,
      value: s.name,
    }));
  };

  handleTitleChange = e => {
    this.setState({ title: e.target.value });
  };

  handleSave = e => {
    e.preventDefault();
    this.props
      .mutate({ variables: { title: this.state.title, query: this.props.currentQuery } })
      .then(({ data: { saveQuestion: { id, title } } }) => {
        const questionPath = kebabCase(`${id}-${title}`);
        this.props.history.push(`/questions/${questionPath}`);
      })
      .catch(({ error }) => {
        console.log(error);
      });
  };

  render() {
    const {
      handleDataSourcesChange,
      handleRunClick,
      handleQueryChange,
      dataSources,
      selectedDataSource,
      currentQuery,
      querySuccess,
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
          <Form.Input
            placeholder="Untitled visualization"
            onChange={this.handleTitleChange}
            value={this.state.title}
            width={13}
          />
          <Button
            icon="play"
            primary
            content="Run"
            disabled={!(selectedDataSource && currentQuery)}
            onClick={handleRunClick}
          />
          <Button icon="save" primary content="Save" disabled={!querySuccess} onClick={this.handleSave} />
        </Form.Group>
        <Form.TextArea placeholder="Write your query here" value={currentQuery} onChange={handleQueryChange} />
      </Form>
    );
  }
}

const SaveQuestion = gql`
  mutation SaveQuestion($title: String!, $query: String!) {
    saveQuestion(title: $title, query: $query) {
      id
      title
    }
  }
`;

const QuestionForm = withRouter(graphql(SaveQuestion)(QuestionFormContainer));

export default QuestionForm;
