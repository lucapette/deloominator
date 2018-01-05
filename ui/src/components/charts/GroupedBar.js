//@flow
import React, {Component} from 'react';
import VegaLite from 'react-vega-lite';

type Props = {
  values: Array<{}>,
  width: number,
  onNewView: Function,
};

class GroupedBar extends Component<Props> {
  render() {
    const {values, onNewView, width} = this.props;
    const data = {
      values: values,
    };

    const [x, y, z] = Object.keys(values[0]);

    const spec = {
      description: 'A multiple columns bar chart.',
      mark: 'bar',
      transform: [{calculate: `datum.${x}`, as: x}],
      encoding: {
        column: {type: 'ordinal', field: y},
        x: {type: 'nominal', scale: {rangeStep: 9}, axis: {title: ''}, field: x},
        y: {type: 'quantitative', axis: {grid: false, title: z}, field: z},
        color: {type: 'nominal', field: x},
      },
      config: {facet: {cell: {strokeWidth: 0}}},
    };

    return <VegaLite spec={spec} data={data} onNewView={onNewView} />;
  }
}

export default GroupedBar;
