The scope of this document is to provide an overview of how to use
`deloominator`. As long as `deloominator` is not in beta (you can check the
current status
[here](https://github.com/lucapette/deloominator/milestones?direction=desc&sort=completeness&state=open)),
this document serves as a brief spec of the features we want to build.

Each of the following document sections represents a section of the
application. Each section has a corresponding menu item in the top navigation:

- [Playground](#playground)
- [Q&A](#q&a)
- [Dashboards](#dashboards)
- [Reports](#reports)

# Playground

The playground is the place where users explore their data sets. The interface
presents an SQL editor, a data source selector and a call to action that runs
queries. It looks like this:

![playground](/docs/img/playground.png)

Running a query in the playground has four possible outcomes:

- Error
  There is a problem with the query.
- No data displayed
  The query is fine, but there is no data to return.
- A chart displayed
  `deloominator` chose to display the data as a chart. Please refer to [this
  document](/docs/charts.md) for further information.
- Raw data displayed
  `deloominator` renders only a table with raw data in case no chart was detected.

`deloominator` uses a simple [algorithm](/docs/charts.md#algorithm) to
determine which graph to show.

Users can save their queries via the `save` button. As the goal of each
visualisation is to answer a question, `deloominator` calls a visualisation
`Q&A` and asks to provide the following information in order to save a
visualisation:

- a `title` that uniquely describes the `Q&A`.
- a `description` that will render markdown.

# Q&A

The `Q&A` section presents the list of the questions users saved in
`deloominator`. It looks like this:

![q&a mock](/docs/img/q-and-a.png)

Users can search for existing `q&a`, edit and duplicate existing ones.

# Dashboards

# Reports
