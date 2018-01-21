import React, {Component, Fragment} from 'react';
import DocumentTitle from 'react-document-title';
import {Grid} from 'semantic-ui-react';

import ApiClient from '../../services/ApiClient';
import QueryResult from '../../components/QueryResult';
import QuestionForm from './QuestionForm';

import debounce from 'lodash/debounce';

class Playground extends Component {
  constructor(props) {
    super(props);

    this.state = {
      currentDataSource: '',
      currentQuery: '',
      dataSource: '',
      query: '',
      showResult: false,
      querySuccess: false,
      variables: [],
    };

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
        <Fragment>
          <Grid.Row>
            <Grid.Column>
              <QuestionForm
                saveEnabled={!settings.isReadOnly}
                handleDataSourcesChange={this.handleDataSourcesChange}
                handleQueryChange={this.handleQueryChange}
                handleRunClick={this.handleRunClick}
                dataSource={this.state.currentDataSource}
                query={this.state.currentQuery}
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
                  dataSource={this.state.currentDataSource}
                  query={this.state.query}
                  variables={this.state.variables}
                />
              )}
            </Grid.Column>
          </Grid.Row>
        </Fragment>
      </DocumentTitle>
    );
  }
}

export default Playground;
