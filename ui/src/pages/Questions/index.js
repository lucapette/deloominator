import React from "react";
import { Container } from "semantic-ui-react";

import Question from "./Question";

const Questions = ({ match }) => {
  const id = match.params.id;
  if (!id) {
    return <Container>Welcome to Questions section</Container>;
  }

  return <Question id={id} />;
};

export default Questions;
