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

type KredentialStore interface {
	Store(kred *models.Kredential) error
	Load(name string) (*models.Kredential, error)
	List() ([]*models.Kredential, error)
	Delete(name string) error
}

type FileKredentialStore struct {
	storageDirectory      string
	storageDirPermissions os.FileMode
}

func NewFileKredentialStore(storageDirectory string, storageDirPermissions os.FileMode) *FileKredentialStore {
	return &FileKredentialStore{
		storageDirectory:      storageDirectory,
		storageDirPermissions: storageDirPermissions,
	}
}

func (s *FileKredentialStore) getAndExpandStorageDirectory() (string, error) {
	destDir, err := fileutil.ExpandPath(s.storageDirectory)
	if err != nil {
		return "", fmt.Errorf("expandig storage directory: %w", err)
	}

	return destDir, nil
}

func (s *FileKredentialStore) Store(kred *models.Kredential) error {
	storageDir, err := s.getAndExpandStorageDirectory()
	if err != nil {
		return err
	}

	// ensure the directory exists
	if err = os.MkdirAll(storageDir, s.storageDirPermissions); err != nil {
		return fmt.Errorf("creating storage directory: %w", err)
	}

	filename := filepath.Join(storageDir, kred.Name)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file '%s': %w", filename, err)
	}
	defer file.Close()

	_, err = kubernetes_util.WriteKubernetesConfig(file, *kred.Config)
	if err != nil {
		return fmt.Errorf("writing '%s' file %w", filename, err)
	}
	return nil
}

func (s *FileKredentialStore) Load(name string) (*models.Kredential, error) {
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

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file '%s': %w", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file content '%s': %w", filename, err)
	}

	config, err := kubernetes_util.ReadKubernetesConfig(data)
	if err != nil {
		return nil, err
	}

	return models.NewKredential(name, config), nil
}

func (s *FileKredentialStore) List() ([]*models.Kredential, error) {
	var kreds []*models.Kredential

	storageDir, err := s.getAndExpandStorageDirectory()
	if err != nil {
		return kreds, err
	}

	pathInfo, err := os.Stat(storageDir)
	if err != nil {
		return kreds, fmt.Errorf("storage directory does not exists")
	}

	if !pathInfo.IsDir() {
		return kreds, fmt.Errorf("storage directory is actually a file")
	}

	files, err := os.ReadDir(storageDir)
	if err != nil {
		return kreds, fmt.Errorf("reading storage directory: %w", err)
	}

	if len(files) == 0 {
		return kreds, nil
	}

	for _, file := range files {
		kred, err := s.Load(file.Name())
		if err != nil {
			// should not stop the listing process if a file is invalid
			continue
		}
		kreds = append(kreds, kred)
	}

	return kreds, nil
}

func (s *FileKredentialStore) Delete(name string) error {
	storageDir, err := s.getAndExpandStorageDirectory()
	if err != nil {
		return err
	}

	filename := filepath.Join(storageDir, name)
	if err = os.Remove(filename); err != nil {
		return fmt.Errorf("deleting file '%s': %w", filename, err)
	}

	return nil
}
