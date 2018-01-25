import React, {Component} from 'react';
import {connect} from 'react-redux';

import gql from 'graphql-tag';
import {Container, Loader} from 'semantic-ui-react';
import {graphql} from 'react-apollo';
import {BrowserRouter as Router} from 'react-router-dom';

import 'semantic-ui-css/semantic.min.css';
import 'semantic-ui-css/semantic.min.js';

import './app.css';

import Routes from './Routes';

import * as actions from './actions/settings';

import NavMenu from './layout/NavMenu';
import Footer from './layout/Footer';

const SettingsQuery = gql`
  {
    settings {
      isReadOnly
    }
  }
`;

class RootContainer extends Component {
  componentWillUpdate(nextProps) {
    const {data: {settings}, setSettings} = nextProps;
    if (settings) {
      setSettings(settings);
    }
  }

  render() {
    const {data: {loading, error}} = this.props;

    if (loading) {
      return <Loader active />;
    }
    if (error) {
      return <p>error!</p>;
    }
    return (
      <Router>
        <div className="page">
          <NavMenu />
          <Container className="content">
            <Routes />
          </Container>
          <Footer />
        </div>
      </Router>
    );
  }
}

const Root = connect(null, {setSettings: actions.setSettings})(graphql(SettingsQuery)(RootContainer));
export default Root;
