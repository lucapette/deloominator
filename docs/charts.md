# Charts

The scope of this document is to present the core functionality of
`deloominator`: the automatic chart display.

## Algorithm

The algorithm uses a convention based on the following characteristics of the
returned dataset:

- The number of columns returned
- The types of the columns and their order

For example, let's consider the following dataset:

| country | total |
|---------|-------|
| A       | 10    |
| B       | 42    |
| C       | 24    |

It has *2* columns and the sequece of types is:

- `Text` - `Number`

In such a case, `deloominator` renders a [Simple Bar](#simple-bar).

## Supported charts

Here is a list of the charts `deloominator` automatically detects.

### Simple Bar

*2* columns with one of the following type sequences:

- `Text`   - `Number`
- `Number` - `Number`

### Simple Line

**2** columns with one of the following type sequences:

- `Number` - `Time`

### Grouped Bar

**3** columns with one of the following type sequences:

- `Text` - `Text` - `Number`

### Multi Line

**3** columns with one of the following type sequences:

- `Time` - `Text` - `Number`