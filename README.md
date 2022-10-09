# pginsert-bench

Benchmarking the performance of PostgreSQL's insert performance.

Using different methods:

- `COPY`.
- Batched with `INSERT INTO ... VALUES (...), (...), ...` syntax.
- `INSERT` with `INSERT INTO ... VALUES (...), (...), ...` syntax.
- `INSERT` with `INSERT INTO ... SELECT (UNNEST(...))` syntax.

# Benchmark Results

## 2 Columns

| rows    | cols | batch | copy | unnest | values |
|---------|------|-------|------|--------|--------|
| 400     | 2    | 7     | 4    | 3      | 3      |
| 2000    | 2    | 13    | 7    | 6      | 18     |
| 10000   | 2    | 54    | 35   | 24     | 49     |
| 50000   | 2    | 322   | 107  | 128    | x      |
| 250000  | 2    | 1711  | 440  | 599    | x      |
| 1000000 | 2    | 5302  | 1301 | 2280   | x      |

![](https://docs.google.com/spreadsheets/d/e/2PACX-1vS2UnZfxk7tGOpETnXislo9b_ruOkkhAmOcuvzwN5DQh5LKUkKZ9woQkLxx8ttBhWtR9HVXzAc3eXZB/pubchart?oid=1256092944&format=image)

## 3 Columns

| rows    | cols | batch | copy | unnest | values |
|---------|------|-------|------|--------|--------|
| 400     | 3    | 7     | 5    | 2      | 4      |
| 2000    | 3    | 11    | 8    | 6      | 12     |
| 10000   | 3    | 54    | 22   | 28     | 50     |
| 50000   | 3    | 254   | 86   | 124    | x      |
| 250000  | 3    | 1333  | 473  | 611    | x      |
| 1000000 | 3    | 6103  | 1470 | 2720   | x 		  |

![](https://docs.google.com/spreadsheets/d/e/2PACX-1vS2UnZfxk7tGOpETnXislo9b_ruOkkhAmOcuvzwN5DQh5LKUkKZ9woQkLxx8ttBhWtR9HVXzAc3eXZB/pubchart?oid=1680973688&format=image)

## 4 Columns

| rows    | cols | batch | copy | unnest | values |
|---------|------|-------|------|--------|--------|
| 400     | 4    | 6     | 3    | 8      | 5      |
| 2000    | 4    | 16    | 9    | 10     | 20     |
| 10000   | 4    | 67    | 31   | 35     | 91     |
| 50000   | 4    | 299   | 126  | 220    | x      |
| 250000  | 4    | 1912  | 554  | 781    | x      |
| 1000000 | 4    | 6882  | 1707 | 3816   | x      |

![](https://docs.google.com/spreadsheets/d/e/2PACX-1vS2UnZfxk7tGOpETnXislo9b_ruOkkhAmOcuvzwN5DQh5LKUkKZ9woQkLxx8ttBhWtR9HVXzAc3eXZB/pubchart?oid=689894188&format=image)
