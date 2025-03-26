package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/SticketInya/kredentials/internal/fileutil"
	"github.com/SticketInya/kredentials/internal/kubernetesutil"
	"github.com/SticketInya/kredentials/models"
)

const (
	zipeArchiveExtension string = ".zip"
)

type ZipArchiveStore struct {
	storageDirPermissions os.FileMode
}

func NewZipArchiveStore(storageDirectoryPermissions os.FileMode) *ZipArchiveStore {
	return &ZipArchiveStore{
		storageDirPermissions: storageDirectoryPermissions,
	}
}

func (s *ZipArchiveStore) Store(path string, name string, kredentials []*models.Kredential) error {
	if err := fileutil.EnsureDirectory(path, s.storageDirPermissions); err != nil {
		return fmt.Errorf("cannot create archive directory: %w", err)
	}

	formattedName := s.formatArchiveName(name)
	archive, err := os.Create(filepath.Join(path, formattedName))
	if err != nil {
		return fmt.Errorf("creating zip archive '%s': %w", formattedName, err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	for _, kred := range kredentials {
		w, err := zipWriter.Create(kred.Name)
		if err != nil {
			return fmt.Errorf("cannot create kredential '%s' in archive: %w", kred.Name, err)
		}

		if _, err := kubernetesutil.WriteKubernetesConfig(w, *kred.Config); err != nil {
			return fmt.Errorf("cannot write kredential '%s' to archive: %w", kred.Name, err)
		}
	}

	return nil
}

func (s *ZipArchiveStore) Load(path string) ([]*models.Kredential, error) {
	var kredentials []*models.Kredential

	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return kredentials, fmt.Errorf("cannot open archive '%s': %w", path, err)
	}
	defer zipReader.Close()

	kreds := []*models.Kredential{}
	for _, file := range zipReader.File {
		scan, err := file.Open()
		if err != nil {
			return kredentials, fmt.Errorf("cannot open kredential '%s' in archive: %w", file.Name, err)
		}
		defer scan.Close()

		kred, err := s.loadKredentialFromArchive(file.Name, scan)
		if err != nil {
			return kredentials, fmt.Errorf("reading kredential '%s': %w", file.Name, err)
		}
		kreds = append(kreds, kred)
	}
	kredentials = kreds

	return kredentials, nil
}

func (s *ZipArchiveStore) loadKredentialFromArchive(name string, r io.Reader) (*models.Kredential, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read kredential: %w", err)
	}

	config, err := kubernetesutil.ReadKubernetesConfig(data)
	if err != nil {
		return nil, fmt.Errorf("cannot parse kredential: %w", err)
	}

	return models.NewKredential(name, config), nil
}

func (s *ZipArchiveStore) formatArchiveName(name string) string {
	if ext := filepath.Ext(name); ext == zipeArchiveExtension {
		return name
	}

	return name + zipeArchiveExtension
}
