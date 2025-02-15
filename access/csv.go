package access

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
)

func getTextFromCsv(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.LazyQuotes = true
	headers, err := r.Read()
	if err != nil {
		return "", err
	}

	var result []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		m := make(map[string]interface{})
		for i, header := range headers {
			header = strings.TrimSpace(header)
			value := strings.TrimSpace(record[i])

			if intValue, convErr := strconv.Atoi(value); convErr == nil {
				m[header] = intValue
			} else {
				m[header] = value
			}
		}

		jsonData, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		result = append(result, string(jsonData))
	}

	return strings.Join(result, "\n"), nil
}
