package templates

import (
	"fmt"

	"github.com/SticketInya/kredentials/internal/formatterutil"
	"github.com/SticketInya/kredentials/models"
)

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

func BuildKredentialNodeList(kreds []*models.Kredential) KredentialNodeListTemplate {
	var template KredentialNodeListTemplate

	for _, kred := range kreds {
		template.Items = append(template.Items, KredentialNodeTemplate{
			Name:     kred.Name,
			Clusters: formatterutil.JoinMapKeysGeneric(kred.Config.Clusters, formatterutil.DefaultListSeparator),
			Contexts: formatterutil.JoinMapKeysGeneric(kred.Config.Contexts, formatterutil.DefaultListSeparator),
		})
	}

	return template
}
