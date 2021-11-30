package policy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/olekukonko/tablewriter"
)

var ReadFile = func(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func WriteTableToFile(filename string, rows [][]string) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return fmt.Errorf("Failed to make directory: %w", err)
	}

	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetRowLine(true)
	for _, r := range rows {
		table.Append(defaultHeader)
		table.Rich(defaultHeader, []tablewriter.Colors{{tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}, {tablewriter.Bold}})
		table.Append(r)
	}
	table.Render() // Send output

	regex := regexp.MustCompile("(\n\\| APIGROUP)")
	contents := regex.ReplaceAllString(buf.String(), "\n\n${1}")

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	fmt.Fprintln(file, contents)
	defer file.Close()

	return nil
}
