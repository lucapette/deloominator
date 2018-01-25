//@flow

import type {Settings} from '../types';

type State = {
  +isReadOnly: boolean,
};

export type SetSettingsAction = {type: 'SET_SETTINGS', settings: Settings};

const settings = (state: State = {isReadOnly: false}, action: SetSettingsAction) => {
  switch (action.type) {
    case 'SET_SETTINGS':
      return {...state, ...action.settings};
    default:
      return state;
  }
};

export default settings;
