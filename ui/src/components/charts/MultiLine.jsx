//@flow
import React, {Component} from 'react';
import VegaLite from 'react-vega-lite';

type Props = {
  width: number,
  values: Array<{}>,
  onNewView: Function,
};

class MultiLine extends Component<Props> {
  render() {
    const {width, values, onNewView} = this.props;
    const data = {
      values: values,
    };

    const [x, y, z] = Object.keys(this.props.values[0]);

    const spec = {
      description: 'A multiple lines chart with embedded data.',
      mark: 'line',
      width: width || 1000,
      encoding: {
        x: {type: 'temporal', axis: {shortTimeLabels: true}, field: x},
        y: {type: 'quantitative', field: z},
        color: {type: 'nominal', field: y},
      },
    };

    return <VegaLite spec={spec} data={data} onNewView={onNewView} />;
  }
}

export default MultiLine;
