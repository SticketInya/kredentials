package kredentials

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/internal/fileutil"
)

const (
	kredentialConfigDirPermissions = 0755
	DefaultConfigStorageDir        = "~/.kredentials/configs/"
)

type Kredential struct {
	Name   string
	Config *KubernetesConfig
}

func NewKredential(name string) *Kredential {
	return &Kredential{
		Name: name,
	}
}

func (k *Kredential) getConfigStorageDir() string {
	// TODO: Add env support
	return DefaultConfigStorageDir
}

func (k *Kredential) StoreKredential(path string) error {
	destDir, err := fileutil.ExpandPath(k.getConfigStorageDir())
	if err != nil {
		return fmt.Errorf("expandig config storage directory %w", err)
	}

	if err = os.MkdirAll(destDir, kredentialConfigDirPermissions); err != nil {
		return fmt.Errorf("creating config storage directory %w", err)
	}

	if err = k.WriteKubernetesConfig(destDir, k.Name); err != nil {
		return fmt.Errorf("writing '%s' to config storage directory %w", k.Name, err)
	}

	return nil
}
