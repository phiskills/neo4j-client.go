package neo4j

import "strings"

type Operation interface {
	Then(Operation) Operation
	And(Operation) Operation
	Or(Operation) Operation
	XOr(Operation) Operation
	eval() (string, Records)
	isComposite() bool
	invert() Operation
}

type operation struct {
	value      string
	params     Records
	operations []Operation
	operators  []string
}

func (o operation) Then(operation Operation) Operation {
	return o.next(", ", operation)
}

func (o operation) And(condition Operation) Operation {
	return o.next(" AND ", condition)
}

func (o operation) Or(condition Operation) Operation {
	return o.next(" OR ", condition)
}

func (o operation) XOr(condition Operation) Operation {
	return o.next(" XOR ", condition)
}

func (o operation) eval() (string, Records) {
	operations := []string{o.value}
	for i, next := range o.operations {
		operation, params := next.eval()
		operations = append(operations, o.operators[i])
		if next.isComposite() {
			operation = "(" + operation + ")"
		}
		operations = append(operations, operation)
		o.params = o.params.Merge(params)
	}
	operation := strings.Join(operations, "")
	return operation, o.params
}

func (o operation) isComposite() bool {
	return len(o.operations) > 0 && o.value != "NOT"
}

func (o operation) invert() Operation {
	return operation{value: ""}.next("NOT ", o)
}

func (o operation) next(operator string, operation Operation) Operation {
	if operation == nil {
		return o
	}
	o.operations = append(o.operations, operation)
	o.operators = append(o.operators, operator)
	return o
}
