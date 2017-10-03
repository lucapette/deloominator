//@flow
import React, { Component } from "react";
import { gql, graphql } from "react-apollo";
import { Button, Container, Message, Loader, Segment, Divider } from "semantic-ui-react";

import Chart from "./Chart";
import Table from "./Table";

class QueryResultContainer extends Component {
  view: Object;

  componentWillUpdate(nextProps, nextState) {
    const { data: { loading, error, query }, handleQuerySuccess } = nextProps;
    if (handleQuerySuccess) {
      handleQuerySuccess(!(loading || error) && !(query != null && query.__typename == "queryError"));
    }
  }

  onNewView = (view: Object) => {
    this.view = view;
  };

  exportCSV = (e: MouseEvent) => {
    e.preventDefault();
    const { source, input } = this.props;
    console.log(this.props);

    fetch("http://localhost:3000/export/csv", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        Source: source,
        Query: input,
      }),
    })
      .then(response => response.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(blob);
        let a = document.createElement("a");
        a.href = url;
        a.download = "chart.csv";
        a.click();
      });
  };

  exportPNG = (e: MouseEvent) => {
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
    const { data: { loading, error, query }, handleQuerySuccess } = this.props;

    if (loading) {
      return (
        <Container>
          <Loader active inline="centered">
            Loading
          </Loader>
        </Container>
      );
    }

    if (error) {
      return <p>Error!</p>;
    }

    if (query.__typename == "queryError") {
      return (
        <Message negative>
          <Message.Header>There is a problem with your query.</Message.Header>
          <p>{query.message}</p>
        </Message>
      );
    }

    if (query.chartName !== "UnknownChart") {
      return (
        <Container>
          <Segment padded>
            <Button.Group basic floated="right">
              <Button onClick={this.exportPNG}>PNG</Button>
              <Button onClick={this.exportCSV}>CSV</Button>
            </Button.Group>
            <Chart name={query.chartName} columns={query.columns} rows={query.rows} onNewView={this.onNewView} />
            <Divider section horizontal>
              Raw data
            </Divider>
            <Table columns={query.columns} rows={query.rows} />
          </Segment>
        </Container>
      );
    }

    return (
      <Container>
        <Segment padded>
          <Table columns={query.columns} rows={query.rows} />
        </Segment>
      </Container>
    );
  }
}

const Query = gql`
  query Query($source: String!, $input: String!) {
    query(source: $source, input: $input) {
      ... on results {
        chartName
        columns {
          name
          type
        }
        rows {
          cells {
            value
          }
        }
      }
      ... on queryError {
        message
      }
    }
  }
`;

const QueryResult = graphql(Query, {
  options: ({ source, input }) => ({ variables: { source, input } }),
})(QueryResultContainer);

export default QueryResult;
