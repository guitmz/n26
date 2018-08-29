package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type tblWriter struct {
	*tablewriter.Table
}

func NewTableWriter() *tblWriter {
	return &tblWriter{tablewriter.NewWriter(os.Stdout)}
}

func (table *tblWriter) WriteData(header []string, data [][]string) error {
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
	return nil
}
