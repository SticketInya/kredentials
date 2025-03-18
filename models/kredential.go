package models

type Kredential struct {
	Name   string
	Config *KubernetesConfig
}

func NewKredential(name string, config *KubernetesConfig) *Kredential {
	return &Kredential{
		Name:   name,
		Config: config,
	}
}
