package shop

import (
	"os"

	"github.com/shyim/go-composer"
)

// ReadComposerAuth reads a Composer auth.json (via go-composer), then merges
// Shopwell-specific environment configuration on top: the COMPOSER_AUTH env var
// (handled by go-composer's MergeEnv) and the SHOPWELL_PACKAGES_TOKEN convenience
// variable, which registers a bearer token for packages.shopwell.cn.
//
// A missing auth.json is not an error; an empty, env-populated Auth is returned.
func ReadComposerAuth(authFile string) (*composer.Auth, error) {
	auth, err := composer.ReadAuth(authFile)
	if err != nil {
		return nil, err
	}

	if err := auth.MergeEnv(); err != nil {
		return nil, err
	}

	if token := os.Getenv("SHOPWELL_PACKAGES_TOKEN"); token != "" {
		if auth.BearerAuth == nil {
			auth.BearerAuth = map[string]string{}
		}
		auth.BearerAuth["packages.shopwell.cn"] = token
	}

	return auth, nil
}
