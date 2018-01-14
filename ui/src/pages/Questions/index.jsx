//@flow
import React from 'react';

import Question from './Question';
import QuestionList from './QuestionList';

import type {Match} from 'react-router';

import * as types from '../../types';

type Props = {
  settings: types.Settings,
  match: Match,
};

const Questions = (props: Props) => {
  const {match, settings} = props;
  const question = match.params.question;
  if (question) {
    const [id] = question.split('-');

    return <Question id={id} />;
  }

  return settings.isReadOnly ? <p>No questions in read-only mode</p> : <QuestionList />;
};

export default Questions;
