package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SticketInya/kredentials/internal/assert"
	"github.com/SticketInya/kredentials/models"
	"github.com/SticketInya/kredentials/storage"
)

func TestFileKubernetesConfigStore_Store(t *testing.T) {

	t.Run("store valid config", func(t *testing.T) {
		testCases := []struct {
			configName string
			config     models.KubernetesConfig
		}{
			{
				configName: "test-config",
				config: models.KubernetesConfig{
					APIVersion: "v1",
				},
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()
			store := storage.NewFileKubernetesConfigStore(tempDir, 0755)
			assert.NoError(t, store.Store(tc.configName, tc.config))

			// Verify file was created
			filePath := filepath.Join(tempDir, tc.configName)
			_, err := os.Stat(filePath)
			assert.NoError(t, err)
		}
	})

	t.Run("store with unwritable directory", func(t *testing.T) {
		testCases := []struct {
			configName string
			config     models.KubernetesConfig
		}{
			{
				configName: "test-config",
				config: models.KubernetesConfig{
					APIVersion: "v1",
				},
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()
			unwritableDir := filepath.Join(tempDir, "unwritable")
			if err := os.Mkdir(unwritableDir, 0500); err != nil {
				t.Fatal(err)
			}

			store := storage.NewFileKubernetesConfigStore(unwritableDir, 0755)
			assert.Error(t, store.Store(tc.configName, tc.config), os.ErrPermission)
		}

	})
}
