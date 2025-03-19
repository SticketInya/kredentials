package kredentials

import (
	"fmt"

	"github.com/SticketInya/kredentials/models"
	"github.com/SticketInya/kredentials/storage"
)

const (
	kubernetesConfigFilename       string = "config"
	kubernetesConfigBackupFilename string = "config.last"
)

type KredentialManager struct {
	kredStore   storage.KredentialStore
	configStore storage.KubernetesConfigStore
}

func NewKredentialManager(kredentialStore storage.KredentialStore, configStore storage.KubernetesConfigStore) *KredentialManager {
	return &KredentialManager{
		kredStore:   kredentialStore,
		configStore: configStore,
	}
}

func (m *KredentialManager) AddKredential(name string, path string) error {
	kredentials, err := m.ListKredentials()
	if err != nil {
		return err
	}

	for _, existing := range kredentials {
		if name == existing.Name {
			return ErrKredentialConflict{name}
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
