//@flow
import React, {Component} from 'react';
import {gql, graphql} from 'react-apollo';
import DocumentTitle from 'react-document-title';
import {Grid} from 'semantic-ui-react';

import ApiClient from '../../services/ApiClient';
import QueryResult from '../../components/QueryResult';
import QuestionForm from './QuestionForm';

class PlaygroundPage extends Component {
  state = {
    selectedDataSource: '',
    currentQuery: '',
    dataSource: '',
    query: '',
    showResult: false,
    querySuccess: false,
    variables: {},
  };

  handleQuerySuccess = value => {
    this.setState({querySuccess: value});
  };

  handleDataSourcesChange = (e, {value}) => {
    this.setState({selectedDataSource: value});
  };

  handleRunClick = e => {
    e.preventDefault();
    this.setState({
      showResult: true,
      query: this.state.currentQuery,
      dataSource: this.state.selectedDataSource,
    });
  };

  handleQueryChange = query => {
    // TODO Use settings endpoint to contruct a regex and run the endpoint only
    // for known variables

    ApiClient.post('query/evaluate', {query: query})
      .then(response => response.json())
      .then(variables => {
        this.setState({variables});
      });

    this.setState({currentQuery: query});
  };

  handleVariableChange = (key, value) => {
    let variables = Object.assign({}, this.state.variables);

    variables[key] = value;

    this.setState({variables});
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
                dataSources={this.props.data.dataSources}
                handleDataSourcesChange={this.handleDataSourcesChange}
                handleQueryChange={this.handleQueryChange}
                handleRunClick={this.handleRunClick}
                selectedDataSource={this.state.selectedDataSource}
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
                  source={this.state.dataSource}
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

const Query = gql`
  {
    dataSources {
      name
    }
  }
`;

const Playground = graphql(Query)(PlaygroundPage);

export default Playground;
