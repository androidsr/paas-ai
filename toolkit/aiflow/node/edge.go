package node

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type Edge struct {
	SourceNodeId string `json:"sourceNodeId"`
	TargetNodeId string `json:"targetNodeId"`
	LoopType     string `json:"loopType"`
	LoopVariable string `json:"loopVariable"`
	Condition    string `json:"condition"`
}

func NewEdge(sourceNodeId, targetNodeId, loopType, loopVariable, condition string) *Edge {
	return &Edge{SourceNodeId: sourceNodeId, TargetNodeId: targetNodeId, LoopType: loopType, LoopVariable: loopVariable, Condition: condition}
}

func (e *Edge) EvaluateCondition(context map[string]interface{}) bool {
	if e.Condition == "" {
		return true
	}
	ok, err := evaluateExpression(e.Condition, context)
	if err != nil {
		return false
	}
	return ok
}

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
