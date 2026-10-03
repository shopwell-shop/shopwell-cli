package extension

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/shopwell-shop/shopwell-cli/logging"
)

func setupShopwellInTemp(ctx context.Context, minVersion string) (string, error) {
	dir, err := os.MkdirTemp("", "extension")
	if err != nil {
		return "", err
	}

	branch := "v" + strings.ToLower(minVersion)

	if minVersion == DevVersionNumber {
		branch = "trunk"
	}

	logging.FromContext(ctx).Infof("Cloning shopwell with branch: %s into %s", branch, dir)

	gitCheckoutCmd := exec.CommandContext(ctx, "git", "clone", "https://github.com/shopwell-shop/shopwell.git", "--depth=1", "-b", branch, dir)
	gitCheckoutCmd.Stdout = os.Stdout
	gitCheckoutCmd.Stderr = os.Stderr
	err = gitCheckoutCmd.Run()
	if err != nil {
		return "", err
	}

	return dir, nil
}
