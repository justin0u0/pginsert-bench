package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jessevdk/go-flags"
)

func main() {
	var args benchArgs
	if _, err := flags.Parse(&args); err != nil {
		panic(err)
	}

	ctx := context.Background()

	pgc, err := pgx.Connect(ctx, args.URL)
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

	var writer writer
	if args.Save {
		fmt.Println("Saving benchmark results to", args.SavePath+".csv")

		w, err := newCSVWriter(args.SavePath + ".csv")
		if err != nil {
			panic(err)
		}
		defer w.close()

		header := []string{"rows", "cols"}
		for _, inserter := range inserters {
			header = append(header, inserter.name())
		}
		if err := w.write(header); err != nil {
			panic(err)
		}

		writer = w
	}

	for _, c := range configs {
		fmt.Println("config:", c)
		if err := c.setup(ctx, pgc); err != nil {
			panic(err)
		}

		fmt.Printf("Inserting %d rows into %s\n", c.rows, c.tableName)

		row := []string{strconv.Itoa(c.rows), strconv.Itoa(len(c.cols))}

		for _, inserter := range inserters {
			duration, err := inserter.insert(ctx, pgc, c)
			if err != nil {
				fmt.Printf("%s: %s\n", inserter.name(), err)
			} else {
				fmt.Printf("%s: %s\n", inserter.name(), duration)
			}

			row = append(row, strconv.FormatInt(int64(duration/time.Millisecond), 10))
		}

		if writer != nil {
			if err := writer.write(row); err != nil {
				panic(err)
			}
		}

		fmt.Println("")
	}
}
