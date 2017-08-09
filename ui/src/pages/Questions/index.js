import React from "react";
import { Container } from "semantic-ui-react";

import Question from "./Question";

const Questions = ({ match }) => {
  const question = match.params.question;
  if (!question) {
    return <Container>Welcome to Questions section</Container>;
  }

  const [id] = question.split("-");

  return <Question id={id} />;
};

export default Questions;
