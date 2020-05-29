package neo4j

import (
	"strings"
)

type Direction int

const (
	NoDirection Direction = iota
	FromOriginToDestination
	FromDestinationToOrigin
)

type Relationship struct {
	Id, Type  string
	Props     Records
	Direction Direction
}

type Path struct {
	Origin       *Node
	Relationship *Relationship
	Destination  Destination
}

type Destination interface {
	extends() (string, Records)
}

func (r *Relationship) Property(name string) Property {
	return property{
		name:  r.Id + "." + name,
		alias: r.Id + "_" + name,
	}
}

func (r *Relationship) Properties(names ...string) []Property {
	var props []Property
	for _, name := range names {
		props = append(props, r.Property(name))
	}
	return props
}

func (r *Relationship) eval() (string, Records) {
	if r == nil {
		return "--", Records{}
	}
	kind := ""
	if r.Type != "" {
		kind = ":" + r.Type
	}
	var props []string
	params := Records{}
	for _, key := range r.Props.Keys() {
		alias := r.Id + "_" + key
		prop := key + ": $" + alias
		props = append(props, prop)
		params[alias] = r.Props[key]
	}
	content := ""
	if len(props) > 0 {
		content = "{" + strings.Join(props, ", ") + "}"
	}
	arrow := ""
	if r.Id != "" || kind != "" || content != "" {
		arrow = "[" + r.Id + kind + content + "]"
	}
	switch r.Direction {
	case FromOriginToDestination:
		arrow = "-" + arrow + "->"
	case FromDestinationToOrigin:
		arrow = "<-" + arrow + "-"
	default:
		arrow = "-" + arrow + "-"
	}
	return arrow, params
}

func (p *Path) eval() (string, Records) {
	oOperation, oParams := "()", Records{}
	if p == nil {
		return oOperation, oParams
	}
	if p.Origin != nil {
		oOperation, oParams = p.Origin.eval()
	}
	if p.Destination == nil {
		return oOperation, oParams
	}
	rOperation, rParams := p.Relationship.eval()
	dOperation, dParams := p.Destination.extends()
	operation := oOperation + rOperation + dOperation
	params := oParams.Merge(rParams).Merge(dParams)
	return operation, params
}

func (p *Path) extends() (string, Records) {
	return p.eval()
}
