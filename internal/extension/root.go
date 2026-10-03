package extension

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/shyim/go-version"

	"github.com/shopwell-shop/shopwell-cli/internal/archiver"
	"github.com/shopwell-shop/shopwell-cli/internal/validation"
)

const (
	TypePlatformApp    = "app"
	TypePlatformPlugin = "plugin"
	TypeShopwellBundle = "shopwell-bundle"

	ComposerTypePlugin = "shopwell-platform-plugin"
	ComposerTypeApp    = "shopwell-app"
	ComposerTypeBundle = "shopwell-bundle"
)

func GetExtensionByFolder(ctx context.Context, path string) (Extension, error) {
	if _, err := os.Stat(path + "/plugin.xml"); err == nil {
		return nil, errors.New("shopwell 5 is not supported. Please use https://github.com/FriendsOfShopware/FroshPluginUploader instead")
	}

	if _, err := os.Stat(path + "/manifest.xml"); err == nil {
		return newApp(ctx, path)
	}

	if _, err := os.Stat(path + "/composer.json"); err != nil {
		return nil, fmt.Errorf("unknown extension type: no composer.json or manifest.xml found in %s", path)
	}

	var ext Extension

	ext, err := newPlatformPlugin(ctx, path)
	if err != nil {
		if errors.Is(err, ErrPlatformInvalidType) {
			ext, err = newShopwellBundle(ctx, path)
		} else {
			return nil, err
		}
	}

	return ext, err
}

func GetExtensionByZip(ctx context.Context, filePath string) (Extension, error) {
	dir, err := os.MkdirTemp("", "extension")
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	file, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, err
	}

	err = archiver.Unzip(file, dir)
	if err != nil {
		return nil, err
	}

	fileName := file.File[0].Name

	if strings.Contains(fileName, "..") {
		return nil, errors.New("invalid zip file")
	}

	extName := strings.Split(fileName, "/")[0]
	return GetExtensionByFolder(ctx, fmt.Sprintf("%s/%s", dir, extName))
}

type ExtensionTranslated struct {
	German  string `json:"german"`
	English string `json:"english"`
}

type ExtensionChangelog struct {
	German     string `json:"german"`
	English    string `json:"english"`
	Changelogs map[string]string
}

type ExtensionMetadata struct {
	Name        string
	Label       ExtensionTranslated
	Description ExtensionTranslated
}

type Extension interface {
	GetName() (string, error)
	GetComposerName() (string, error)
	// Deprecated: use the list variation instead
	GetResourcesDir() string
	GetResourcesDirs() []string

	GetIconPath() string

	// GetRootDir Returns the root folder where the code is located plugin -> src, app ->
	GetRootDir() string
	GetSourceDirs() []string
	GetVersion() (*version.Version, error)
	GetLicense() (string, error)
	GetShopwellVersionConstraint() (*version.Constraints, error)
	GetType() string
	GetPath() string
	GetChangelog() (*ExtensionChangelog, error)
	GetMetaData() *ExtensionMetadata
	UpdateMetaData(*ExtensionMetadata) error
	GetExtensionConfig() *Config
	Validate(context.Context, validation.Check)
}

// getShopwellVersionConstraintFromComposer is a shared helper for composer-based extensions
// (PlatformPlugin and ShopwellBundle) to extract the Shopwell compatibility constraint from
// the composer requirements. It intentionally ignores the build-time override
// (config.Build.ShopwellVersionConstraint) so that the reported compatibility always reflects
// what the extension actually supports (e.g. for account store uploads). Use
// GetShopwellBuildVersionConstraint / GetShopwellVersionConstraintForBuild when the build-time
// override should be honored.
func getShopwellVersionConstraintFromComposer(composerRequire map[string]string) (*version.Constraints, error) {
	shopwellConstraintString, ok := composerRequire["shopwell/core"]
	if !ok {
		return nil, errors.New("require.shopwell/core is required")
	}

	shopwellConstraint, err := version.NewConstraint(shopwellConstraintString)
	if err != nil {
		return nil, err
	}

	return &shopwellConstraint, nil
}

// GetShopwellBuildVersionConstraint returns the Shopwell version constraint configured as a
// build-time override via config.Build.ShopwellVersionConstraint. It handles a nil config
// gracefully and returns (nil, nil) when no override is configured. This is the single source
// of truth for the build-time constraint and is never mixed into an extension's reported
// compatibility constraint.
func GetShopwellBuildVersionConstraint(config *Config) (*version.Constraints, error) {
	if config != nil && config.Build.ShopwellVersionConstraint != "" {
		constraint, err := version.NewConstraint(config.Build.ShopwellVersionConstraint)
		if err != nil {
			return nil, err
		}

		return &constraint, nil
	}

	return nil, nil //nolint:nilnil // nil constraint signals "no build override configured", not an error
}

// GetShopwellVersionConstraintForBuild resolves the constraint that should be used when building
// assets for the given extension. It honors the build-time override
// (config.Build.ShopwellVersionConstraint) when set and otherwise falls back to the extension's
// reported compatibility constraint.
func GetShopwellVersionConstraintForBuild(ext Extension) (*version.Constraints, error) {
	buildConstraint, err := GetShopwellBuildVersionConstraint(ext.GetExtensionConfig())
	if err != nil {
		return nil, err
	}

	if buildConstraint != nil {
		return buildConstraint, nil
	}

	return ext.GetShopwellVersionConstraint()
}
