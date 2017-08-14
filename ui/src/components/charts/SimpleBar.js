import React, { Component } from 'react';
import VegaLite from 'react-vega-lite';

const spec = {
  "description": "A simple bar chart with embedded data.",
  "mark": "bar",
  "encoding": {
    "x": {"type": "ordinal"},
    "y": {"type": "quantitative"}
  }
};

export default class SimpleBar extends Component {
  render() {
    const data = {
      "values": this.props.values
    };

    const [x, y] = Object.keys(this.props.values[0]);

    spec["encoding"]["x"]["field"] = x;
    spec["encoding"]["y"]["field"] = y;

    spec["width"] = this.props.width || 1000;

    return <VegaLite spec={spec} data={data} onNewView={this.props.onNewView} />;
  }
}
