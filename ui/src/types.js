//@flow
export type Row = {cells: Array<{value: string}>};

export type Column = {name: string};

export type Variable = {name: string, value: string, isControllable: boolean};

export type Settings = {isReadOnly: boolean};
