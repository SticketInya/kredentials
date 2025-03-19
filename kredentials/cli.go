package kredentials

import (
	"os"

	"github.com/SticketInya/kredentials/formatter"
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
