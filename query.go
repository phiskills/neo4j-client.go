package neo4j

import (
	"fmt"
	"strings"
)

type Query interface {
	Custom(string, Records) Query
	Create(structure structure) Query
	Set(Data, Records) Query
	Delete(...string) Query
	Match(structure structure) Query
	Merge(structure structure) Query
	OnCreate() Query
	OnMatch() Query
	Optional() Query
	OrderBy(...Property) Query
	Desc() Query
	Limit(int) Query
	Skip(int) Query
	Where(Operation) Query
	With(...string) Query
	Return(...Property) Query
	eval() (string, Records)
	String() string
}

type query struct {
	operations []string
	params     Records
}

type structure interface {
	eval() (string, Records)
}

type invertible interface {
	invert() Operation
}

func Not(i invertible) Operation {
	return i.invert()
}

func (q query) String() string {
	query, params := q.eval()
	return fmt.Sprintf("[Query] %s\n[Params] %v\n", query, params)
}

func (q query) Custom(operation string, params Records) Query {
	return query{
		operations: append(q.operations, operation),
		params:     q.params.Merge(params),
	}
}

func (q query) Create(structure structure) Query {
	return q.primary("CREATE", structure)
}

func (q query) Set(data Data, props Records) Query {
	if data == nil || len(props) == 0 {
		return q
	}
	var operations []Operation
	for _, key := range props.Keys() {
		value := props[key]
		operations = append(operations, data.Property(key).IsEqual(value))
	}
	operation, params := chain(operations).eval()
	return q.Custom("SET "+operation, params)
}

func (q query) Delete(ids ...string) Query {
	operation := strings.Join(ids, ", ")
	return q.Custom("DELETE "+operation, Records{})
}

func (q query) Match(structure structure) Query {
	return q.primary("MATCH", structure)
}

func (q query) Merge(structure structure) Query {
	return q.primary("MERGE", structure)
}

func (q query) OnCreate() Query {
	return q.Custom("ON CREATE", Records{})
}

func (q query) OnMatch() Query {
	return q.Custom("ON MATCH", Records{})
}

func (q query) Optional() Query {
	return q.Custom("OPTIONAL", Records{})
}

func (q query) OrderBy(props ...Property) Query {
	if len(props) == 0 {
		return q
	}
	var operations []Operation
	for _, prop := range props {
		operations = append(operations, prop.Get())
	}
	operation, _ := chain(operations).eval()
	return q.Custom("ORDER BY "+operation, Records{})
}

func (q query) Desc() Query {
	return q.Custom("DESC", Records{})
}

func (q query) Limit(limit int) Query {
	operation := fmt.Sprintf("LIMIT %d", limit)
	return q.Custom(operation, Records{})
}

func (q query) Skip(limit int) Query {
	operation := fmt.Sprintf("SKIP %d", limit)
	return q.Custom(operation, Records{})
}

func (q query) Where(condition Operation) Query {
	operation, params := condition.eval()
	return q.Custom("WHERE "+operation, params)
}

func (q query) With(ids ...string) Query {
	operation := strings.Join(ids, ", ")
	return q.Custom("WITH "+operation, Records{})
}

func (q query) Return(props ...Property) Query {
	if len(props) == 0 {
		return q
	}
	var operations []Operation
	for _, prop := range props {
		operations = append(operations, prop.Get())
	}
	operation, _ := chain(operations).eval()
	return q.Custom("RETURN "+operation, Records{})
}

func (q query) eval() (string, Records) {
	operation := strings.Join(q.operations, " ")
	return operation, q.params
}

func (q query) primary(instruction string, structure structure) Query {
	value, params := structure.eval()
	operation := instruction + " " + value
	return q.Custom(operation, params)
}

func chain(operations []Operation) Operation {
	var first Operation
	for _, operation := range operations {
		current := operation
		if first == nil {
			first = current
			continue
		}
		first = first.Then(current)
	}
	return first
}
