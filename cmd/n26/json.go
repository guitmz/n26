package main

import (
	"encoding/json"
	"fmt"

	"github.com/guitmz/n26"
)

type jsonWriter struct{}

func (w jsonWriter) WriteTransactions(t *n26.Transactions) error {
	formatted, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(formatted))
	return nil
}
