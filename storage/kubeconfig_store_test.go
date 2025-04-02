package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SticketInya/kredentials/internal/assert"
	"github.com/SticketInya/kredentials/models"
	"github.com/SticketInya/kredentials/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	defaultDirectoryPermission os.FileMode = 0755
)

var (
	validTestKubernetesConfig models.KubernetesConfig = models.KubernetesConfig{
		Clusters: map[string]*api.Cluster{
			"test-cluster": {
				Server:     "test-server",
				Extensions: map[string]runtime.Object{},
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"test-auth": {
				Token:      "test-token",
				Exec:       &api.ExecConfig{},
				Extensions: map[string]runtime.Object{},
			},
		},
		Contexts: map[string]*api.Context{
			"test-context": {
				Namespace:  "test-ns",
				Extensions: map[string]runtime.Object{},
			},
		},
		Extensions: map[string]runtime.Object{},
		Preferences: api.Preferences{
			Colors:     true,
			Extensions: map[string]runtime.Object{},
		},
	}
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
			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
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

			store := storage.NewFileKubernetesConfigStore(unwritableDir, defaultDirectoryPermission)
			assert.Error(t, store.Store(tc.configName, tc.config), os.ErrPermission)
		}

	})
}

func TestFileKubernetestConfigStore_Load(t *testing.T) {

	t.Run("load valid kubernetes config", func(t *testing.T) {
		testCases := []struct {
			configName string
			config     models.KubernetesConfig
		}{
			{
				configName: "test-config",
				config:     validTestKubernetesConfig,
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			assert.NoError(t, store.Store(tc.configName, tc.config))

			config, err := store.Load(tc.configName)
			assert.NoError(t, err)
			assert.DeepEqual(t, config, &tc.config)
		}
	})

	t.Run("load unreadable file", func(t *testing.T) {
		testCases := []struct {
			configName string
			config     models.KubernetesConfig
		}{
			{
				configName: "test-config",
				config: models.KubernetesConfig{
					Clusters: map[string]*api.Cluster{
						"test-cluster": {
							Server: "test-server",
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			assert.NoError(t, store.Store(tc.configName, tc.config))

			os.Chmod(filepath.Join(tempDir, tc.configName), 0333)

			_, err := store.Load(tc.configName)
			assert.Error(t, err, os.ErrPermission)

		}
	})

	t.Run("load invalid config", func(t *testing.T) {
		testCases := []struct {
			configName string
			config     []byte
		}{
			{
				configName: "test-config",
				config:     []byte{0x00},
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			file, err := os.Create(filepath.Join(tempDir, tc.configName))
			assert.NoError(t, err)

			_, err = file.Write(tc.config)
			assert.NoError(t, err)

			_, err = store.Load(tc.configName)
			assert.Error(t, err, storage.ErrKubernetesConfigInvalid)
		}
	})

	t.Run("load missing file", func(t *testing.T) {
		testCases := []struct {
			configName string
		}{
			{
				configName: "test-config",
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)

			_, err := store.Load(tc.configName)
			assert.Error(t, err, storage.ErrKubernetesConfigNotFound)
		}
	})

	t.Run("path to directory", func(t *testing.T) {
		testCases := []struct {
			configName string
		}{
			{
				configName: "test-config",
			},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			assert.NoError(t, os.MkdirAll(filepath.Join(tempDir, tc.configName), defaultDirectoryPermission))
			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)

			_, err := store.Load(tc.configName)
			assert.Error(t, err, storage.ErrKubernetesConfigIsDir)
		}
	})
}

func TestFileKubernetesConfigStore_LoadFromPath(t *testing.T) {
	t.Run("loading valid config from valid path", func(t *testing.T) {
		testCases := []struct {
			path   string
			config models.KubernetesConfig
		}{
			{path: "/configs/test-config", config: validTestKubernetesConfig},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			actualPath := filepath.Join(tempDir, tc.path)

			assert.NoError(t, os.MkdirAll(filepath.Dir(actualPath), defaultDirectoryPermission))
			assert.NoError(t, store.Store(tc.path, tc.config))

			config, err := store.LoadFromPath(actualPath)
			assert.NoError(t, err)
			assert.DeepEqual(t, config, &tc.config)
		}
	})

	t.Run("loading from invalid path", func(t *testing.T) {
		testCases := []struct {
			path string
		}{
			{path: "this/is/invalid/path"},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			actualPath := filepath.Join(tempDir, tc.path)

			_, err := store.LoadFromPath(actualPath)
			assert.Error(t, err, storage.ErrKubernetesConfigNotFound)
		}
	})

	t.Run("loading from unreadable path", func(t *testing.T) {
		testCases := []struct {
			path   string
			config models.KubernetesConfig
		}{
			{path: "super/duper/nested/dir/config", config: validTestKubernetesConfig},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			actualPath := filepath.Join(tempDir, tc.path)

			assert.NoError(t, os.MkdirAll(filepath.Dir(actualPath), defaultDirectoryPermission))
			assert.NoError(t, store.Store(tc.path, tc.config))
			assert.NoError(t, os.Chmod(actualPath, 0333))

			_, err := store.LoadFromPath(actualPath)
			assert.Error(t, err, os.ErrPermission)
		}
	})

	t.Run("loading invalid config", func(t *testing.T) {
		testCases := []struct {
			path   string
			config []byte
		}{
			{path: "path/to/config", config: []byte{0x01}},
		}

		for _, tc := range testCases {
			tempDir := t.TempDir()

			store := storage.NewFileKubernetesConfigStore(tempDir, defaultDirectoryPermission)
			actualPath := filepath.Join(tempDir, tc.path)
			assert.NoError(t, os.MkdirAll(filepath.Dir(actualPath), defaultDirectoryPermission))

			file, err := os.Create(actualPath)
			assert.NoError(t, err)

			_, err = file.Write(tc.config)
			assert.NoError(t, err)

			_, err = store.LoadFromPath(actualPath)
			assert.Error(t, err, storage.ErrKubernetesConfigInvalid)
		}
	})
}
