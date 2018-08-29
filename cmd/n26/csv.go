package main

import (
	"encoding/csv"
	"io"
)

type csvWriter struct {
	*csv.Writer
}

func NewCsvWriter(target io.Writer) (*csvWriter, error) {
	writer := csv.NewWriter(target)
	return &csvWriter{writer}, nil
}

func (w *csvWriter) WriteData(header []string, data [][]string) error {
	if err := w.Write(header); err != nil {
		return err
	}
	if err := w.WriteAll(data); err != nil {
		return err
	}
	return nil
}
