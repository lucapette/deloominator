import React from "react";
import { Container, Segment, List, Divider } from "semantic-ui-react";

const Footer = () => {
  return (
    <Segment vertical style={{ margin: "5em 0em 0em", padding: "5em 0em" }}>
      <Divider section />
      <Container textAlign="center">
        <List horizontal divided link>
          <List.Item as="a" href="https://github.com/lucapette/deloominator/issues/new">
            Report a problem
          </List.Item>
        </List>
      </Container>
    </Segment>
  );
};

export default Footer;
