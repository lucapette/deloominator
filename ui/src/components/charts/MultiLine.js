import React, {Component} from 'react';
import VegaLite from 'react-vega-lite';

const spec = {
  description: 'A multiple lines chart with embedded data.',
  mark: 'line',
  encoding: {
    x: {type: 'temporal', axis: {shortTimeLabels: true}},
    y: {type: 'quantitative'},
    color: {type: 'nominal'},
  },
};

export default class SimpleLine extends Component {
  render() {
    const data = {
      values: this.props.values,
    };

    const [x, y, z] = Object.keys(this.props.values[0]);

    spec['encoding']['x']['field'] = x;
    spec['encoding']['y']['field'] = z;
    spec['encoding']['color']['field'] = y;

    spec['width'] = this.props.width || 1000;

    return <VegaLite spec={spec} data={data} onNewView={this.props.onNewView} />;
  }
}
