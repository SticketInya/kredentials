package kubernetesutil

import (
	"fmt"
	"io"

	"github.com/SticketInya/kredentials/models"
	"k8s.io/client-go/tools/clientcmd"
)

func WriteKubernetesConfig(w io.Writer, config models.KubernetesConfig) (n int, err error) {
	data, err := clientcmd.Write(config)
	if err != nil {
		return 0, fmt.Errorf("writing kubernetes config: %w", err)
	}

	return w.Write(data)
}

func ReadKubernetesConfig(data []byte) (*models.KubernetesConfig, error) {
	config, err := clientcmd.Load(data)
	if err != nil {
		return nil, fmt.Errorf("reading kubernetes config: %w", err)
	}
	return config, nil
}
