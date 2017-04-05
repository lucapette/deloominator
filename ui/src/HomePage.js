//@flow
import React, { Component } from 'react'
import DocumentTitle from 'react-document-title'

export default class HomePage extends Component {
  render() {
    return (
      <DocumentTitle title='Home'>
        <div id='welcome'>
          Welcome to deloominator!
        </div>
      </DocumentTitle>
    )
  }
}
