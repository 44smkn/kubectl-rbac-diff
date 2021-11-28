package policy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

func WriteTableToFile(filename string, rows [][]string, header []string) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return fmt.Errorf("Failed to make directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	table := tablewriter.NewWriter(file)
	table.SetRowLine(true)
	table.SetHeader(header)
	for _, r := range rows {
		table.Append(r)
	}
	table.Render() // Send output

	return nil
}
