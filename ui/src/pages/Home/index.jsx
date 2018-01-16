//@flow

import React, {Component} from 'react';
import DocumentTitle from 'react-document-title';

type Props = {};

class Home extends Component<Props> {
  render() {
    return (
      <DocumentTitle title="Home">
        <div>Welcome to deloominator!</div>
      </DocumentTitle>
    );
  }
}

export default Home;
