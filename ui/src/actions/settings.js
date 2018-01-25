//@flow

import type {Settings} from '../types';
import type {SetSettingsAction} from '../reducers/settings';

export const setSettings = (settings: Settings): SetSettingsAction => {
  return {
    type: 'SET_SETTINGS',
    settings,
  };
};
