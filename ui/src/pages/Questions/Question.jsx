import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import DocumentTitle from 'react-document-title';
import {Container, Loader, Grid, Header, Dropdown} from 'semantic-ui-react';

import ApiClient from '../../services/ApiClient';

import QueryResult from '../../components/QueryResult';
import QueryVariables from '../../components/QueryVariables';
import Menu from 'semantic-ui-react/dist/commonjs/collections/Menu/Menu';

class QuestionContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      variables: [],
    };
  }

  componentWillReceiveProps(nextProps) {
    const {data: {loading, error, question}} = nextProps;

    if (loading || error) {
      return;
    }

    const {variables} = question;

    this.setState({variables: variables.map(v => ({name: v.name, value: v.value, isControllable: v.isControllable}))});
  }

  handleVariableChange = (key, value) => {
    const {variables} = this.state;
    const index = variables.findIndex(e => e['name'] == key);
    this.setState({variables: variables.map((item, i) => (index !== i ? item : {...item, value: value}))});
  };

  exportCSV = e => {
    e.preventDefault();
    const {data: {question: {dataSource, query}}} = this.props;

    ApiClient.post('export/csv', {Source: dataSource, Query: query})
      .then(response => response.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(new Blob([blob], {type: {type: 'text/csv;charset=utf-8;'}}));
        let a = document.createElement('a');
        a.href = url;
        a.download = 'chart.csv';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
      });
  };

  onNewView = view => {
    this.view = view;
  };

  exportPNG = e => {
    e.preventDefault();
    this.view.toImageURL('png').then(url => {
      let link = document.createElement('a');
      link.setAttribute('href', url);
      link.setAttribute('target', '_blank');
      link.setAttribute('download', 'chart.png');
      link.dispatchEvent(new MouseEvent('click'));
    });
  };

  render() {
    const {data: {loading, error, question}} = this.props;

    if (loading) {
      return (
        <Container>
          <Loader active inline="centered">
            Loading
          </Loader>
        </Container>
      );
    }

    if (error) {
      return <p>Error!</p>;
    }

    const {title, dataSource, query} = question;

    const {variables} = this.state;

    return (
      <DocumentTitle title={title}>
        <Container>
          <Header as="h1">
            {title}
            <Header.Subheader>{dataSource}</Header.Subheader>
          </Header>
          <Grid.Row>
            <Menu borderless secondary>
              <Menu.Item>
                <QueryVariables variables={variables} handleVariableChange={this.handleVariableChange} />
              </Menu.Item>
              <Menu.Menu position="right">
                <Dropdown item icon="download">
                  <Dropdown.Menu>
                    <Dropdown.Item onClick={this.exportPNG}>PNG</Dropdown.Item>
                    <Dropdown.Item onClick={this.exportCSV}>CSV</Dropdown.Item>
                  </Dropdown.Menu>
                </Dropdown>
              </Menu.Menu>
            </Menu>
          </Grid.Row>
          <QueryResult source={dataSource} query={query} variables={variables} onNewView={this.onNewView} />
        </Container>
      </DocumentTitle>
    );
  }
}

const Query = gql`
  query Question($id: ID!) {
    question(id: $id) {
      id
      title
      query
      dataSource
      variables {
        name
        value
        isControllable
      }
    }
  }
`;

const Question = graphql(Query, {
  options: ({id}) => ({variables: {id}}),
})(QuestionContainer);

export default Question;
