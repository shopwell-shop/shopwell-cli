package extension

import (
	"os"
	"path/filepath"
	"strings"
)

func PlatformPath(projectRoot, component, subPath string) string {
	return filepath.Join(projectRoot, PlatformRelPath(projectRoot, component, subPath))
}

func PlatformRelPath(projectRoot, component, subPath string) string {
	if _, err := os.Stat(filepath.Join(projectRoot, "src", "Core", "composer.json")); err == nil {
		return filepath.Join("src", component, subPath)
	} else if _, err := os.Stat(filepath.Join(projectRoot, "vendor", "shopwell", "platform")); err == nil {
		return filepath.Join("vendor", "shopwell", "platform", "src", component, subPath)
	}

	return filepath.Join("vendor", "shopwell", strings.ToLower(component), subPath)
}

// projectRequiresBuild checks if the project is a contribution project aka shopwell/shopwell.
func projectRequiresBuild(projectRoot string) bool {
	// We work inside Shopwell itself
	if _, err := os.Stat(filepath.Join(projectRoot, "src", "Core", "composer.json")); err == nil {
		return true
	}

	// vendor/shopwell/platform does never have assets pre-build
	if _, err := os.Stat(filepath.Join(projectRoot, "vendor", "shopwell", "platform", "composer.json")); err == nil {
		return true
	}

	return false
}
