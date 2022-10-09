package main

import (
	"encoding/csv"
	"os"
)

type writer interface {
	write([]string) error
	close() error
}

type csvWriter struct {
	file   *os.File
	writer *csv.Writer
}

var _ writer = (*csvWriter)(nil)

func newCSVWriter(path string) (*csvWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &csvWriter{
		file:   f,
		writer: csv.NewWriter(f),
	}, nil
}

func (w *csvWriter) write(row []string) error {
	return w.writer.Write(row)
}

func (w *csvWriter) close() error {
	w.writer.Flush()
	return w.file.Close()
}
