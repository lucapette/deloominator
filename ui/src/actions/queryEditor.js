export const resetQueryEditor = () => {
  return {
    type: 'RESET_QUERY_EDITOR',
  };
};

export const setQueryWasSuccessful = value => {
  return {
    type: 'SET_QUERY_WAS_SUCCESSFUL',
    queryWasSuccessful: value,
  };
};

export const setInputValue = (key, value) => {
  return {
    type: 'SET_INPUT_VALUE',
    key,
    value,
  };
};

export const setVariable = (key, value) => {
  return {
    type: 'SET_VARIABLE',
    key,
    value,
  };
};

export const setVariables = variables => {
  return {
    type: 'SET_VARIABLES',
    variables,
  };
};
