import React, {Component} from 'react';
import gql from 'graphql-tag';
import {graphql} from 'react-apollo';
import {Form} from 'semantic-ui-react';

import {sortBy} from 'lodash';

class DataSourcesContainer extends Component<Props> {
  dataSourcesOptions = dataSources => {
    return sortBy(dataSources || [], ['name'], ['asc']).map(s => ({
      name: s.name,
      text: s.name,
      value: s.name,
    }));
  };

  render() {
    const {data: {loading, error, dataSources}, handleDataSourcesChange, dataSource} = this.props;

    if (error) {
      return <p>Error!</p>;
    }

    return (
      <Form.Dropdown
        loading={loading}
        placeholder="Data Source"
        search
        selection
        onChange={handleDataSourcesChange}
        options={this.dataSourcesOptions(dataSources)}
        value={dataSource}
      />
    );
  }
}

const Query = gql`
  {
    dataSources {
      name
    }
  }
`;

const DataSources = graphql(Query)(DataSourcesContainer);

export default DataSources;
