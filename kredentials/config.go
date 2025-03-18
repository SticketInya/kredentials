package kredentials

import "os"

const (
	defaultKredentialStorageDir     string      = "~/.kredentials/configs"
	defaultKubernetesStorageDir     string      = "~/.kube"
	defaultStorageDirPermissions    os.FileMode = 0755
	defaultKubernetesDirPermissions os.FileMode = 0755
)

type KredentialsConfig struct {
	kredentialStorageDir            string
	kredentialStorageDirPermissions os.FileMode
	kubernetesStorageDir            string
	kubernetesStorageDirPermissions os.FileMode
}

func NewKredentialsDefaultConfig() *KredentialsConfig {
	return &KredentialsConfig{
		kredentialStorageDir:            defaultKredentialStorageDir,
		kredentialStorageDirPermissions: defaultStorageDirPermissions,
		kubernetesStorageDir:            defaultKubernetesStorageDir,
		kubernetesStorageDirPermissions: defaultKubernetesDirPermissions,
	}
}
