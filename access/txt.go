package access

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/androidsr/sc-go/shttp"
	"github.com/androidsr/sc-go/sno"
	"github.com/tmc/langchaingo/schema"
)

func GetSupportedFileTypes() []string {
	return []string{".txt", ".md", ".csv", ".yaml", ".docx", ".pdf"}
}

func GetParsedTextFromUrl(url string, ext string) ([]schema.Document, error) {
	var path string
	var err error
	if !strings.HasPrefix(url, "http") {
		path = url
	} else {
		data, err := shttp.Get(url, shttp.HTML, nil)
		if err != nil {
			return nil, err
		}
		path = "/tmp/" + sno.GetString() + "web.txt"
		os.WriteFile(path, data, 0644)
		defer func() {
			err = os.Remove(path)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
			}
		}()
	}
	var res string
	if ext == ".txt" || ext == ".md" || ext == ".yaml" {
		res, err = getTextFromPlain(path)
	} else if ext == ".csv" {
		res, err = getTextFromCsv(path)
	} else if ext == ".docx" {
		res, err = GetTextFromDocx(path)
	} else if ext == ".pdf" {
		res, err = getTextFromPdf(path)
	} else if ext == ".xlsx" {
		res, err = getTextFromXlsx(path)
	} else if ext == ".pptx" {
		res, err = getTextFromPptx(path)
	} else {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	p, err := GetSplitProvider("Default")
	if err != nil {
		return nil, err
	}
	textSections, err := p.SplitText(res)
	if err != nil {
		return nil, err
	}
	result := make([]schema.Document, 0)
	var item schema.Document
	for i, v := range textSections {
		fmt.Println(i, v)
		item = schema.Document{}
		item.Metadata = map[string]any{"page": fmt.Sprint(i), "filename": filepath.Base(path)}
		item.PageContent = v
		result = append(result, item)
	}
	return result, nil
}
