package aiflow

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

func evaluateExpression(expression string, context map[string]interface{}) (bool, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return false, fmt.Errorf("表达式语法错误: %v", err)
	}
	result, err := expr.Evaluate(context)
	if err != nil {
		return false, fmt.Errorf("表达式执行失败: %v", err)
	}
	booleanResult, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("表达式未返回布尔值")
	}
	return booleanResult, nil
}
