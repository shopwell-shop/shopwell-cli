package shop

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyim/go-composer"
)

// DetectShopwellVersion returns the installed Shopwell version, falling back to
// composer.lock when no installed version can be read. It returns an empty
// string if neither file contains a Shopwell version.
func DetectShopwellVersion(projectRoot string) string {
	if installed := installedShopwellVersion(projectRoot); installed != "" {
		return strings.TrimPrefix(installed, "v")
	}

	lock, err := composer.ReadLock(filepath.Join(projectRoot, "composer.lock"))
	if err != nil {
		return ""
	}
	for _, name := range []string{"shopwell/core", "shopwell/platform"} {
		if pkg := lock.GetPackage(name); pkg != nil && pkg.Version != "" {
			return strings.TrimPrefix(pkg.Version, "v")
		}
	}
	return ""
}

func installedShopwellVersion(projectRoot string) string {
	data, err := os.ReadFile(filepath.Join(projectRoot, "vendor", "composer", "installed.json"))
	if err != nil {
		return ""
	}
	type installedPackage struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	var installed struct {
		Packages []installedPackage `json:"packages"`
	}
	if err := json.Unmarshal(data, &installed); err != nil {
		if err := json.Unmarshal(data, &installed.Packages); err != nil {
			return ""
		}
	}
	for _, name := range []string{"shopwell/core", "shopwell/platform"} {
		for _, pkg := range installed.Packages {
			if pkg.Name == name && pkg.Version != "" {
				return pkg.Version
			}
		}
	}
	return ""
}
