//@flow
import React, {Component} from 'react';
import {
  ApolloClient,
  ApolloProvider,
  createNetworkInterface,
  graphql,
  gql,
  IntrospectionFragmentMatcher,
} from 'react-apollo';
import {} from 'react-apollo';
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

const fm = new IntrospectionFragmentMatcher({
  introspectionQueryResultData: {
    __schema: {
      types: [
        {
          kind: 'UNION',
          name: 'QueryResult',
          possibleTypes: [{name: 'queryError'}, {name: 'results'}],
        },
      ],
    },
  },
});

const networkInterface = createNetworkInterface({
  uri: 'http://localhost:3000/graphql',
});

const client = new ApolloClient({
  fragmentMatcher: fm,
  networkInterface: networkInterface,
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

const mountNode = document.getElementById('root');

ReactDOM.render(
  <ApolloProvider client={client}>
    <App />
  </ApolloProvider>,
  mountNode,
);
