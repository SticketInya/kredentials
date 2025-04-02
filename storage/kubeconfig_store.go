package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/SticketInya/kredentials/internal/fileutil"
	"github.com/SticketInya/kredentials/internal/kubernetesutil"
	"github.com/SticketInya/kredentials/models"
)

var (
	ErrKubernetesConfigInvalid    = errors.New("invalid kubernetes config")
	ErrKubernetesConfigNotFound   = errors.New("kubernetes config not found")
	ErrKubernetesConfigAccess     = errors.New("kubernetes config file access error")
	ErrKubernetesConfigIsDir      = errors.New("expected file but got directory")
	ErrKubernetesConfigCannotOpen = errors.New("cannot open kubernetes config file")
	ErrKubernetesConfigCannotRead = errors.New("cannot read kubernetes config file")
)

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
		return "", fmt.Errorf("%w expanding kubernetes directory: %v", ErrKubernetesConfigAccess, err)
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

	_, err = kubernetesutil.WriteKubernetesConfig(file, config)
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
		return nil, fmt.Errorf("reading file '%s': %w", filename, err)
	}

	return kubernetesutil.ReadKubernetesConfig(data)
}

func (s *FileKubernetesConfigStore) Load(name string) (*models.KubernetesConfig, error) {
	storageDir, err := s.getAndExpandStorageDirectory()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(storageDir, name)
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: '%s'", ErrKubernetesConfigNotFound, filename)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%w: '%s'", ErrKubernetesConfigIsDir, filename)
	}

	config, err := s.loadFileAndParse(filename)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrPermission):
			return nil, err
		default:
			return nil, fmt.Errorf("%w: %v", ErrKubernetesConfigInvalid, err)
		}
	}

	return config, nil
}

func (s *FileKubernetesConfigStore) LoadFromPath(path string) (*models.KubernetesConfig, error) {
	targetPath, err := fileutil.ExpandPath(path)
	if err != nil {
		return nil, fmt.Errorf("%w expanding path '%s':%v", ErrKubernetesConfigAccess, path, err)
	}

	if !fileutil.CheckFileExists(targetPath) {
		return nil, fmt.Errorf("%w: '%s'", ErrKubernetesConfigNotFound, targetPath)
	}

	config, err := s.loadFileAndParse(targetPath)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrPermission):
			return nil, fmt.Errorf("%w: '%s'", err, targetPath)
		default:
			return nil, fmt.Errorf("%w: %v", ErrKubernetesConfigInvalid, err)
		}
	}

	return config, nil
}
