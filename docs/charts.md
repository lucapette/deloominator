# Charts

The scope of this document is to present the core functionality of
`deloominator`: the automatic chart display.

## How it works

`deloominator detects chart based on the following characteristics of the
data set:

- The number of columns returned
- The types of the columns and their order

For example, let's consider the following data set:

| country | total |
|---------|-------|
| A       | 10    |
| B       | 42    |
| C       | 24    |

It has *2* columns and the sequence of types is:

- `Text` - `Number`

In such a case, `deloominator` renders a [Simple Bar](#simple-bar).

## Supported charts

Here is a list of the charts `deloominator` automatically detects.

### Simple Bar

*2* columns with one of the following type sequences:

- `Text`   - `Number`
- `Number` - `Number`

![simple bar](/img/simple-bar.png)

### Simple Line

**2** columns with one of the following type sequences:

- `Number` - `Time`

![simple line](/img/simple-line.png)

### Grouped Bar

**3** columns with one of the following type sequences:

- `Text` - `Text` - `Number`

![grouped bar](/img/grouped-bar.png)

### Multi Line

**3** columns with one of the following type sequences:

- `Time` - `Text` - `Number`

![multi-line](/img/multi-line.png)
