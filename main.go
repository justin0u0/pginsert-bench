package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

type config struct {
	tableName string
	rows      int
	cols      []string
}

func (c *config) String() string {
	return fmt.Sprintf("tableName: %s, rows: %d, cols: %s", c.tableName, c.rows, c.cols)
}

func (c *config) setup(ctx context.Context, pgc *pgx.Conn, cols int) error {
	tableCols := ""
	for _, col := range c.cols {
		tableCols += fmt.Sprintf(", %s BIGINT NOT NULL", col)
	}

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id BIGSERIAL PRIMARY KEY%s)", c.tableName, tableCols)
	if _, err := pgc.Exec(ctx, sql); err != nil {
		return err
	}

	if _, err := pgc.Exec(ctx, "TRUNCATE TABLE "+c.tableName); err != nil {
		return err
	}

	return nil
}

func benchBatch(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
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

func benchCopy(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
	start := time.Now()

	rows := make([][]any, c.rows)
	for i := 0; i < c.rows; i++ {
		rows[i] = []interface{}{i, i}
	}

	if _, err := pgc.CopyFrom(ctx, pgx.Identifier{c.tableName}, c.cols, pgx.CopyFromRows(rows)); err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func benchInsert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
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

func benchUnnestInsert(ctx context.Context, pgc *pgx.Conn, c *config) (time.Duration, error) {
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

func main() {
	ctx := context.Background()
	url := "postgres://postgres:nYcjZh9pHohei4XJA97lOWBG@localhost:5432/postgres?sslmode=disable"

	pgc, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}
	defer pgc.Close(ctx)

	configs := []*config{
		{
			tableName: "test_2_cols",
			rows:      400,
			cols:      []string{"col1", "col2"},
		},
		{
			tableName: "test_2_cols",
			rows:      2000,
			cols:      []string{"col1", "col2"},
		},
		{
			tableName: "test_2_cols",
			rows:      10000,
			cols:      []string{"col1", "col2"},
		},
		{
			tableName: "test_2_cols",
			rows:      50000,
			cols:      []string{"col1", "col2"},
		},
		{
			tableName: "test_2_cols",
			rows:      250000,
			cols:      []string{"col1", "col2"},
		},
		{
			tableName: "test_2_cols",
			rows:      1000000,
			cols:      []string{"col1", "col2"},
		},
	}

	for _, c := range configs {
		fmt.Println("config:", c)
		if err := c.setup(ctx, pgc, len(c.cols)); err != nil {
			panic(err)
		}

		fmt.Printf("Inserting %d rows into %s\n", c.rows, c.tableName)

		{
			d, err := benchInsert(ctx, pgc, c)
			if err != nil {
				fmt.Printf("Insert: %s\n", err.Error())
			} else {
				fmt.Printf("Insert: %s\n", d)
			}
		}

		{
			d, err := benchBatch(ctx, pgc, c)
			if err != nil {
				fmt.Printf("Batch: %s\n", err.Error())
			} else {
				fmt.Printf("Batch: %s\n", d)
			}
		}

		{
			d, err := benchCopy(ctx, pgc, c)
			if err != nil {
				fmt.Printf("Copy: %s\n", err.Error())
			} else {
				fmt.Printf("Copy: %s\n", d)
			}
		}

		{
			d, err := benchUnnestInsert(ctx, pgc, c)
			if err != nil {
				fmt.Printf("Unnest Insert: %s\n", err.Error())
			} else {
				fmt.Printf("Unnest Insert: %s\n", d)
			}
		}

		fmt.Println("")
	}
}
