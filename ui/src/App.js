//@flow
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Menu, Grid } from 'semantic-ui-react';
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';

import {
  BrowserRouter as Router,
  Route,
  NavLink
} from 'react-router-dom';

import 'semantic-ui-css/semantic.min.css';
import 'semantic-ui-css/semantic.min.js';

import './app.css';

import HomePage from './HomePage';
import PlaygroundPage from './PlaygroundPage';

const networkInterface = createNetworkInterface({
  uri: 'http://localhost:3000/graphql'
});

const client = new ApolloClient({
  networkInterface: networkInterface
});

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <Menu>
            <Container>
              <NavLink exact to='/' className='item' activeClassName='active'>Home</NavLink>
              <NavLink to='/playground' className='item' activeClassName='active'>Playground</NavLink>
            </Container>
          </Menu>

          <Route exact path='/' component={HomePage}/>
          <Route path='/playground' component={PlaygroundPage}/>
        </div>
      </Router>
    )
  }
}

const mountNode = document.getElementById('root');

ReactDOM.render(
  <ApolloProvider client={client}>
    <App />
  </ApolloProvider>,
  mountNode
);
