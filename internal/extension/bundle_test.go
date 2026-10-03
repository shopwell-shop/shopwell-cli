package extension

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shopwell-shop/shopwell-cli/internal/testhelper"
)

func TestCreateBundleEmptyFolder(t *testing.T) {
	dir := t.TempDir()

	bundle, err := newShopwellBundle(t.Context(), dir)
	assert.Error(t, err)
	assert.Nil(t, bundle)
}

func TestCreateBundleInvalidComposerType(t *testing.T) {
	dir := testhelper.ExtensionDir(t, testhelper.ComposerJSON{
		Name: "shopwell/invalid",
		Type: "invalid",
	})

	bundle, err := newShopwellBundle(t.Context(), dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "composer.json type is not shopwell-bundle")
	assert.Nil(t, bundle)
}

func TestCreateBundleMissingName(t *testing.T) {
	dir := testhelper.ExtensionDir(t, testhelper.ComposerJSON{
		Name: "shopwell/invalid",
		Type: "shopwell-bundle",
	})

	bundle, err := newShopwellBundle(t.Context(), dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "composer.json does not contain shopwell-bundle-name")
	assert.Nil(t, bundle)
}

func TestCreateBundle(t *testing.T) {
	dir := testhelper.ExtensionDir(t, testhelper.ComposerJSON{
		Name:    "shopwell/invalid",
		Version: "1.0.0",
		Type:    "shopwell-bundle",
		Extra:   map[string]any{"shopwell-bundle-name": "TestBundle"},
		Psr4:    map[string]string{`TestBundle\`: "src/"},
	})

	bundle, err := newShopwellBundle(t.Context(), dir)
	assert.NoError(t, err)

	name, err := bundle.GetName()
	assert.NoError(t, err)

	assert.Equal(t, "TestBundle", name)
	assert.Equal(t, path.Join(dir, "src"), bundle.GetRootDir())
	assert.Equal(t, dir, bundle.GetPath())
	assert.Equal(t, path.Join(dir, "src", "Resources"), bundle.GetResourcesDir())
	assert.Equal(t, path.Join(dir, "src", "Resources"), bundle.GetResourcesDirs()[0])
	assert.Equal(t, TypeShopwellBundle, bundle.GetType())

	_, err = bundle.GetChangelog()
	// changelog is missing
	assert.Error(t, err)

	version, err := bundle.GetVersion()
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", version.String())

	// does nothing
	bundle.Validate(getTestContext(), &testCheck{})

	assert.Equal(t, "FALLBACK", bundle.GetMetaData().Description.German)
}
