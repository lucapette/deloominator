//@flow

import type {SetQuerySuccessAction, ResetAction, SetKeyValueAction, SetVariablesAction} from '../reducers/queryEditor';

import type {Variable} from '../types';

export const resetQueryEditor = (): ResetAction => {
  return {
    type: 'RESET_QUERY_EDITOR',
  };
};

export const setQueryWasSuccessful = (value: boolean): SetQuerySuccessAction => {
  return {
    type: 'SET_QUERY_WAS_SUCCESSFUL',
    queryWasSuccessful: value,
  };
};

export const setInputValue = (key: string, value: string): SetKeyValueAction => {
  return {
    type: 'SET_INPUT_VALUE',
    key,
    value,
  };
};

export const setVariable = (key: string, value: string): SetKeyValueAction => {
  return {
    type: 'SET_VARIABLE',
    key,
    value,
  };
};

export const setVariables = (variables: Array<Variable>): SetVariablesAction => {
  return {
    type: 'SET_VARIABLES',
    variables,
  };
};
