//@flow
import React, { Component } from 'react'
import DocumentTitle from 'react-document-title'

import { Container } from 'semantic-ui-react'

export default class Home extends Component {
  render() {
    return (
      <DocumentTitle title='Home'>
        <Container id='welcome'>
          Welcome to deloominator!
        </Container>
      </DocumentTitle>
    )
  }
}
