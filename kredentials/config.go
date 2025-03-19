package kredentials

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SticketInya/kredentials/internal/fileutil"
)

const (
	defaultKredentialDir            string      = ".kredentials"
	defaultKredentialStorageDir     string      = "configs"
	defaultKubernetesStorageDir     string      = ".kube"
	defaultStorageDirPermissions    os.FileMode = 0755
	defaultKubernetesDirPermissions os.FileMode = 0755
	defaultArchiveDirPermissions    os.FileMode = 0755

	// Environment variables
	kubernetesConfigCustomDirEnvKey string = "KUBECONFIG"
	xdgConfigHomeDirEnvKey          string = "XDG_CONFIG_HOME"
	kredentialConfigCustomDirEnvKey string = "KREDENTIAL_CONFIG_HOME"
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
	archiveStorageDirPermissions    os.FileMode

	versionConfig VersionConfig
}

func NewKredentialsDefaultConfig(versionConfig VersionConfig) (*KredentialsConfig, error) {
	kredentialStorageDir, err := getKredentialStorageDirectory()
	if err != nil {
		return nil, err
	}

	kubernetesStorageDir, err := getKubernetesStorageDirectory()
	if err != nil {
		return nil, err
	}

	return &KredentialsConfig{
		kredentialStorageDir:            kredentialStorageDir,
		kredentialStorageDirPermissions: defaultStorageDirPermissions,
		kubernetesStorageDir:            kubernetesStorageDir,
		kubernetesStorageDirPermissions: defaultKubernetesDirPermissions,
		archiveStorageDirPermissions:    defaultArchiveDirPermissions,

		versionConfig: versionConfig,
	}, nil
}

func getUserHomeDir() (string, error) {
	if xdgConfigHomeDir, ok := os.LookupEnv(xdgConfigHomeDirEnvKey); ok {
		xdgHomedir, err := fileutil.ExpandPath(xdgConfigHomeDir)
		if err != nil {
			return "", fmt.Errorf("failed to expand xdg home directory: %w", err)
		}
		return xdgHomedir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user home directory: %w", err)
	}
	return homeDir, nil
}

func getKubernetesStorageDirectory() (string, error) {
	if customPath, ok := os.LookupEnv(kubernetesConfigCustomDirEnvKey); ok {
		return fileutil.ExpandPath(customPath)
	}

	homeDir, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, defaultKubernetesStorageDir), nil
}

func getKredentialStorageDirectory() (string, error) {
	if customPath, ok := os.LookupEnv(kredentialConfigCustomDirEnvKey); ok {
		return fileutil.ExpandPath(filepath.Join(customPath, defaultKredentialStorageDir))
	}

	homeDir, err := getUserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, defaultKredentialDir, defaultKredentialStorageDir), nil
}
