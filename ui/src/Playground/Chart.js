import React, { Component } from "react";

import SimpleBar from "../Charts/SimpleBar";
import SimpleLine from "../Charts/SimpleLine";
import GroupedBar from "../Charts/GroupedBar";
import MultiLine from "../Charts/MultiLine";

import { map, zipObject } from "lodash";

const CHARTS = {
  SimpleBar: SimpleBar,
  SimpleLine: SimpleLine,
  GroupedBar: GroupedBar,
  MultiLine: MultiLine,
};

export default class Chart extends Component {
  onNewView = view => {
    this.view = view;
  };

  onClick = e => {
    e.preventDefault();
    this.view
      .toImageURL("png")
      .then(url => {
        let link = document.createElement("a");
        link.setAttribute("href", url);
        link.setAttribute("target", "_blank");
        link.setAttribute("download", "chart.png");
        link.dispatchEvent(new MouseEvent("click"));
      })
      .catch(function(error) {
        /* error handling */
      });
  };

  render() {
    const { chartName, columns, rows } = this.props.data;

    const Handler = CHARTS[chartName];

    const columnNames = map(columns, "name");

    const values = rows.map(row => {
      const cells = map(row.cells, "value");

      return zipObject(columnNames, cells);
    });

    return (
      <div>
        <Handler values={values} onNewView={this.onNewView} />
        <a href="" onClick={this.onClick}>
          Download as PNG
        </a>
      </div>
    );
  }
}
