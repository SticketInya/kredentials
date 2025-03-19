package kredentials

import "os"

const (
	defaultKredentialStorageDir     string      = "~/.kredentials/configs"
	defaultKubernetesStorageDir     string      = "~/.kube"
	defaultStorageDirPermissions    os.FileMode = 0755
	defaultKubernetesDirPermissions os.FileMode = 0755
)

type VersionConfig struct {
	ApplicationVersion string
	CommitHash         string
	BuildDate          string
}

func NewVersionConfig(applicationVersion string, commitHash string, buildDate string) VersionConfig {
	return VersionConfig{
		ApplicationVersion: applicationVersion,
		CommitHash:         commitHash,
		BuildDate:          buildDate,
	}
}

type KredentialsConfig struct {
	kredentialStorageDir            string
	kredentialStorageDirPermissions os.FileMode
	kubernetesStorageDir            string
	kubernetesStorageDirPermissions os.FileMode

	versionConfig VersionConfig
}

func NewKredentialsDefaultConfig(versionConfig VersionConfig) *KredentialsConfig {
	return &KredentialsConfig{
		kredentialStorageDir:            defaultKredentialStorageDir,
		kredentialStorageDirPermissions: defaultStorageDirPermissions,
		kubernetesStorageDir:            defaultKubernetesStorageDir,
		kubernetesStorageDirPermissions: defaultKubernetesDirPermissions,

		versionConfig: versionConfig,
	}
}
