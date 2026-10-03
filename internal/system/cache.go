package system

import (
	"os"
	"path"
)

// GetShopwellCliCacheDir returns the base cache directory for shopwell-c
func GetShopwellCliCacheDir() string {
	if dir := os.Getenv("SHOPWELL_CLI_CACHE_DIR"); dir != "" {
		return dir
	}

	cacheDir, _ := os.UserCacheDir()

	return path.Join(cacheDir, "shopwell-cli")
}
