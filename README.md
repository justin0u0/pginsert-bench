# pginsert-bench

Benchmarking the performance of PostgreSQL's insert performance.

Using different methods:

- `COPY`.
- Batched with `INSERT INTO ... VALUES (...), (...), ...` syntax.
- `INSERT` with `INSERT INTO ... VALUES (...), (...), ...` syntax.
- `INSERT` with `INSERT INTO ... SELECT (UNNEST(...))` syntax.
