import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import DocumentTitle from 'react-document-title';
import {Grid} from 'semantic-ui-react';

import ApiClient from '../../services/ApiClient';
import QueryResult from '../../components/QueryResult';
import QuestionForm from './QuestionForm';

import debounce from 'lodash/debounce';

class Playground extends Component {
  state = {
    currentDataSource: '',
    currentQuery: '',
    dataSource: '',
    query: '',
    showResult: false,
    querySuccess: false,
    variables: [],
  };

  evalQuery: Function;

  constructor(props) {
    super(props);
    this.evalQuery = debounce(this.evalQuery, 200, {trailing: true});
  }

  evalQuery(query, variables) {
    ApiClient.post('query/evaluate', {query, variables})
      .then(response => response.json())
      .then(({variables}) => {
        this.setState({variables});
      });
  }

  handleQuerySuccess = value => {
    this.setState({querySuccess: value});
  };

  handleDataSourcesChange = (e, {value}) => {
    this.setState({currentDataSource: value});
  };

  handleRunClick = e => {
    e.preventDefault();
    this.setState({showResult: true, query: this.state.currentQuery});
  };

  handleQueryChange = query => {
    this.setState({currentQuery: query});
    this.evalQuery(query, this.state.variables);
  };

  handleVariableChange = (key, value) => {
    const {variables} = this.state;
    const index = variables.findIndex(e => e['name'] == key);
    this.setState({
      variables: variables.map((item, i) => (index !== i ? item : {...item, value: value})),
    });
  };

  render() {
    const {settings} = this.props;

    return (
      <DocumentTitle title="Playground">
        <div>
          <Grid.Row>
            <Grid.Column>
              <QuestionForm
                saveEnabled={!settings.isReadOnly}
                handleDataSourcesChange={this.handleDataSourcesChange}
                handleQueryChange={this.handleQueryChange}
                handleRunClick={this.handleRunClick}
                currentDataSource={this.state.currentDataSource}
                currentQuery={this.state.currentQuery}
                querySuccess={this.state.querySuccess}
                variables={this.state.variables}
                handleVariableChange={this.handleVariableChange}
              />
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              {this.state.showResult && (
                <QueryResult
                  handleQuerySuccess={this.handleQuerySuccess}
                  source={this.state.currentDataSource}
                  query={this.state.query}
                  variables={this.state.variables}
                />
              )}
            </Grid.Column>
          </Grid.Row>
        </div>
      </DocumentTitle>
    );
  }
}

export default Playground;
