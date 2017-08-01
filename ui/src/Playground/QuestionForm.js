//@flow
import { sortBy } from "lodash";
import React, { Component } from "react";
import { gql, graphql } from "react-apollo";
import { Button, Form } from "semantic-ui-react";

class QuestionFormContainer extends Component {
  dataSourcesOptions = (dataSources: [string]) => {
    return sortBy(dataSources || [], ["name"], ["asc"]).map(s => ({
      name: s.name,
      text: s.name,
      value: s.name,
    }));
  };

  handleSave = e => {
    e.preventDefault();
    this.props.mutate({ variables: { title: "nonworkingmutation", query: this.props.currentQuery } });
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
          <Form.Input placeholder="Untitled visualization" width={13} />
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
      query
    }
  }
`;

const QuestionForm = graphql(SaveQuestion)(QuestionFormContainer);

export default QuestionForm;
