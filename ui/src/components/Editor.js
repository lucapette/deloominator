import React, {Component} from 'react';
import CodeMirror from 'react-codemirror';

import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/sql/sql';

const options = {
  mode: 'text/x-sql',
  lineNumbers: true,
};

class Editor extends Component {
  render() {
    const {code, onChange} = this.props;
    return <CodeMirror className="thirteen wide field" options={options} value={code} onChange={onChange} />;
  }
}

export default Editor;
