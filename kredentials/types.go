package kredentials

import "github.com/SticketInya/kredentials/models"

type KubernetesConfigStore interface {
	Store(name string, config models.KubernetesConfig) error
	Load(name string) (*models.KubernetesConfig, error)
	LoadFromPath(path string) (*models.KubernetesConfig, error)
}

type KredentialStore interface {
	Store(kred *models.Kredential) error
	Load(name string) (*models.Kredential, error)
	List() ([]*models.Kredential, error)
	Delete(name string) error
}

type ArchiveStore interface {
	Store(path string, name string, kredentials []*models.Kredential) error
	Load(path string) ([]*models.Kredential, error)
}
