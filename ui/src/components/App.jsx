import React from 'react';
import gql from 'graphql-tag';
import {Container, Loader} from 'semantic-ui-react';
import {graphql} from 'react-apollo';
import {BrowserRouter as Router, Route} from 'react-router-dom';

import 'semantic-ui-css/semantic.min.css';
import 'semantic-ui-css/semantic.min.js';

import './app.css';

import NavMenu from './../layout/NavMenu';
import Footer from './../layout/Footer';

import Home from './../pages/Home';
import Playground from './../pages/Playground';
import Questions from './../pages/Questions';

const SettingsQuery = gql`
  {
    settings {
      isReadOnly
    }
  }
`;

const App = graphql(SettingsQuery)(({data: {loading, error, settings}}) => {
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
          <Route path="/" exact component={Home} />
          <Route path="/playground" component={() => <Playground settings={settings} />} />
          <Route path="/questions" exact component={routeProps => <Questions {...routeProps} settings={settings} />} />
          <Route
            path="/questions/:question"
            component={routeProps => <Questions {...routeProps} settings={settings} />}
          />
        </Container>

        <Footer settings={settings} />
      </div>
    </Router>
  );
});

export default App;
