package access

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"

	"github.com/gonfva/docxlib"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/xuri/excelize/v2"
)

type Access struct {
	Path string
	file *os.File
	size int64
}

/**句号拆分文本**/
type FullStopTextSplitter struct {
}

func (m FullStopTextSplitter) SplitText(text string) ([]string, error) {
	return strings.Split(text, "。"), nil
}

/**换行拆分文本**/
type LineFeedTextSplitter struct {
}

func (m LineFeedTextSplitter) SplitText(text string) ([]string, error) {
	return strings.Split(text, "\n"), nil
}

func AutoLoad(filename string) ([]schema.Document, error) {
	var docs []schema.Document
	var err error
	switch filename[strings.LastIndex(filename, "."):] {
	case ".pdf":
		docs, err = NewPath(filename).ReadPDF()
	case ".html":
		docs, err = NewPath(filename).ReadHTML()
	case ".csv":
		docs, err = NewPath(filename).ReadCSV()
	case ".docx":
		docs, err = NewPath(filename).ReadWORD()
	case ".xlsx":
		docs, err = NewPath(filename).ReadXLSX()
	case ".md":
		docs, err = NewPath(filename).ReadMD()
	default:
		docs, err = NewPath(filename).ReadTEXT()
	}
	return docs, err
}

func NewPath(path string) *Access {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("读取文件失败")
		return nil
	}
	finfo, _ := f.Stat()
	return &Access{Path: path, file: f, size: finfo.Size()}
}

func NewFile(file *os.File, size int64) *Access {
	return &Access{file: file, size: size}
}

/** 读取PDF文件内容 **/
func (m *Access) ReadPDF() ([]schema.Document, error) {
	defer m.file.Close()
	p := documentloaders.NewPDF(m.file, m.size)
	docs, err := p.Load(context.Background())
	if err != nil {
		return nil, err
	}
	return docs, err
}

/** 读取CSV文件内容 **/
func (m *Access) ReadCSV() ([]schema.Document, error) {
	defer m.file.Close()
	loader := documentloaders.NewCSV(m.file)
	docs, err := loader.Load(context.Background())
	if err != nil {
		return nil, err
	}
	return docs, err
}

/** 读取HTML文件内容 **/
func (m *Access) ReadHTML() ([]schema.Document, error) {
	defer m.file.Close()
	loader := documentloaders.NewHTML(m.file)
	docs, err := loader.LoadAndSplit(context.Background(), LineFeedTextSplitter{})
	if err != nil {
		return nil, err
	}
	return docs, err
}

/** 读取TEXT文件内容 **/
func (m *Access) ReadTEXT() ([]schema.Document, error) {
	defer m.file.Close()
	loader := documentloaders.NewText(m.file)
	docs, err := loader.LoadAndSplit(context.Background(), FullStopTextSplitter{})
	if err != nil {
		return nil, err
	}
	return docs, err
}

/** 读取MD文件内容 **/
func (m *Access) ReadMD() ([]schema.Document, error) {
	defer m.file.Close()
	loader := documentloaders.NewText(m.file)
	docs, err := loader.LoadAndSplit(context.Background(), FullStopTextSplitter{})
	if err != nil {
		return nil, err
	}
	return docs, err
}

/** 读取WORD文件内容 **/
func (m *Access) ReadWORD() ([]schema.Document, error) {
	defer m.file.Close()
	doc, err := docxlib.Parse(m.file, m.size)
	if err != nil {
		panic(err)
	}
	result := make([]schema.Document, 0)
	for i, para := range doc.Paragraphs() {
		for _, child := range para.Children() {
			if child.Run != nil {
				result = append(result, schema.Document{PageContent: child.Run.Text.Text, Metadata: map[string]any{"page": i}})
			}
			if child.Link != nil {
				id := child.Link.ID
				text := child.Link.Run.InstrText
				link, err := doc.References(id)
				if err != nil {
					result = append(result, schema.Document{PageContent: text, Metadata: map[string]any{"page": i, "id": id}})
				} else {
					result = append(result, schema.Document{PageContent: text, Metadata: map[string]any{"page": i, "link": link}})
				}
			}
		}
	}
	return result, err
}

/** 读取Excel文件内容 **/
func (m *Access) ReadXLSX() ([]schema.Document, error) {
	defer m.file.Close()
	f, err := excelize.OpenReader(m.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	result := make([]schema.Document, 0)
	for i, name := range f.GetSheetList() {
		rows, err := f.GetRows(name)
		if err != nil {
			log.Println(err)
			continue
		}
		content := bytes.Buffer{}
		for _, row := range rows {
			content.WriteString(strings.Join(row, ","))
		}
		result = append(result, schema.Document{PageContent: content.String(), Metadata: map[string]any{"page": i, "sheetName": name}})
	}
	return result, err
}

func (doc *Access) NotionMD(filename string) ([]schema.Document, error) {
	notion := documentloaders.NewNotionDirectory(filename)
	return notion.Load()
}
