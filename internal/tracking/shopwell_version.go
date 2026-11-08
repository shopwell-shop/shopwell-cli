package tracking

import (
	"os"
	"regexp"
	"strings"

	"github.com/shopwell-shop/shopwell-cli/internal/shop"
)

// ShopwellVersionLookup holds one optional background lookup result.
type ShopwellVersionLookup <-chan string

// Ready consumes the version if available, without waiting for filesystem I/O.
func (lookup ShopwellVersionLookup) Ready() string {
	select {
	case version := <-lookup:
		return version
	default:
		return ""
	}
}

// StartShopwellVersion starts a best-effort lookup for eligible project commands.
func StartShopwellVersion(commandName string) ShopwellVersionLookup {
	return startShopwellVersion(commandName, projectShopwellVersion)
}

func startShopwellVersion(commandName string, lookup func() string) ShopwellVersionLookup {
	if _, disabled := os.LookupEnv("DO_NOT_TRACK"); disabled || !strings.HasPrefix(commandName, "project.") || commandName == "project.create" {
		return nil
	}
	result := make(chan string, 1)
	go func() {
		// Buffer the result so finishing never depends on a consumer.
		result <- lookup()
	}()
	return result
}

// Keep arbitrary development branch names out of telemetry.
var shopwellVersionPattern = regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+(-rc\d+)?|dev-trunk)$`)

// ReportableShopwellVersion returns version if it is a release, a release
// candidate or dev-trunk, and "" for anything else, such as other development
// branches.
func ReportableShopwellVersion(version string) string {
	if !shopwellVersionPattern.MatchString(version) {
		return ""
	}
	return version
}

func projectShopwellVersion() string {
	projectRoot, err := shop.FindClosestShopwellProject(false)
	if err != nil {
		return ""
	}
	return ReportableShopwellVersion(shop.DetectShopwellVersion(projectRoot))
}
