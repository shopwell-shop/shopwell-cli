package extension

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shyim/go-version"

	"github.com/shopwell-shop/shopwell-cli/internal/validation"
)

type ShopwellBundle struct {
	path     string
	Composer shopwellBundleComposerJson
	config   *Config
}

func newShopwellBundle(ctx context.Context, path string) (*ShopwellBundle, error) {
	composerJsonFile := path + "/composer.json"
	if _, err := os.Stat(composerJsonFile); err != nil {
		return nil, err
	}

	jsonFile, err := os.ReadFile(composerJsonFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read composer.json: %w", err)
	}

	var composerJson shopwellBundleComposerJson
	err = json.Unmarshal(jsonFile, &composerJson)
	if err != nil {
		return nil, fmt.Errorf("cannot parse composer.json: %w", err)
	}

	if composerJson.Type != "shopwell-bundle" {
		return nil, errors.New("composer.json type is not shopwell-bundle")
	}

	if composerJson.Extra.BundleName == "" {
		return nil, errors.New("composer.json does not contain shopwell-bundle-name in extra")
	}

	cfg, err := readExtensionConfig(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("cannot read extension config: %w", err)
	}

	extension := ShopwellBundle{
		Composer: composerJson,
		path:     path,
		config:   cfg,
	}

	return &extension, nil
}

type composerAutoload struct {
	Psr4 map[string]string `json:"psr-4"`
}

type shopwellBundleComposerJson struct {
	Name     string                          `json:"name"`
	Type     string                          `json:"type"`
	License  string                          `json:"license"`
	Version  string                          `json:"version"`
	Require  map[string]string               `json:"require"`
	Extra    shopwellBundleComposerJsonExtra `json:"extra"`
	Suggest  map[string]string               `json:"suggest"`
	Autoload composerAutoload                `json:"autoload"`
}

type shopwellBundleComposerJsonExtra struct {
	BundleName string `json:"shopwell-bundle-name"`
}

func (p ShopwellBundle) GetComposerName() (string, error) {
	return p.Composer.Name, nil
}

// GetRootDir returns the src directory of the bundle.
func (p ShopwellBundle) GetRootDir() string {
	return filepath.Join(p.path, "src")
}

func (p ShopwellBundle) GetSourceDirs() []string {
	var result []string

	for _, val := range p.Composer.Autoload.Psr4 {
		result = append(result, filepath.Join(p.path, val))
	}

	return result
}

// GetResourcesDir returns the resources directory of the shopwell bundle.
func (p ShopwellBundle) GetResourcesDir() string {
	return filepath.Join(p.GetRootDir(), "Resources")
}

func (p ShopwellBundle) GetResourcesDirs() []string {
	var result []string

	for _, val := range p.GetSourceDirs() {
		result = append(result, filepath.Join(val, "Resources"))
	}

	return result
}

func (p ShopwellBundle) GetName() (string, error) {
	return p.Composer.Extra.BundleName, nil
}

func (p ShopwellBundle) GetExtensionConfig() *Config {
	return p.config
}

func (p ShopwellBundle) GetShopwellVersionConstraint() (*version.Constraints, error) {
	return getShopwellVersionConstraintFromComposer(p.Composer.Require)
}

func (ShopwellBundle) GetType() string {
	return TypeShopwellBundle
}

func (p ShopwellBundle) GetVersion() (*version.Version, error) {
	return version.NewVersion(p.Composer.Version)
}

func (p ShopwellBundle) GetChangelog() (*ExtensionChangelog, error) {
	return parseExtensionMarkdownChangelog(p)
}

func (p ShopwellBundle) GetLicense() (string, error) {
	return p.Composer.License, nil
}

func (p ShopwellBundle) GetPath() string {
	return p.path
}

func (p ShopwellBundle) GetIconPath() string {
	return ""
}

func (p ShopwellBundle) GetMetaData() *ExtensionMetadata {
	return &ExtensionMetadata{
		Label: ExtensionTranslated{
			German:  "FALLBACK",
			English: "FALLBACK",
		},
		Description: ExtensionTranslated{
			German:  "FALLBACK",
			English: "FALLBACK",
		},
	}
}

func (p ShopwellBundle) UpdateMetaData(_ *ExtensionMetadata) error {
	return nil
}

func (p ShopwellBundle) Validate(c context.Context, check validation.Check) {
	// ShopwellBundle validation is currently empty but signature updated to match interface
}
