The scope of this document is to provide an overview of how to use
`deloominator`. As long as `deloominator` is not in beta (you can check the
current status
[here](https://github.com/lucapette/deloominator/milestones?direction=desc&sort=completeness&state=open)),
this document serves as a brief spec of the features we want to build.

Each of the following document sections represents a section of the
application. Each section has a corresponding menu item in the top navigation:

- [Playground](#playground)
  - [Editor](#editor)
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
- Raw data displayed `deloominator`
  renders only a table with raw data in case no chart was detected.

`deloominator` uses a simple [algorithm](/docs/charts.md#algorithm) to determine
which graph to show.

You can give a name to your experiment and then hit the `save` button which will
show you the saved `q&a`:

![save q&a](/docs/img/question.png)

## Editor

In `deloominator` editor, you can use _"query variables"_. A query variable is a
word surrounded by `{` and `}`. By default, the following variables are
predefined:

- `{start_date}`
- `{end_date}`
- `{today}`
- `{yesterday}`

Variables are evaluated by `deloominator` right before query execution.
Questions that contain variables come with additional UI elements to
interactively control the value of variables before query execution.

Here is an example for `{today}`:

![today](/docs/img/today.png)

Some query variables require user input (like `{start_date}`), those have a
sensible default to get you started when writing a query.

# Q&A

The `Q&A` section presents the list of the questions users saved in
`deloominator`. It looks like this:

![q&as](/docs/img/questions.png)

Users can search for existing `q&a`, view, edit, and duplicate existing ones.

## Viewing a Q&A

# Dashboards

# Reports
