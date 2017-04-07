//@flow
import React, { Component } from 'react'
import DocumentTitle from 'react-document-title'

import { Button, Container, Form, Grid } from 'semantic-ui-react'

import { gql, graphql } from 'react-apollo';

class Playground extends Component {
  dataSourcesOptions = (data) => {
    return data.dataSources ? data.dataSources.map((s) => ({text: s.name, value: s.name})) : [];
  }

  render() {
    return (
      <DocumentTitle title='Playground'>
        <Container>
          <Grid.Row>
            <Grid.Column>
              <Form>
                <Form.Group>
                  <Form.Dropdown placeholder='Data Source' search selection options={this.dataSourcesOptions(this.props.data)} />
                  <Button icon='play' primary content='Run'/>
                </Form.Group>
                <Form.TextArea placeholder='Write your query here' />
              </Form>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              Results goes here
            </Grid.Column>
          </Grid.Row>
        </Container>
      </DocumentTitle>
    )
  }
}

const Query = gql`{dataSources { name }}`;

const PlaygroundPage = graphql(Query)(Playground);

export default PlaygroundPage;
