//@flow
import React, { Component } from 'react'
import DocumentTitle from 'react-document-title'

export default class HomePage extends Component {
  render() {
    return (
      <DocumentTitle title='Home'>
        <div>
          Welcome to deloominator!
        </div>
      </DocumentTitle>
    )
  }
}
