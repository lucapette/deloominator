//@flow
import React from 'react';
import {Container, Segment, List, Message} from 'semantic-ui-react';

import * as types from '../types';

type Props = {
  settings: types.Settings,
};

const Footer = (props: Props) => {
  const {settings} = props;
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
      {settings.isReadOnly && (
        <Message attached="bottom" color="orange" className="read-only">
          <p>Running in read-only mode. Learn more here!</p>
        </Message>
      )}
    </div>
  );
};

export default Footer;
