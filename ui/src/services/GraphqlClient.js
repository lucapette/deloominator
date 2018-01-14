import {ApolloClient} from 'apollo-client';
import {createHttpLink} from 'apollo-link-http';
import {InMemoryCache, IntrospectionFragmentMatcher} from 'apollo-cache-inmemory';

const fragmentMatcher = new IntrospectionFragmentMatcher({
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

const GraphqlClient = ({port}) => {
  return new ApolloClient({
    link: createHttpLink({uri: `http://localhost:${port}/graphql`}),
    cache: new InMemoryCache({fragmentMatcher}),
  });
};

export default GraphqlClient;
