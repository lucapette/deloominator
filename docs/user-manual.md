The scope of this document is to provide an overview of how to use
`deloominator`. As long as `deloominator` is not in beta (you can check the
current status
[here](https://github.com/lucapette/deloominator/milestones?direction=desc&sort=completeness&state=open)),
this document serves as a brief spec of the features we want to build.

Each of the following sections of the document represents a section of the
application. Each section has a corresponding menu item in the top navigation:

- [Playground](#playground)
- [Q&A](#q&a)
- [Dashboards](#dashboards)
- [Reports](#reports)

# Playground

The playground is the place where users explore their data sets. The interface
presents a SQL editor, a data source selector and  call to action that runs
queries. It looks like this:

![playground mockup](/docs/img/playground.png)

Running a query in the playground has four possible outcomes:

- Error
  There is a problem with the query.
- No data display
  The query is fine but there are no data to return.
- Chart display
  `deloominator` chose to display the data as a chart. Please refer to [this
  document](/docs/charts.md) for further information.
- Raw data display
  `deloominator` will use this option only if there is no way to display the
  data as a chart.

`deloominator` uses a simple [algorithm](/docs/charts.md#algorithm) to
determine which graph to show.

# Q&A

# Dashboards

# Reports
