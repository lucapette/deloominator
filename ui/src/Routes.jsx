import React, {Fragment} from 'react';
import {Route} from 'react-router-dom';

import {idFromSlug} from './helpers/routing';

import Home from './pages/Home';
import Playground from './pages/Playground';
import Questions from './pages/Questions';
import Question from './pages/Questions/Question';
import QuestionEdit from './pages/Questions/QuestionEdit';

const Routes = props => {
  const {settings} = props;
  return (
    <Fragment>
      <Route path="/" exact component={Home} />
      <Route path="/playground" component={() => <Playground settings={settings} />} />
      <Route path="/questions" exact component={() => <Questions settings={settings} />} />
      <Route
        path="/questions/:id"
        exact
        component={({match}) => <Question id={idFromSlug(match.params.id)} settings={settings} />}
      />
      <Route
        path="/questions/:id/edit"
        exact
        component={({match}) => <QuestionEdit id={idFromSlug(match.params.id)} settings={settings} />}
      />
    </Fragment>
  );
};

export default Routes;
