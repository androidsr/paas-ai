package access

import (
	"strings"

	"github.com/ledongthuc/pdf"
)

func getTextFromPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	totalPage := r.NumPage()
	var mergedTexts []string
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		var mergedSentence string

		texts := p.Content().Text
		for _, text := range texts {
			if text.Y == lastTextStyle.Y {
				mergedSentence += text.S
			} else {
				if mergedSentence != "" {
					mergedTexts = append(mergedTexts, mergedSentence)
				}
				lastTextStyle = text
				mergedSentence = text.S
			}
		}

		if mergedSentence != "" {
			mergedTexts = append(mergedTexts, mergedSentence)
		}
	}

	mergedText := strings.Join(mergedTexts, "\n")
	return mergedText, nil
}
