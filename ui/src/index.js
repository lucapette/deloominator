import React from 'react';
import ReactDOM from 'react-dom';
import {ApolloProvider} from 'react-apollo';

import AppConfig from './services/AppConfig';
import GraphqlClient from './services/GraphqlClient';
import {Provider} from 'react-redux';
import {createStore} from 'redux';

import appReducers from './reducers';

import Root from './Root';

let store = createStore(appReducers, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());

ReactDOM.render(
  <ApolloProvider client={GraphqlClient({port: AppConfig.port()})}>
    <Provider store={store}>
      <Root />
    </Provider>
  </ApolloProvider>,
  AppConfig.root(),
);
