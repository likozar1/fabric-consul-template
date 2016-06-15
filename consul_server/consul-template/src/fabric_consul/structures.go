package main

type ConsulNodes []struct {
	Node string `json:"Node"`
	Address string `json:"Address"`
}

type Nodes struct {
		nodes [] string
}

//nodes order by prefixes
type prefixedNodes map[string]*Nodes

//add node to map of prefixes
func (m *prefixedNodes) append(prefix string, node string) {
	if _, exists := (*m)[prefix]; !exists {
		(*m)[prefix] = &Nodes{}
	}
	(*m)[prefix].append(node)
}

func (al *Nodes) append(arLog string) {
	al.nodes = append(al.nodes, arLog)
}


