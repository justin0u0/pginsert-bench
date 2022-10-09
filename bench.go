package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

type inserter interface {
	name() string
	insert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error)
}

// INSERT INTO table (cols1, cols2, ...) VALUES (...); INSERT INTO table (cols1, cols2, ...) VALUES (...); ...
type batchInserter struct{}

var _ inserter = (*batchInserter)(nil)

func (*batchInserter) name() string {
	return "batch"
}

func (*batchInserter) insert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
	start := time.Now()

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", c.tableName, strings.Join(c.cols, ","), preparePGXParameters(1, 1, len(c.cols)))

	batch := &pgx.Batch{}
	for i := 0; i < c.rows; i++ {
		args := make([]any, 0, len(c.cols))
		for j := 0; j < len(c.cols); j++ {
			args = append(args, i)
		}
		batch.Queue(sql, args...)
	}

	br := pgc.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < c.rows; i++ {
		if _, err := br.Exec(); err != nil {
			return 0, err
		}
	}

	return time.Since(start), nil
}

// COPY table (cols1, cols2, ...) FROM ...;
type copyInserter struct{}

var _ inserter = (*copyInserter)(nil)

func (*copyInserter) name() string {
	return "copy"
}

func (*copyInserter) insert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
	start := time.Now()

	rows := make([][]any, c.rows)
	for i := 0; i < c.rows; i++ {
		for j := 0; j < len(c.cols); j++ {
			rows[i] = append(rows[i], i)
		}
	}

	if _, err := pgc.CopyFrom(ctx, pgx.Identifier{c.tableName}, c.cols, pgx.CopyFromRows(rows)); err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

// INSERT INTO table (cols1, cols2, ...) VALUES (...), (...), ...;
type valuesInserter struct{}

var _ inserter = (*valuesInserter)(nil)

func (*valuesInserter) name() string {
	return "values"
}

func (*valuesInserter) insert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
	start := time.Now()

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", c.tableName, strings.Join(c.cols, ","), preparePGXParameters(1, c.rows, len(c.cols)))
	args := make([]any, 0, c.rows*len(c.cols))
	for i := 0; i < c.rows; i++ {
		for j := 0; j < len(c.cols); j++ {
			args = append(args, i)
		}
	}

	if _, err := pgc.Exec(ctx, sql, args...); err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

// INSERT INTO table (cols1, cols2, ...) SELECT * FROM UNNEST(cols1::BIGINT[], cols2::BIGINT[], ...);
type unnestInserter struct{}

var _ inserter = (*unnestInserter)(nil)

func (*unnestInserter) name() string {
	return "unnest"
}

func (*unnestInserter) insert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
	start := time.Now()

	unnestCols := ""
	for i := range c.cols {
		unnestCols += fmt.Sprintf("$%d::BIGINT[], ", i+1)
	}
	unnestCols = unnestCols[:len(unnestCols)-2]

	sql := fmt.Sprintf("INSERT INTO %s (%s) SELECT * FROM UNNEST(%s)", c.tableName, strings.Join(c.cols, ","), unnestCols)
	values := make([][]any, len(c.cols))
	for i := 0; i < len(c.cols); i++ {
		values[i] = make([]any, c.rows)
	}
	for i := 0; i < c.rows; i++ {
		for j := 0; j < len(c.cols); j++ {
			values[j][i] = i
		}
	}

	args := make([]any, 0, len(c.cols))
	for _, v := range values {
		args = append(args, v)
	}

	if _, err := pgc.Exec(ctx, sql, args...); err != nil {
		return 0, err
	}

	return time.Since(start), nil
}
