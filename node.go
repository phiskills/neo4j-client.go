package neo4j

import "strings"

type Node struct {
	Id     string
	Labels []string
	Props  Records
}

func (n *Node) Property(name string) Property {
	return property{
		name:  n.Id + "." + name,
		alias: n.Id + "_" + name,
	}
}

func (n *Node) Properties(names ...string) []Property {
	var props []Property
	for _, name := range names {
		props = append(props, n.Property(name))
	}
	return props
}

func (n *Node) eval() (string, Records) {
	if n == nil {
		return "()", Records{}
	}
	kind := ""
	for _, label := range n.Labels {
		kind += ":" + label
	}
	var props []string
	params := Records{}
	for _, key := range n.Props.Keys() {
		alias := n.Id + "_" + key
		prop := key + ": $" + alias
		props = append(props, prop)
		params[alias] = n.Props[key]
	}
	content := ""
	if len(props) > 0 {
		content = "{" + strings.Join(props, ", ") + "}"
	}
	node := "(" + n.Id + kind + content + ")"
	return node, params
}

func (n *Node) extends() (string, Records) {
	return n.eval()
}
