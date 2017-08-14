//@flow
import { map, zipObject } from "lodash";
import React, { Component } from "react";

import SimpleBar from "./charts/SimpleBar";
import SimpleLine from "./charts/SimpleLine";
import GroupedBar from "./charts/GroupedBar";
import MultiLine from "./charts/MultiLine";

import * as Types from "../types";

const CHARTS = {
  SimpleBar: SimpleBar,
  SimpleLine: SimpleLine,
  GroupedBar: GroupedBar,
  MultiLine: MultiLine,
};

class Chart extends Component {
  view: Object;
  props: {
    name: string,
    rows: Array<Types.Row>,
    columns: Array<Types.Column>,
  };

  onNewView = (view: Object) => {
    this.view = view;
  };

  onClick = (e: MouseEvent) => {
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
    const { name, columns, rows } = this.props;

    const Handler = CHARTS[name];

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
export default Chart;
