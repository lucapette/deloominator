# Charts

The scope of this document is to present the core functionality of
`deloominator`: the automatic chart display.

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

The data set has _2_ columns and they have the types _Number_ and _Text_, in
this order. `deloominator` would render a [Simple Bar](#simple-bar) in this case.

## Supported charts

Here is a list of all the charts `deloominator` detects automatically.

### Simple Bar

- *2* Columns
- Types:
  - `Number`
  - `Text`

### Simple Line

- **2** Columns
- Types:
  - `Number`
  - `Time`

### Multiple Bar

- **3** Columns
- Types:
  - `Number`
  - `Text`
  - `Text`

### Multiple Lines

- **3** Columns
- Types:
  - `Number`
  - `Time`
  - `Text`