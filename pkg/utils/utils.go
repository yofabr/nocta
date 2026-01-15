package utils

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"
)

func ExportToJSON(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func ExportToCSV(headers []string, records [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return err
	}

	return writer.WriteAll(records)
}

func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

func NormalizeProtocol(protocol string) string {
	return strings.ToLower(strings.TrimSpace(protocol))
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
