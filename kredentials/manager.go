package kredentials

import (
	"errors"
	"fmt"

	"github.com/SticketInya/kredentials/internal/fileutil"
	"github.com/SticketInya/kredentials/models"
)

const (
	kubernetesConfigFilename       string = "config"
	kubernetesConfigBackupFilename string = "config.last"
	kredentialBackupFilename       string = "kredential_backup"
)

type KredentialManager struct {
	kredStore    KredentialStore
	configStore  KubernetesConfigStore
	archiveStore ArchiveStore
}

func NewKredentialManager(
	kredentialStore KredentialStore,
	configStore KubernetesConfigStore,
	archiveStore ArchiveStore,
) *KredentialManager {
	return &KredentialManager{
		kredStore:    kredentialStore,
		configStore:  configStore,
		archiveStore: archiveStore,
	}
}

func (m *KredentialManager) AddKredential(name string, path string, options models.AddKredentialOptions) error {
	kredentials, err := m.ListKredentials()
	if err != nil {
		return err
	}

	if !options.OverwriteExisting {
		for _, existing := range kredentials {
			if name == existing.Name {
				return ErrKredentialConflict{name}
			}
		}
	}

	kubernetesConfig, err := m.configStore.LoadFromPath(path)
	if err != nil {
		return fmt.Errorf("loading kubernetes config: %w", err)
	}

	kred := models.NewKredential(name, kubernetesConfig)
	return m.kredStore.Store(kred)
}

func (m *KredentialManager) LoadKredential(name string) (*models.Kredential, error) {
	return m.kredStore.Load(name)
}

func (m *KredentialManager) ListKredentials() ([]*models.Kredential, error) {
	return m.kredStore.List()
}

func (m *KredentialManager) DeleteKredential(name string) {
	_ = m.kredStore.Delete(name)
}

func (m *KredentialManager) UseKredential(name string) error {
	kred, err := m.kredStore.Load(name)
	if err != nil {
		return err
	}

	if err = m.createKubernetesConfigBackup(); err != nil {
		return err
	}

	err = m.configStore.Store(kubernetesConfigFilename, *kred.Config)
	if err != nil {
		return fmt.Errorf("overwriting kubernetes config: %w", err)
	}
	return nil
}

func (m *KredentialManager) RevertKredential() error {
	lastConfig, err := m.configStore.Load(kubernetesConfigBackupFilename)
	if err != nil {
		return fmt.Errorf("loading last kubernetes config: %w", err)
	}

	if err = m.createKubernetesConfigBackup(); err != nil {
		return err
	}

	if err = m.configStore.Store(kubernetesConfigFilename, *lastConfig); err != nil {
		return fmt.Errorf("reverting kubernetes config: %w", err)
	}

	return nil
}

func (m *KredentialManager) CreateKredentialBackup(path string) error {
	kredentials, err := m.ListKredentials()
	if err != nil {
		return fmt.Errorf("collecting kredentials: %w", err)
	}

	expandedPath, err := fileutil.ExpandPath(path)
	if err != nil {
		return fmt.Errorf("extending path '%s': %w", path, err)
	}

	if err = m.archiveStore.Store(expandedPath, kredentialBackupFilename, kredentials); err != nil {
		return fmt.Errorf("creating archive: %w", err)
	}

	return nil
}

func (m *KredentialManager) RestoreKredentialBackup(path string) error {
	expandedPath, err := fileutil.ExpandPath(path)
	if err != nil {
		return fmt.Errorf("expanding path '%s': %w", path, err)
	}

	kreds, err := m.archiveStore.Load(expandedPath)
	if err != nil {
		return fmt.Errorf("restoring from archive: %w", err)
	}

	errs := []error{}
	for _, kred := range kreds {
		if err := m.kredStore.Store(kred); err != nil {
			errs = append(errs, fmt.Errorf("restoring kredential '%s': %w", kred.Name, err))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to restore kredentials: %w", errors.Join(errs...))
	}

	return nil
}

func (m *KredentialManager) createKubernetesConfigBackup() error {
	currentConfig, err := m.configStore.Load(kubernetesConfigFilename)
	if err != nil {
		return fmt.Errorf("loading current kubernetes configuration: %w", err)
	}

	err = m.configStore.Store(kubernetesConfigBackupFilename, *currentConfig)
	if err != nil {
		return fmt.Errorf("saving kubernetes config backup: %w", err)
	}

	return nil
}
