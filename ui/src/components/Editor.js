//@flow
import React, {Component} from 'react';
import CodeMirror from 'react-codemirror';

import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/sql/sql';

const options = {
  mode: 'text/x-sql',
  lineNumbers: true,
};

type Props = {
  code: string,
  onChange: Function,
};

class Editor extends Component<Props> {
  render() {
    const {code, onChange} = this.props;
    return <CodeMirror className="thirteen wide field" options={options} value={code} onChange={onChange} />;
  }
}

export default Editor;
