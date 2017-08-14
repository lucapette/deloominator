import React, { Component } from 'react';
import VegaLite from 'react-vega-lite';

const spec = {
  "description": "A multiple columns bar chart.",
  "mark": "bar",
  "encoding": {
    "column": {"type": "ordinal"},
    "x": { "type": "nominal", "scale": { "rangeStep": 9 }, "axis": {"title": ""}},
    "y": {"type": "quantitative", "axis": { "grid": false}},
    "color": {"type": "nominal"},
  },
  "config": {"facet": {"cell": {"strokeWidth": 0}}}
};

export default class GroupedBar extends Component {
  render() {
    const data = {
      "values": this.props.values
    };

    const [x, y, z] = Object.keys(this.props.values[0]);

    spec["encoding"]["x"]["field"] = x;
    spec["encoding"]["color"]["field"] = x;
    spec["encoding"]["y"]["field"] = z;
    spec["encoding"]["y"]["axis"]["title"] = z;
    spec["encoding"]["column"]["field"] = y;

    spec["transform"] = [
      { "calculate": `datum.${x}`, "as": x}
    ];

    return <VegaLite spec={spec} data={data} onNewView={this.props.onNewView} />;
  }
}
