import React from 'react';
import gql from 'graphql-tag';
import {Container, Loader} from 'semantic-ui-react';
import {graphql} from 'react-apollo';
import {BrowserRouter as Router} from 'react-router-dom';

import 'semantic-ui-css/semantic.min.css';
import 'semantic-ui-css/semantic.min.js';

import './app.css';

import Routes from './Routes';

import NavMenu from './layout/NavMenu';
import Footer from './layout/Footer';

const SettingsQuery = gql`
  {
    settings {
      isReadOnly
    }
  }
`;

const Root = graphql(SettingsQuery)(({data: {loading, error, settings}}) => {
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
          <Routes settings={settings} />
        </Container>
        <Footer settings={settings} />
      </div>
    </Router>
  );
});

export default Root;
