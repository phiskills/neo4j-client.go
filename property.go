package neo4j

type Data interface {
	Property(name string) Property
	Properties(names ...string) []Property
}

type Property interface {
	Get() Operation
	IsEqual(interface{}) Operation
	IsNotEqual(interface{}) Operation
	LessThan(interface{}) Operation
	LessEqual(interface{}) Operation
	GreaterThan(interface{}) Operation
	GreaterEqual(interface{}) Operation
	StartsWith(string) Operation
	EndsWith(string) Operation
	Contains(interface{}) Operation
	In([]interface{}) Operation
	Matches(string) Operation
	IsNull() Operation
	IsNotNull() Operation
	invert() Operation
}

type property struct {
	name  string
	alias string
}

func (p property) Get() Operation {
	return operation{value: p.name}
}

func (p property) IsEqual(value interface{}) Operation {
	return p.apply("=", value)
}

func (p property) IsNotEqual(value interface{}) Operation {
	return p.apply("<>", value)
}

func (p property) LessThan(value interface{}) Operation {
	return p.apply("<", value)
}

func (p property) LessEqual(value interface{}) Operation {
	return p.apply("<=", value)
}

func (p property) GreaterThan(value interface{}) Operation {
	return p.apply(">", value)
}

func (p property) GreaterEqual(value interface{}) Operation {
	return p.apply(">=", value)
}

func (p property) StartsWith(value string) Operation {
	return p.apply("STARTS WITH", value)
}

func (p property) EndsWith(value string) Operation {
	return p.apply("ENDS WITH", value)
}

func (p property) Contains(value interface{}) Operation {
	return p.apply("CONTAINS", value)
}

func (p property) In(values []interface{}) Operation {
	return p.apply("IN", values)
}

func (p property) Matches(regex string) Operation {
	return p.apply("=~", regex)
}

func (p property) IsNull() Operation {
	return operation{value: p.name + " IS NULL"}
}

func (p property) IsNotNull() Operation {
	return operation{value: p.name + " IS NOT NULL"}
}

func (p property) apply(operator string, value interface{}) Operation {
	return operation{
		value:  p.name + " " + operator + " $" + p.alias,
		params: Records{p.alias: value},
	}
}

func (p property) invert() Operation {
	return operation{value: "NOT " + p.name}
}
