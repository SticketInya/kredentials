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

func getConfigStorageDir() string {
	// TODO: Add env support
	return DefaultConfigStorageDir
}

func (k *Kredential) StoreKredential(path string) error {
	destDir, err := fileutil.ExpandPath(getConfigStorageDir())
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
		kred := NewKredential(filepath.Base(destPath))
		err := kred.readKredential(destPath)
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

		kred := NewKredential(file.Name())
		targetFilePath := filepath.Join(destPath, file.Name())
		err := kred.readKredential(targetFilePath)
		if err != nil {
			// TODO: properly handle error
			fmt.Printf("cannot read file '%s' %s", targetFilePath, err)
			continue
		}

		kreds = append(kreds, kred)
	}

	return kreds, nil
}

func (k *Kredential) readKredential(path string) error {
	err := k.ReadKubernetesConfig(path)
	if err != nil {
		return err
	}
	return nil
}

func CheckKredentialInStorage(name string) bool {
	configDir, err := fileutil.ExpandPath(getConfigStorageDir())
	if err != nil {
		return false
	}

	return fileutil.CheckFileExists(filepath.Join(configDir, name))
}

func RetrieveKredentialFromStorage(name string) (*Kredential, error) {
	if !CheckKredentialInStorage(name) {
		return nil, fmt.Errorf("config '%s' does not exists", name)
	}

	storageDir, err := fileutil.ExpandPath(getConfigStorageDir())
	if err != nil {
		return nil, fmt.Errorf("expanding config storage directory %w", err)
	}

	kreds, err := ReadKredentials(filepath.Join(storageDir, name))
	if err != nil {
		return nil, err
	}
	if len(kreds) == 0 {
		return nil, fmt.Errorf("retrieved 0 kredentials")
	}

	// TODO: check if we should handle more than 1 results
	return kreds[0], nil
}

func DeleteKredentialFromStorage(name string) error {
	storageDir, err := fileutil.ExpandPath(getConfigStorageDir())
	if err != nil {
		return fmt.Errorf("expanding config storage directory %w", err)
	}

	targetPath := filepath.Join(storageDir, name)
	if err = os.Remove(targetPath); err != nil {
		return fmt.Errorf("deleting kredential at '%s' %w", targetPath, err)
	}

	return nil
}
