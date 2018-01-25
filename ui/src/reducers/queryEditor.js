//@flow

import type {Variable} from '../types';

type State = {
  +queryWasSuccessful: boolean,
  +dataSource: string,
  +query: string,
  +queryDraft: string,
  +variables: Array<Variable>,
};

export type SetQuerySuccessAction = {type: 'SET_QUERY_WAS_SUCCESSFUL', queryWasSuccessful: boolean};
export type ResetAction = {type: 'RESET_QUERY_EDITOR'};
export type SetKeyValueAction = {type: 'SET_INPUT_VALUE' | 'SET_VARIABLE', key: string, value: string};
export type SetVariablesAction = {type: 'SET_VARIABLES', variables: Array<Variable>};

type Action = SetQuerySuccessAction | ResetAction | SetKeyValueAction | SetVariablesAction;

const initalState = {
  queryWasSuccessful: false,
  dataSource: '',
  query: '',
  queryDraft: '',
  variables: [],
};

const setVariable = (state: State, action: SetKeyValueAction): State => {
  const index = state.variables.findIndex(e => e['name'] == action.key);

  return {
    ...state,
    variables: state.variables.map((item: Variable, i) => (index !== i ? item : {...item, value: action.value})),
  };
};

const queryEditor = (state: State = initalState, action: Action): State => {
  switch (action.type) {
    case 'RESET_QUERY_EDITOR':
      return initalState;
    case 'SET_INPUT_VALUE':
      return {...state, [action.key]: action.value};
    case 'SET_QUERY_WAS_SUCCESSFUL':
      return {...state, queryWasSuccessful: action.queryWasSuccessful};
    case 'SET_VARIABLES':
      return {...state, variables: action.variables};
    case 'SET_VARIABLE':
      return setVariable(state, action);
    default:
      return state;
  }
};

export default queryEditor;
