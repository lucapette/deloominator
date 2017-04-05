//@flow
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Menu } from 'semantic-ui-react'

import {
  BrowserRouter as Router,
  Route,
  NavLink
} from 'react-router-dom'

import 'semantic-ui-css/components/site.min.css';
import 'semantic-ui-css/components/menu.min.css';
import './app.css'

import HomePage from './HomePage'
import PlaygroundPage from './PlaygroundPage'

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <Menu>
            <NavLink exact to='/' className='item' activeClassName='active'>Home</NavLink>
            <NavLink to='/playground' className='item' activeClassName='active'>Playground</NavLink>
          </Menu>

          <Route exact path='/' component={HomePage}/>
          <Route path='/playground' component={PlaygroundPage}/>
        </div>
      </Router>
    )
  }
}

const mountNode = document.getElementById('root');

ReactDOM.render(<App />, mountNode);
