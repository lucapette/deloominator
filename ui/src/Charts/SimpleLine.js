import React, { Component } from 'react';
import VegaLite from 'react-vega-lite';

const spec = {
  "description": "A simple line chart with embedded data.",
  "mark": "line",
  "encoding": {
    "x": {"type": "temporal", "axis": {"shortTimeLabels": true}},
    "y": {"type": "quantitative"}
  }
};

export default class SimpleLine extends Component {
  render() {
    const data = {
      "values": this.props.values
    };

    // not sure this is a safe assumption
    const [y, x] = Object.keys(this.props.values[0]);

    spec["encoding"]["x"]["field"] = x;
    spec["encoding"]["y"]["field"] = y;

    spec["width"] = this.props.width || 1000;

    return (<VegaLite spec={spec} data={data} />);
  }
}