import React, {Fragment} from 'react';
import {Route} from 'react-router-dom';

import {idFromSlug} from './helpers/routing';

import Home from './pages/Home';
import Playground from './pages/Playground';
import Questions from './pages/Questions';
import Question from './pages/Questions/Question';
import QuestionEdit from './pages/Questions/QuestionEdit';

const Routes = () => {
  return (
    <Fragment>
      <Route path="/" exact component={Home} />
      <Route path="/playground" component={() => <Playground />} />
      <Route path="/questions" exact component={() => <Questions />} />
      <Route path="/questions/:id" exact component={({match}) => <Question id={idFromSlug(match.params.id)} />} />
      <Route
        path="/questions/:id/edit"
        exact
        component={({match}) => <QuestionEdit id={idFromSlug(match.params.id)} />}
      />
    </Fragment>
  );
};

export default Routes;
