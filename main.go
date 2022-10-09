package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func main() {
	ctx := context.Background()
	url := "postgres://postgres:nYcjZh9pHohei4XJA97lOWBG@localhost:5432/postgres?sslmode=disable"

	pgc, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}
	defer pgc.Close(ctx)

	configs := []*config{
		{tableName: "test_2_cols", rows: 400, cols: []string{"col1", "col2"}},
		{tableName: "test_2_cols", rows: 2000, cols: []string{"col1", "col2"}},
		{tableName: "test_2_cols", rows: 10000, cols: []string{"col1", "col2"}},
		{tableName: "test_2_cols", rows: 50000, cols: []string{"col1", "col2"}},
		{tableName: "test_2_cols", rows: 250000, cols: []string{"col1", "col2"}},
		{tableName: "test_2_cols", rows: 1000000, cols: []string{"col1", "col2"}},
		{tableName: "test_3_cols", rows: 400, cols: []string{"col1", "col2", "col3"}},
		{tableName: "test_3_cols", rows: 2000, cols: []string{"col1", "col2", "col3"}},
		{tableName: "test_3_cols", rows: 10000, cols: []string{"col1", "col2", "col3"}},
		{tableName: "test_3_cols", rows: 50000, cols: []string{"col1", "col2", "col3"}},
		{tableName: "test_3_cols", rows: 250000, cols: []string{"col1", "col2", "col3"}},
		{tableName: "test_3_cols", rows: 1000000, cols: []string{"col1", "col2", "col3"}},
		{tableName: "test_4_cols", rows: 400, cols: []string{"col1", "col2", "col3", "col4"}},
		{tableName: "test_4_cols", rows: 2000, cols: []string{"col1", "col2", "col3", "col4"}},
		{tableName: "test_4_cols", rows: 10000, cols: []string{"col1", "col2", "col3", "col4"}},
		{tableName: "test_4_cols", rows: 50000, cols: []string{"col1", "col2", "col3", "col4"}},
		{tableName: "test_4_cols", rows: 250000, cols: []string{"col1", "col2", "col3", "col4"}},
		{tableName: "test_4_cols", rows: 1000000, cols: []string{"col1", "col2", "col3", "col4"}},
	}

	inserters := []inserter{
		(*batchInserter)(nil),
		(*copyInserter)(nil),
		(*valuesInserter)(nil),
		(*unnestInserter)(nil),
	}

	for _, c := range configs {
		fmt.Println("config:", c)
		if err := c.setup(ctx, pgc); err != nil {
			panic(err)
		}

		fmt.Printf("Inserting %d rows into %s\n", c.rows, c.tableName)

		for _, inserter := range inserters {
			duration, err := inserter.insert(ctx, pgc, c)
			if err != nil {
				fmt.Printf("%s: %s\n", inserter.name(), err)
			} else {
				fmt.Printf("%s: %s\n", inserter.name(), duration)
			}
		}

		fmt.Println("")
	}
}
