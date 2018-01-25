const initalState = {
  queryWasSuccessful: false,
  dataSource: '',
  query: '',
  queryDraft: '',
  variables: [],
};

const setVariable = (state, action) => {
  const index = state.variables.findIndex(e => e['name'] == action.key);

  return {...state, variables: state.variables.map((item, i) => (index !== i ? item : {...item, value: action.value}))};
};

const queryEditor = (state = initalState, action) => {
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
