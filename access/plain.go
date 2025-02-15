package access

import (
	"bytes"
	"os"
	"strings"
)

func getTextFromPlain(path string) (string, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	res := string(bs)
	datas := strings.Split(res, "。")
	dataResult := bytes.Buffer{}
	for _, v := range datas {
		dataResult.WriteString(v)
		dataResult.WriteString("。")
		dataResult.WriteString("\n\n")
	}
	return dataResult.String(), nil
}
