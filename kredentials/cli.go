package kredentials

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/formatter"
	"github.com/SticketInya/kredentials/internal/fileutil"
	"github.com/SticketInya/kredentials/storage"
)

type KredentialsCli struct {
	config  *KredentialsConfig
	Manager *KredentialManager
	Printer *formatter.StructuredPrinter
}

func NewKredentialsCli(config *KredentialsConfig) *KredentialsCli {
	kredentialsStore := storage.NewFileKredentialStore(config.kredentialStorageDir, config.kredentialStorageDirPermissions)
	kubernetesStore := storage.NewFileKubernetesConfigStore(config.kubernetesStorageDir, config.kubernetesStorageDirPermissions)
	manager := NewKredentialManager(kredentialsStore, kubernetesStore)
	printer := formatter.NewStructuredPrinter(os.Stdout)

	return &KredentialsCli{
		config:  config,
		Manager: manager,
		Printer: printer,
	}
}

func (cli *KredentialsCli) GetVersion() VersionConfig {
	return cli.config.versionConfig
}

func (cli *KredentialsCli) Initialize() error {
	if err := fileutil.EnsureDirectory(cli.config.kubernetesStorageDir, cli.config.kubernetesStorageDirPermissions); err != nil {
		return fmt.Errorf("cannot initialize kredentials cli: %w", err)
	}

	if err := fileutil.EnsureDirectory(cli.config.kredentialStorageDir, cli.config.kredentialStorageDirPermissions); err != nil {
		return fmt.Errorf("cannot initialize kredentials cli: %w", err)
	}

	return nil
}
