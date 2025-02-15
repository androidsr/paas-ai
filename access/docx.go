package access

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

func GetTextFromDocx(path string) (string, error) {
	// 打开zip文件
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", errors.New("无法打开文档")
	}
	defer r.Close()

	var doc Document
	for _, file := range r.File {
		// 找到包含文档内容的XML文件
		if strings.Contains(file.Name, "word/document.xml") {
			rc, err := file.Open()
			if err != nil {
				log.Fatalf("无法打开文件内容: %s", err)
				return "", errors.New("无法打开文件内容")
			}
			defer rc.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, rc)
			if err != nil {
				return "", errors.New("读取文件内容错误")
			}

			err = xml.Unmarshal(buf.Bytes(), &doc)
			if err != nil {
				return "", errors.New("XML解析错误")
			}
		}
	}
	content := bytes.Buffer{}
	// 打印文档段落内容
	for _, para := range doc.Body.Paragraphs {
		for _, text := range para.Texts {
			content.WriteString(text.Content)
		}
		content.WriteString("\n\n")
	}

	// 打印表格内容
	for _, table := range doc.Body.Tables {
		for _, row := range table.Rows {
			for _, cell := range row.Cells {
				for _, para := range cell.Paragraphs {
					for _, text := range para.Texts {
						fmt.Print(text.Content)
						content.WriteString(text.Content)
					}
					fmt.Print("\t") // 单元格间用制表符分隔
				}
			}
			content.WriteString("\n\n")
		}
	}
	return content.String(), nil
}

// 定义用于解析 XML 的结构体
type Document struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

type Body struct {
	XMLName    xml.Name    `xml:"body"`
	Paragraphs []Paragraph `xml:"p"`
	Tables     []Table     `xml:"tbl"`
}

type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Texts   []Text   `xml:"r>t"`
}

type Text struct {
	XMLName xml.Name `xml:"t"`
	Content string   `xml:",chardata"`
}

type Table struct {
	XMLName xml.Name   `xml:"tbl"`
	Rows    []TableRow `xml:"tr"`
}

type TableRow struct {
	XMLName xml.Name    `xml:"tr"`
	Cells   []TableCell `xml:"tc"`
}

type TableCell struct {
	XMLName    xml.Name    `xml:"tc"`
	Paragraphs []Paragraph `xml:"p"`
}
