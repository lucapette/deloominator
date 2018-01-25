import React, {Component} from 'react';
import {Button, Form, Modal} from 'semantic-ui-react';
import {connect} from 'react-redux';
import gql from 'graphql-tag';
import {withRouter} from 'react-router';
import {graphql} from 'react-apollo';

import {urlFor} from '../../helpers/routing';
import * as actions from '../../actions/queryEditor';

class SaveModalContainer extends Component {
  handleInputChange = event => {
    const target = event.target;
    const name = target.name;
    const value = target.value;

    this.props.setInputValue(name, value);
  };

  handleSave = e => {
    e.preventDefault();

    const {query, dataSource, history, mutate, variables, title, description} = this.props;

    mutate({
      refetchQueries: ['questions'],
      variables: {
        title,
        description,
        query,
        dataSource,
        variables,
      },
    }).then(({data: {question}}) => {
      const questionPath = urlFor(question, ['id', 'title']);
      history.push(`/questions/${questionPath}`);
    });
  };

  render() {
    const {queryWasSuccessful} = this.props;

    return (
      <Modal trigger={<Button icon="save" primary content="Save" disabled={!queryWasSuccessful} />}>
        <Modal.Header content="Save question" />
        <Modal.Content>
          <Form>
            <Form.Group widths="equal">
              <Form.Input
                fluid
                required
                name="title"
                label="Title"
                placeholder="Untitled question"
                onChange={this.handleInputChange}
              />
            </Form.Group>
            <Form.Group widths="equal">
              <Form.TextArea name="description" label="Description" onChange={this.handleInputChange} />
            </Form.Group>
          </Form>
          <Modal.Actions>
            <Button primary icon="save" content="Save" onClick={this.handleSave} />
            <Button content="Cancel" />
          </Modal.Actions>
        </Modal.Content>
      </Modal>
    );
  }
}

const mapStateToProps = state => {
  return {...state.queryEditor};
};

const SaveQuestion = gql`
  mutation SaveQuestion(
    $title: String!
    $query: String!
    $dataSource: String!
    $description: String
    $variables: [InputVariable]
  ) {
    question(title: $title, query: $query, dataSource: $dataSource, description: $description, variables: $variables) {
      id
      title
    }
  }
`;

const SaveModal = connect(mapStateToProps, {
  setInputValue: actions.setInputValue,
})(withRouter(graphql(SaveQuestion)(SaveModalContainer)));

export default SaveModal;
