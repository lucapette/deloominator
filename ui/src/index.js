import React from 'react';
import ReactDOM from 'react-dom';
import {ApolloProvider} from 'react-apollo';

import AppConfig from './services/AppConfig';
import GraphqlClient from './services/GraphqlClient';

import App from './components/App';

ReactDOM.render(
  <ApolloProvider client={GraphqlClient({port: AppConfig.port()})}>
    <App />
  </ApolloProvider>,
  AppConfig.root(),
);
