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

![playground](/docs/img/playground.png)

Running a query in the playground has four possible outcomes:

- Error
  There is a problem with the query.
- No data display
  The query is fine but there are no data to return.
- Chart display
  `deloominator` chose to display the data as a chart. Please refer to [this
  document](/docs/charts.md) for further information.
- Raw data display
  `deloominator` renders only a table with raw data in case no chart was detected.

`deloominator` uses a simple [algorithm](/docs/charts.md#algorithm) to
determine which graph to show.

Hitting the `save` button shows the following screen:

![save q&a mock](/docs/img/save-q-and-a.png)

As the goal of each visualization or query is to answer a question,
`deloominator` saves the information as a `Q&A`. The following information is
required to save a `Q&A`:

- a `title` that uniquely describes the `Q&A`
- a `description` that will render as markdown

# Q&A

The `Q&A` section presents the list of the questions users saved in
`deloominator`. It looks like this:

![q&a mock](/docs/img/list-of-q-and-a.png)

Users can search for existing `q&a`, edit and duplicate existing ones.

# Dashboards

# Reports
