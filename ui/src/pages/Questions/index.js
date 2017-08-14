import React from "react";
import { Container } from "semantic-ui-react";

import Question from "./Question";
import QuestionList from "./QuestionList";

const Questions = ({ match }) => {
  const question = match.params.question;
  if (!question) {
    return <QuestionList />;
  }

  const [id] = question.split("-");

  return <Question id={id} />;
};

export default Questions;
