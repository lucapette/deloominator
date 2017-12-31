import React, {Component} from 'react';
import gql from 'graphql-tag';
import {ApolloClient} from 'apollo-client';
import {createHttpLink} from 'apollo-link-http';
import {InMemoryCache} from 'apollo-cache-inmemory';
import {ApolloProvider, graphql} from 'react-apollo';
import ReactDOM from 'react-dom';
import {Container, Menu, Grid, Loader} from 'semantic-ui-react';

import {BrowserRouter as Router, Route, NavLink} from 'react-router-dom';

import 'semantic-ui-css/semantic.min.css';
import 'semantic-ui-css/semantic.min.js';

import './app.css';

import NavMenu from './layout/NavMenu';
import Footer from './layout/Footer';

import Home from './pages/Home';
import Playground from './pages/Playground';
import Questions from './pages/Questions';

const client = new ApolloClient({
  link: createHttpLink({uri: 'http://localhost:3000/graphql'}),
  cache: new InMemoryCache(),
});

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

ReactDOM.render(
  <ApolloProvider client={client}>
    <App />
  </ApolloProvider>,
  document.getElementById('root'),
);
