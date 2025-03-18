package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/SticketInya/kredentials/internal/fileutil"
	"github.com/SticketInya/kredentials/internal/kubernetes_util"
	"github.com/SticketInya/kredentials/models"
)

type KubernetesConfigStore interface {
	Store(name string, config models.KubernetesConfig) error
	Load(name string) (*models.KubernetesConfig, error)
	LoadFromPath(path string) (*models.KubernetesConfig, error)
}

type FileKubernetesConfigStore struct {
	storageDirectory      string
	storageDirPermissions os.FileMode
}

func NewFileKubernetesConfigStore(storageDirectory string, storageDirPermissions os.FileMode) *FileKubernetesConfigStore {
	return &FileKubernetesConfigStore{
		storageDirectory:      storageDirectory,
		storageDirPermissions: storageDirPermissions,
	}
}

func (s *FileKubernetesConfigStore) getAndExpandStorageDirectory() (string, error) {
	destDir, err := fileutil.ExpandPath(s.storageDirectory)
	if err != nil {
		return "", fmt.Errorf("expanding kubernetes directory: %w", err)
	}

	return destDir, nil
}

func (s *FileKubernetesConfigStore) Store(name string, config models.KubernetesConfig) error {
	storageDir, err := s.getAndExpandStorageDirectory()
	if err != nil {
		return err
	}

	// ensure the directory exists
	if err = os.MkdirAll(storageDir, s.storageDirPermissions); err != nil {
		return fmt.Errorf("creating kubernetes directory: %w", err)
	}

	filename := filepath.Join(storageDir, name)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file '%s': %w", filename, err)
	}
	defer file.Close()

	_, err = kubernetes_util.WriteKubernetesConfig(file, config)
	if err != nil {
		return fmt.Errorf("writing '%s' file %w", filename, err)
	}
	return nil
}

func (s *FileKubernetesConfigStore) loadFileAndParse(filename string) (*models.KubernetesConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file '%s': %w", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file content '%s': %w", filename, err)
	}

	return kubernetes_util.ReadKubernetesConfig(data)
}

func (s *FileKubernetesConfigStore) Load(name string) (*models.KubernetesConfig, error) {
	storageDir, err := s.getAndExpandStorageDirectory()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(storageDir, name)
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("file '%s' does not exists", filename)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("file '%s' is a directory", filename)
	}

	return s.loadFileAndParse(filename)
}

func (s *FileKubernetesConfigStore) LoadFromPath(path string) (*models.KubernetesConfig, error) {
	targetPath, err := fileutil.ExpandPath(path)
	if err != nil {
		return nil, fmt.Errorf("expanding target path '%s':%w", path, err)
	}

	if !fileutil.CheckFileExists(targetPath) {
		return nil, fmt.Errorf("file '%s' does not exists", targetPath)
	}

	return s.loadFileAndParse(targetPath)
}
