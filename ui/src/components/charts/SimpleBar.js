//@flow
import React, {Component} from 'react';
import VegaLite from 'react-vega-lite';

type Props = {
  width: number,
  values: Array<{}>,
  onNewView: Function,
};

class SimpleBar extends Component<Props> {
  render() {
    const {values, onNewView, width} = this.props;
    const data = {
      values: values,
    };

    const [x, y] = Object.keys(values[0]);

    const spec = {
      description: 'A simple bar chart with embedded data.',
      mark: 'bar',
      width: width || 1000,
      encoding: {
        x: {type: 'ordinal', field: x},
        y: {type: 'quantitative', field: y},
      },
    };

    return <VegaLite spec={spec} data={data} onNewView={onNewView} />;
  }
}

export default SimpleBar;
