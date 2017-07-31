//@flow
import React, { Component } from "react";

import { Button, Form } from "semantic-ui-react";

import { sortBy } from "lodash";

class QuestionForm extends Component {
  dataSourcesOptions = (dataSources: [string]) => {
    return sortBy(dataSources || [], ["name"], ["asc"]).map(s => ({
      name: s.name,
      text: s.name,
      value: s.name,
    }));
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
          <Button icon="save" primary content="Save" disabled={!querySuccess} onClick={handleRunClick} />
        </Form.Group>
        <Form.TextArea placeholder="Write your query here" value={currentQuery} onChange={handleQueryChange} />
      </Form>
    );
  }
}

export default QuestionForm;
