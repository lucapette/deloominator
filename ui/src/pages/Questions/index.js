import React from "react";

import Question from "./Question";
import QuestionList from "./QuestionList";

const Questions = ({ settings, match }) => {
  const question = match.params.question;
  if (question) {
    const [id] = question.split("-");

    return <Question id={id} />;
  }

  return settings.isReadOnly ? <p>No questions in read-only mode</p> : <QuestionList />;
};

export default Questions;
