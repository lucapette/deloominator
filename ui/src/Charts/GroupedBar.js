import React, { Component } from 'react';
import VegaLite from 'react-vega-lite';

const spec = {
  "description": "A multiple columns bar chart.",
  "mark": "bar",
  "encoding": {
    "column": {"type": "ordinal", "scale": {"padding": 4}, "axis": {"orient": "bottom", "axisWidth": 1, "offset": -8}},
    "x": {"type": "nominal", "scale": {"bandSize": 9}, "axis": null},
    "y": {"type": "quantitative", "axis": { "grid": false}},
    "color": {"type": "nominal", "scale": {"range": "category20"}},
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
      {"field": x, "expr": `datum.${x}`}
    ];

    return (<VegaLite spec={spec} data={data} />);
  }
}
