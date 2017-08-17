//@flow
import React, { Component } from "react";
import { ApolloClient, ApolloProvider, createNetworkInterface } from "react-apollo";
import ReactDOM from "react-dom";
import { Container, Menu, Grid } from "semantic-ui-react";

import { BrowserRouter as Router, Route, NavLink } from "react-router-dom";

import "semantic-ui-css/semantic.min.css";
import "semantic-ui-css/semantic.min.js";

import "./app.css";

import NavMenu from "./layout/NavMenu";
import Footer from "./layout/Footer";

import Home from "./pages/Home";
import Playground from "./pages/Playground";
import Questions from "./pages/Questions";

const networkInterface = createNetworkInterface({
  uri: "http://localhost:3000/graphql",
});

const client = new ApolloClient({
  networkInterface: networkInterface,
});

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <NavMenu />
          <Container>
            <Route exact path="/" component={Home} />
            <Route path="/playground" component={Playground} />
            <Route exact path="/questions" component={Questions} />
            <Route path="/questions/:question" component={Questions} />
          </Container>

          <Footer />
        </div>
      </Router>
    );
  }
}

const mountNode = document.getElementById("root");

ReactDOM.render(
  <ApolloProvider client={client}>
    <App />
  </ApolloProvider>,
  mountNode,
);
