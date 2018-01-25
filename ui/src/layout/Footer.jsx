//@flow
import React from 'react';
import {connect} from 'react-redux';
import {Container, Segment, List, Message} from 'semantic-ui-react';

import * as types from '../types';

type Props = {
  settings: types.Settings,
};

const FooterContainer = (props: Props) => {
  const {settings: {isReadOnly}} = props;
  return (
    <div className="footer">
      <Segment attached>
        <Container>
          <List horizontal divided link>
            <List.Item as="a" href="https://github.com/lucapette/deloominator/issues/new">
              Report a problem
            </List.Item>
          </List>
        </Container>
      </Segment>
      {isReadOnly && (
        <Message attached="bottom" color="orange" className="read-only">
          <p>Running in read-only mode. Learn more here!</p>
        </Message>
      )}
    </div>
  );
};

const mapStateToProps = state => {
  return {
    settings: state.settings,
  };
};

const Footer = connect(mapStateToProps)(FooterContainer);

export default Footer;
