package access

import (
	"github.com/tealeg/xlsx"
)

func getTextFromXlsx(path string) (string, error) {
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return "", err
	}
	var result string
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text, err := cell.FormattedValue()
				if err != nil {
					return "", err
				}
				result += text + " "
			}
			result += "\n"
		}
	}
	return result, nil
}
