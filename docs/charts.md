# Charts

The scope of this document is to present the core functionality of
`deluminator`: the automatic chart display.

## Algorithm

The algorithm uses a convention based on the following characteristics of the
returned dataset:

- The number of columns returned
- The types of the columns and their order

For example, let's say a data set returns the following:

| total | country |
|-------|---------|
| 10    | A       |
| 42    | B       |
| 24    | C       |

The data set has _2_ columns and they have the types _number_ and _string_, in
this order. `deluminator` would render a [pie chart](#pie-chart) in this case.

## Supported charts

This is a list of all the charts `deluminator` supports.

### Pie chart

- *2* Columns
- Columns types:
  - `number`
  - `string`
