import React, {Component} from 'react';
import {graphql} from 'react-apollo';
import gql from 'graphql-tag';
import DocumentTitle from 'react-document-title';
import {Container, Message, Loader, Grid, Header, Form} from 'semantic-ui-react';

import QueryResult from '../../components/QueryResult';
import QueryVariables from '../../components/QueryVariables';

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

    const {title, dataSource, query, variables} = question;

    return (
      <DocumentTitle title={title}>
        <Container>
          <Header as="h1">{title}</Header>
          <Grid.Row>
            <Grid.Column>
              <Form>
                <Form.Group>
                  <QueryVariables variables={this.state.variables} handleVariableChange={this.handleVariableChange} />
                </Form.Group>
              </Form>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              <QueryResult source={dataSource} query={query} variables={this.state.variables} />
            </Grid.Column>
          </Grid.Row>
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
