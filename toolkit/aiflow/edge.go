package aiflow

type Edge struct {
	SourceNodeId string
	TargetNodeId string
	Condition    string
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
