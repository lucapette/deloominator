//@flow
import React, { Component } from 'react'
import DocumentTitle from 'react-document-title'

import { Button, Container, Form, Grid } from 'semantic-ui-react'

export default class PlaygroundPage extends Component {
  constructor() {
    super();
    this.dataSources = [{text: 'stuff', value: 'stuff'}]
  }

  render() {
    return (
      <DocumentTitle title='Playground'>
        <Container>
          <Grid.Row>
            <Grid.Column>
              <Form>
                <Form.Group>
                  <Form.Dropdown placeholder='Data Source' search selection options={this.dataSources} />
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
