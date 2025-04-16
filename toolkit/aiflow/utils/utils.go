package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func GetFlowParameter(input, output map[string]any, expression string) (string, error) {
	if expression == "" {
		return "", nil
	}

	data := map[string]any{}
	var key string
	if strings.HasPrefix(expression, "$input.") && input != nil {
		data = input
		key = strings.TrimPrefix(expression, "$input.")
	} else if strings.HasPrefix(expression, "$output.") && output != nil {
		data = output
		key = strings.TrimPrefix(key, "$output.")
	}

	if len(data) == 0 {
		return "", nil
	}
	// 获取嵌套值
	return getNestedValue(data, key)
}

// getNestedValue 从嵌套的 map 中获取值
func getNestedValue(data map[string]any, key string) (string, error) {
	if key == "" || data == nil {
		return "", nil
	}
	keys := strings.Split(key, ".")
	current := data

	for _, k := range keys {
		if v, ok := current[k]; ok {
			if next, ok := v.(map[string]any); ok {
				current = next
			} else {
				if v, ok := v.(string); ok {
					return v, nil
				} else {
					bs, err := json.Marshal(v)
					if err != nil {
						return fmt.Sprint(v), nil
					}
					return string(bs), nil
				}
			}
		} else {
			return "", fmt.Errorf("key not found: %s", k)
		}
	}
	return "", nil
}

func MatchExpression(content string) []string {
	regex := `(\$input\.|\$output\.)[a-zA-Z0-9_.]+`
	re := regexp.MustCompile(regex)
	matches := re.FindAllString(content, -1)
	return matches
}

func ReplaceExpression(content string, input, output map[string]any) string {
	exprs := MatchExpression(content)
	for _, expr := range exprs {
		value, err := GetFlowParameter(input, output, expr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		content = strings.ReplaceAll(content, expr, fmt.Sprintf("%v", value))
	}
	return content
}

func OutputPrint(emitter chan string, content string) {
	emitter <- content
	emitter <- " \n\n"
}

func OutputError(emitter chan string, nodeName, error string) {
	emitter <- "## " + nodeName + " 执行失败 \n\n"
	emitter <- error
	emitter <- " \n\n"
}

func InputPrint(emitter chan string, nodeName, content string) {
	emitter <- "## " + nodeName + " \n\n"
	emitter <- content
	emitter <- " \n\n"
}
