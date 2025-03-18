package kredentials

import (
	"fmt"
	"os"
	"path/filepath"

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

func ReadKredentials(path string) ([]*Kredential, error) {
	var kreds []*Kredential

	destPath, err := fileutil.ExpandPath(path)
	if err != nil {
		return kreds, fmt.Errorf("expanding target path '%s' %w", path, err)
	}

	pathInfo, err := os.Stat(destPath)
	if err != nil {
		return kreds, fmt.Errorf("file or directory '%s' does not exists", destPath)
	}

	if !pathInfo.IsDir() {
		kred, err := readKredential(destPath)
		if err != nil {
			return kreds, fmt.Errorf("failed to read kredential %w", err)
		}
		return append(kreds, kred), nil
	}

	files, err := os.ReadDir(destPath)
	if err != nil {
		return kreds, fmt.Errorf("cannot read directory '%s'", destPath)
	}

	for _, file := range files {
		if file.IsDir() {
			// we currently do not support nested directories
			continue
		}

		targetFilePath := filepath.Join(destPath, file.Name())
		kred, err := readKredential(targetFilePath)
		if err != nil {
			// TODO: properly handle error
			fmt.Printf("cannot read file '%s' %s", targetFilePath, err)
			continue
		}

		kreds = append(kreds, kred)
	}

	return kreds, nil
}

func readKredential(path string) (*Kredential, error) {
	name := filepath.Base(path)

	kred := NewKredential(name)
	err := kred.ReadKubernetesConfig(path)
	if err != nil {
		return nil, err
	}
	return kred, nil
}
