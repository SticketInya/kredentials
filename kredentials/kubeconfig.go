package kredentials

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SticketInya/kredentials/internal/fileutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	kubernetesConfigFilePermissions = 0755
	kubernetesConfigFileName        = "config"
	kubernetesConfigDirectory       = "~/.kube/"
)

type KubernetesConfig = api.Config

func (k *Kredential) ReadKubernetesConfig(path string) error {
	targetPath, err := fileutil.ExpandPath(path)
	if err != nil {
		return fmt.Errorf("expanding target path '%s' %w", path, err)
	}

	config, err := clientcmd.LoadFromFile(targetPath)
	if err != nil {
		return fmt.Errorf("error reading file '%s' %w", targetPath, err)
	}
	k.Config = config

	return nil
}

func (k *Kredential) WriteKubernetesConfig(path string, name string) error {
	destDir, err := fileutil.ExpandPath(path)
	if err != nil {
		return fmt.Errorf("expanding target directory '%s' %w", destDir, err)
	}

	if err := os.MkdirAll(destDir, kubernetesConfigFilePermissions); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}
	newFilePath := filepath.Join(destDir, name)

	if err = clientcmd.WriteToFile(*k.Config, newFilePath); err != nil {
		return fmt.Errorf("writing config '%s' to '%s' %w", k.Name, newFilePath, err)
	}
	return nil
}

func (k *Kredential) SetAsKubernetesConfig() error {
	return k.WriteKubernetesConfig(kubernetesConfigDirectory, kubernetesConfigFileName)
}
