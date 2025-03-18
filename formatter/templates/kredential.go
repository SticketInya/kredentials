package templates

import "fmt"

type KredentialNodeTemplate struct {
	Name     string
	Clusters string
	Contexts string
}

func (n KredentialNodeTemplate) String() string {
	return fmt.Sprintf("%s\t%s\t%s", n.Name, n.Clusters, n.Contexts)
}

type KredentialNodeListTemplate struct {
	Items []KredentialNodeTemplate
}

func (t KredentialNodeListTemplate) Headers() string {
	return "NAME\tCLUSTERS\tCONTEXTS"
}

func (t KredentialNodeListTemplate) Rows() []string {
	var rows []string
	for _, node := range t.Items {
		rows = append(rows, node.String())
	}
	return rows
}
