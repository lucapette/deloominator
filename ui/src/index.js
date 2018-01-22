import React from 'react';
import ReactDOM from 'react-dom';
import {ApolloProvider} from 'react-apollo';

import AppConfig from './services/AppConfig';
import GraphqlClient from './services/GraphqlClient';

import Root from './Root';

ReactDOM.render(
  <ApolloProvider client={GraphqlClient({port: AppConfig.port()})}>
    <Root />
  </ApolloProvider>,
  AppConfig.root(),
);
