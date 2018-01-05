//@flow
import React, {Component} from 'react';
import VegaLite from 'react-vega-lite';

type Props = {
  width: number,
  values: Array<{}>,
  onNewView: Function,
};

class SimpleLine extends Component<Props> {
  render() {
    const {values, onNewView, width} = this.props;
    const data = {
      values: values,
    };

    const [x, y] = Object.keys(values[0]);
    const spec = {
      description: 'A simple line chart with embedded data.',
      mark: 'line',
      width: width || 1000,
      encoding: {
        x: {type: 'temporal', axis: {shortTimeLabels: true}, field: x},
        y: {type: 'quantitative', field: y},
      },
    };

    return <VegaLite spec={spec} data={data} onNewView={onNewView} />;
  }
}

export default SimpleLine;
