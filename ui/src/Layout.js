//@flow
import React, { Component } from 'react'
import DocumentTitle from 'react-document-title'

export default class Layout extends Component {
  render() {
    return (
      <DocumentTitle title="deloominator">
        <div id="container">
          It works!
        </div>
      </DocumentTitle>
    )
  }
}
