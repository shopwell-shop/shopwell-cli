package extension

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shyim/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shopwell-shop/shopwell-cli/internal/testhelper"
)

func TestGetExtensionByFolder_DetectsApp(t *testing.T) {
	tmpDir := testhelper.NewApp(t, "TestApp")

	ext, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)
	assert.Equal(t, TypePlatformApp, ext.GetType())

	name, err := ext.GetName()
	require.NoError(t, err)
	assert.Equal(t, "TestApp", name)
}

func TestGetExtensionByFolder_DetectsPlatformPlugin(t *testing.T) {
	// Create composer.json for a PlatformPlugin
	tmpDir := testhelper.ExtensionDir(t, testhelper.ComposerJSON{
		Name:        "test/test-plugin",
		Type:        "shopwell-platform-plugin",
		Version:     "1.0.0",
		License:     "MIT",
		Description: "Test plugin",
		Authors:     []string{"Test"},
		Require:     map[string]string{"shopwell/core": "~6.5.0"},
		Psr4:        map[string]string{`Test\TestPlugin\`: "src/"},
		PluginClass: `Test\TestPlugin\TestPlugin`,
		Label:       map[string]string{"de-DE": "Test Plugin", "en-GB": "Test Plugin"},
		Extra: map[string]any{
			"description":      map[string]string{"de-DE": "Ein Test Plugin", "en-GB": "A test plugin"},
			"manufacturerLink": map[string]string{"de-DE": "https://example.com", "en-GB": "https://example.com"},
			"supportLink":      map[string]string{"de-DE": "https://example.com/support", "en-GB": "https://example.com/support"},
		},
	})

	ext, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)
	assert.Equal(t, TypePlatformPlugin, ext.GetType())

	name, err := ext.GetName()
	require.NoError(t, err)
	assert.Equal(t, "TestPlugin", name)
}

func TestGetExtensionByFolder_DetectsShopwellBundle(t *testing.T) {
	// Create composer.json for a ShopwellBundle
	tmpDir := testhelper.ExtensionDir(t, testhelper.ComposerJSON{
		Name:    "test/test-bundle",
		Type:    "shopwell-bundle",
		Version: "1.0.0",
		License: "MIT",
		Require: map[string]string{"shopwell/core": "~6.5.0"},
		Psr4:    map[string]string{`Test\TestBundle\`: "src/"},
		Extra:   map[string]any{"shopwell-bundle-name": "TestBundle"},
	})

	ext, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)
	assert.Equal(t, TypeShopwellBundle, ext.GetType())

	name, err := ext.GetName()
	require.NoError(t, err)
	assert.Equal(t, "TestBundle", name)
}

func TestGetExtensionByFolder_RejectsShopwell5Plugin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create plugin.xml for a Shopwell 5 plugin
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "plugin.xml"), []byte("<plugin></plugin>"), 0644))

	_, err := GetExtensionByFolder(t.Context(), tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shopwell 5 is not supported")
}

func TestGetExtensionByFolder_RejectsUnknownType(t *testing.T) {
	tmpDir := t.TempDir()

	// Empty directory - no manifest.xml, no composer.json
	_, err := GetExtensionByFolder(t.Context(), tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown extension type")
}

func TestGetExtensionByFolder_PrefersManifestOverComposer(t *testing.T) {
	// Create both manifest.xml and composer.json
	tmpDir := testhelper.NewApp(t, "TestApp")

	testhelper.WriteFile(t, filepath.Join(tmpDir, "composer.json"), testhelper.ComposerJSON{
		Name:    "test/test-plugin",
		Type:    "shopwell-platform-plugin",
		Version: "1.0.0",
	}.String())

	ext, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)
	// Should detect as App since manifest.xml is checked first
	assert.Equal(t, TypePlatformApp, ext.GetType())
}

func TestUpdateMetaData_PlatformPlugin(t *testing.T) {
	tmpDir := testhelper.ExtensionDir(t, testhelper.ComposerJSON{
		Name:        "test/test-plugin",
		Type:        "shopwell-platform-plugin",
		Version:     "1.0.0",
		License:     "MIT",
		Description: "Test plugin",
		Authors:     []string{"Test"},
		Require:     map[string]string{"shopwell/core": "~6.5.0"},
		Psr4:        map[string]string{`Test\TestPlugin\`: "src/"},
		PluginClass: `Test\TestPlugin\TestPlugin`,
		Label:       map[string]string{"de-DE": "Altes Label DE", "en-GB": "Old Label EN"},
		Extra: map[string]any{
			"description": map[string]string{"de-DE": "Alte Beschreibung", "en-GB": "Old description"},
		},
	})

	ext, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)

	err = ext.UpdateMetaData(&ExtensionMetadata{
		Label: ExtensionTranslated{
			German:  "Neues Label DE",
			English: "New Label EN",
		},
		Description: ExtensionTranslated{
			German:  "Neue Beschreibung",
			English: "New description",
		},
	})
	require.NoError(t, err)

	// Re-read the extension to verify the changes were persisted
	ext2, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)

	meta := ext2.GetMetaData()
	assert.Equal(t, "Neues Label DE", meta.Label.German)
	assert.Equal(t, "New Label EN", meta.Label.English)
	assert.Equal(t, "Neue Beschreibung", meta.Description.German)
	assert.Equal(t, "New description", meta.Description.English)
}

func TestUpdateMetaData_App(t *testing.T) {
	tmpDir := t.TempDir()

	manifestContent := `<?xml version="1.0" encoding="UTF-8"?>
<manifest xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://raw.githubusercontent.com/shopwell-shop/shopwell/trunk/src/Core/Framework/App/Manifest/Schema/manifest-2.0.xsd">
    <meta>
        <name>TestApp</name>
        <label>Old Label EN</label>
        <label lang="de-DE">Altes Label DE</label>
        <future-label-metadata>keep</future-label-metadata>
        <description>Old description</description>
        <description lang="de-DE">Alte Beschreibung</description>
        <author>Test Author</author>
        <copyright>(c) Test</copyright>
        <version>1.0.0</version>
        <license>MIT</license>
    </meta>
</manifest>`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "manifest.xml"), []byte(manifestContent), 0644))

	ext, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)

	err = ext.UpdateMetaData(&ExtensionMetadata{
		Label: ExtensionTranslated{
			German:  "Neues Label DE",
			English: "New Label EN",
		},
		Description: ExtensionTranslated{
			German:  "Neue Beschreibung",
			English: "New description",
		},
	})
	require.NoError(t, err)

	// Re-read the extension to verify the changes were persisted
	ext2, err := GetExtensionByFolder(t.Context(), tmpDir)
	require.NoError(t, err)

	meta := ext2.GetMetaData()
	assert.Equal(t, "Neues Label DE", meta.Label.German)
	assert.Equal(t, "New Label EN", meta.Label.English)
	assert.Equal(t, "Neue Beschreibung", meta.Description.German)
	assert.Equal(t, "New description", meta.Description.English)

	manifestBytes, err := os.ReadFile(filepath.Join(tmpDir, "manifest.xml"))
	require.NoError(t, err)
	manifest := string(manifestBytes)
	assert.Contains(t, manifest, "<future-label-metadata>keep</future-label-metadata>")
	assert.Contains(t, manifest, `xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`)
	assert.Contains(t, manifest, `xsi:noNamespaceSchemaLocation="https://raw.githubusercontent.com/shopwell-shop/shopwell/trunk/src/Core/Framework/App/Manifest/Schema/manifest-2.0.xsd"`)
	assert.NotContains(t, manifest, "_xmlns")
}

func TestGetShopwellVersionConstraintFromComposer(t *testing.T) {
	t.Run("uses composer require", func(t *testing.T) {
		constraint, err := getShopwellVersionConstraintFromComposer(map[string]string{
			"shopwell/core": "~6.5.0",
		})
		require.NoError(t, err)
		assert.NotNil(t, constraint)
		assert.True(t, constraint.Check(version.Must(version.NewVersion("6.5.0.0"))))
	})

	t.Run("ignores the build override and always uses composer require", func(t *testing.T) {
		// The build override must never leak into the reported compatibility constraint,
		// otherwise account uploads would report the wrong compatible Shopwell versions.
		constraint, err := getShopwellVersionConstraintFromComposer(map[string]string{
			"shopwell/core": "~6.4.0",
		})
		require.NoError(t, err)
		assert.NotNil(t, constraint)
		assert.True(t, constraint.Check(version.Must(version.NewVersion("6.4.0.0"))))
	})

	t.Run("returns error when shopwell/core not in require", func(t *testing.T) {
		_, err := getShopwellVersionConstraintFromComposer(map[string]string{
			"php": ">=8.1",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "shopwell/core is required")
	})

	t.Run("returns error for invalid constraint in composer", func(t *testing.T) {
		_, err := getShopwellVersionConstraintFromComposer(map[string]string{
			"shopwell/core": "invalid[constraint",
		})
		assert.Error(t, err)
	})
}

func TestGetShopwellBuildVersionConstraint(t *testing.T) {
	t.Run("returns constraint when build override is set", func(t *testing.T) {
		config := &Config{
			Build: ConfigBuild{
				ShopwellVersionConstraint: "~6.5.0",
			},
		}

		constraint, err := GetShopwellBuildVersionConstraint(config)
		require.NoError(t, err)
		require.NotNil(t, constraint)
		assert.True(t, constraint.Check(version.Must(version.NewVersion("6.5.0.0"))))
	})

	t.Run("returns nil when build override is not set", func(t *testing.T) {
		constraint, err := GetShopwellBuildVersionConstraint(&Config{})
		require.NoError(t, err)
		assert.Nil(t, constraint)
	})

	t.Run("handles nil config", func(t *testing.T) {
		constraint, err := GetShopwellBuildVersionConstraint(nil)
		require.NoError(t, err)
		assert.Nil(t, constraint)
	})

	t.Run("returns error for invalid constraint", func(t *testing.T) {
		config := &Config{
			Build: ConfigBuild{
				ShopwellVersionConstraint: "invalid[constraint",
			},
		}

		_, err := GetShopwellBuildVersionConstraint(config)
		assert.Error(t, err)
	})
}

func mustConstraint(t *testing.T, s string) *version.Constraints {
	t.Helper()
	c, err := version.NewConstraint(s)
	require.NoError(t, err)
	return &c
}

func TestGetShopwellVersionConstraintForBuild(t *testing.T) {
	t.Run("uses build override when set", func(t *testing.T) {
		ext := &mockExtension{
			config:     &Config{Build: ConfigBuild{ShopwellVersionConstraint: "~6.5.0"}},
			constraint: mustConstraint(t, "~6.4.0"),
		}

		constraint, err := GetShopwellVersionConstraintForBuild(ext)
		require.NoError(t, err)
		require.NotNil(t, constraint)
		assert.True(t, constraint.Check(version.Must(version.NewVersion("6.5.0.0"))))
		assert.False(t, constraint.Check(version.Must(version.NewVersion("6.4.0.0"))))
	})

	t.Run("falls back to compatibility constraint when no override", func(t *testing.T) {
		ext := &mockExtension{
			config:     &Config{},
			constraint: mustConstraint(t, "~6.4.0"),
		}

		constraint, err := GetShopwellVersionConstraintForBuild(ext)
		require.NoError(t, err)
		require.NotNil(t, constraint)
		assert.True(t, constraint.Check(version.Must(version.NewVersion("6.4.0.0"))))
	})
}
