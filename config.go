package main

import (
	"context"
	"fmt"

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

func (c *config) setup(ctx context.Context, pgc *pgx.Conn) error {
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
