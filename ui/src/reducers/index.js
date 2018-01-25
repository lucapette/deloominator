import {combineReducers} from 'redux';

import queryEditor from './queryEditor';
import settings from './settings';

const appReducers = combineReducers({queryEditor, settings});

export default appReducers;
