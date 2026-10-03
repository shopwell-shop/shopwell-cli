package shop

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindClosestShopwellProjectFallback(t *testing.T) {
	t.Setenv("PROJECT_ROOT", "")
	t.Chdir(t.TempDir())

	_, err := FindClosestShopwellProject(false)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "cannot find Shopwell project")
	}

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		return
	}

	root, err := FindClosestShopwellProject(true)
	if assert.NoError(t, err) {
		assert.Equal(t, cwd, root)
	}
}

func TestFindClosestShopwellProject(t *testing.T) {
	t.Setenv("PROJECT_ROOT", "")
	projectDir := createShopwellProject(t, "composer.json")
	nestedDir := filepath.Join(projectDir, "custom", "plugins", "Example")
	require.NoError(t, os.MkdirAll(nestedDir, 0o755))
	t.Chdir(nestedDir)

	found, err := FindClosestShopwellProject(false)

	require.NoError(t, err)
	assert.Equal(t, projectDir, found)
}

func TestFindClosestShopwellProjectFromLockFile(t *testing.T) {
	t.Setenv("PROJECT_ROOT", "")
	projectDir := createShopwellProject(t, "composer.lock")
	t.Chdir(projectDir)

	found, err := FindClosestShopwellProject(false)

	require.NoError(t, err)
	assert.Equal(t, projectDir, found)
}

func TestFindClosestShopwellProjectUsesEnvironmentOverride(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "partly-configured-project")
	t.Setenv("PROJECT_ROOT", projectDir+string(filepath.Separator))

	found, err := FindClosestShopwellProject(false)

	require.NoError(t, err)
	assert.Equal(t, filepath.Clean(projectDir), found)
}

func TestFindClosestShopwellProjectNotFound(t *testing.T) {
	t.Setenv("PROJECT_ROOT", "")
	t.Chdir(t.TempDir())

	_, err := FindClosestShopwellProject(false)

	assert.ErrorContains(t, err, "cannot find Shopwell project in ")
	assert.ErrorContains(t, err, "bin/console")
}

func TestFindClosestShopwellProjectRequiresConsole(t *testing.T) {
	t.Setenv("PROJECT_ROOT", "")
	projectDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(projectDir, "composer.json"),
		[]byte(`{"require":{"shopwell/core":"*" }}`),
		0o644,
	))
	t.Chdir(projectDir)

	_, err := FindClosestShopwellProject(false)

	assert.ErrorContains(t, err, "cannot find Shopwell project")
}

func createShopwellProject(t *testing.T, composerFile string) string {
	t.Helper()

	projectDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(projectDir, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "bin", "console"), nil, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(projectDir, composerFile),
		[]byte(`{"require":{"shopwell/core":"*" }}`),
		0o644,
	))
	return projectDir
}
