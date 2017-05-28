# Charts

The scope of this document is to present the core functionality of
`deloominator`: the automatic chart display.

- [How it works](#how-it-works)
- [Supported charts](#supported-charts)
  - [Two-column charts](#two-column-charts)
    -[Simple Bar](#simple-bar)
    -[Simple Line](#simple-line)
    -[Simple Bar](#simple-bar)
  - [Three-column charts](#three-column-charts)
    -[Grouped Bar](#grouped-bar)
    -[Multi Line](#Multi-line)

# How it works

`deloominator` detects chart based on the following characteristics of the data
set:

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

# Supported charts

Here is a list of the charts `deloominator` automatically detects. They're
organized by the number of columns of the data set. Each charts comes with an
example query that is based on the [Sakila Sample
Database](https://dev.mysql.com/doc/sakila/en/).

## Two-Column charts

### Simple Bar

One of the following type sequences:

- `Text`   - `Number`
- `Number` - `Number`

Example:

```sql
SELECT rating, count(*) total FROM film GROUP BY 1
```

![simple bar](img/simple-bar.png)

### Simple Line

Type sequence:

- `Time` - `Number`

Example:

```sql
SELECT last_update , count(*) total FROM film group by 1
```

![simple line](img/simple-line.png)

## Three-column charts

### Grouped Bar

Type sequences:

- `Text` - `Text` - `Number`

Example:

```sql
SELECT rating, c.name, count(*) FROM film f
INNER JOIN film_category fc ON fc.film_id = f.film_id
INNER JOIN category c ON fc.category_id = c.category_id
GROUP BY 1, 2
```

![grouped bar](img/grouped-bar.png)

### Multi Line

Type sequence:

- `Time` - `Text` - `Number`

Example:

```sql
SELECT f.last_update, c.name, count(*) FROM film f
INNER JOIN film_category fc ON fc.film_id = f.film_id
INNER JOIN category c ON fc.category_id = c.category_id
WHERE c.name IN ('Sci-Fi', 'Drama', 'Travel', 'Documentary')
GROUP BY 1, 2
```

![multi-line](img/multi-line.png)
