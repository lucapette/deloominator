import React, {Component} from 'react';
import {Button, Form} from 'semantic-ui-react';
import {connect} from 'react-redux';
import debounce from 'lodash/debounce';

import * as actions from '../../actions/queryEditor';

import ApiClient from '../../services/ApiClient';

import Editor from '../../components/Editor';
import DataSources from '../../components/DataSources';
import QueryVariables from '../../components/QueryVariables';

import SaveModal from './SaveModal';

class QuestionFormContainer extends Component {
  constructor(props) {
    super(props);

    this.evalQuery = debounce(this.evalQuery, 200, {trailing: true});
  }

  evalQuery() {
    const {queryDraft, variables, setVariables} = this.props;

    ApiClient.post('query/evaluate', {query: queryDraft, variables})
      .then(response => response.json())
      .then(({variables}) => {
        setVariables(variables);
      });
  }

  handleVariableChange = (key, value) => {
    this.props.setVariable(key, value);
  };

  handleDataSourcesChange = (e, {value}) => {
    this.props.setInputValue('dataSource', value);
  };

  handleRunClick = e => {
    e.preventDefault();
    const {setInputValue, queryDraft} = this.props;
    setInputValue('query', queryDraft);
  };

  handleQueryChange = queryDraft => {
    this.props.setInputValue('queryDraft', queryDraft);
    this.evalQuery();
  };

  render() {
    const {saveEnabled, queryDraft, dataSource, variables} = this.props;

    return (
      <Form>
        <Form.Group>
          <DataSources handleDataSourcesChange={this.handleDataSourcesChange} dataSource={dataSource} />
          <Button
            icon="play"
            primary
            content="Run"
            disabled={!(dataSource && queryDraft)}
            onClick={this.handleRunClick}
          />
          {saveEnabled && <SaveModal handleInputChange={this.handleInputChange} handleSave={this.handleSave} />}
        </Form.Group>
        <Form.Group>
          <QueryVariables variables={variables} handleVariableChange={this.handleVariableChange} />
        </Form.Group>
        <Form.Group widths={16}>
          <Editor code={queryDraft} onChange={this.handleQueryChange} />
        </Form.Group>
      </Form>
    );
  }
}

const mapStateToProps = state => {
  return {...state.queryEditor};
};

const QuestionForm = connect(mapStateToProps, {
  setInputValue: actions.setInputValue,
  setVariables: actions.setVariables,
  setVariable: actions.setVariable,
})(QuestionFormContainer);

export default QuestionForm;
