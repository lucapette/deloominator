import React from "react";
import { Container, Segment, List, Message, Divider } from "semantic-ui-react";

const Footer = ({ settings }) => {
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
