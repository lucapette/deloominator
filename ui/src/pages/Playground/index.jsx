import React, {Component, Fragment} from 'react';
import DocumentTitle from 'react-document-title';
import {connect} from 'react-redux';
import {Grid} from 'semantic-ui-react';

import QueryResult from '../../components/QueryResult';
import QuestionForm from './QuestionForm';

import * as actions from '../../actions/queryEditor';

class PlaygroundContainer extends Component {
  componentWillUnmount() {
    this.props.resetQueryEditor();
  }
  render() {
    const {settings, dataSource, query, variables} = this.props;

    return (
      <DocumentTitle title="Playground">
        <Fragment>
          <Grid.Row>
            <Grid.Column>
              <QuestionForm saveEnabled={!settings.isReadOnly} />
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              {dataSource && query && <QueryResult dataSource={dataSource} query={query} variables={variables} />}
            </Grid.Column>
          </Grid.Row>
        </Fragment>
      </DocumentTitle>
    );
  }
}

const mapStateToProps = state => {
  return {
    ...state.queryEditor,
  };
};

const Playground = connect(mapStateToProps, {resetQueryEditor: actions.resetQueryEditor})(PlaygroundContainer);

export default Playground;
