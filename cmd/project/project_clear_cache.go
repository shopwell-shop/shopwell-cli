package project

import (
	"os"

	"github.com/spf13/cobra"

	adminSdk "github.com/shopwell-shop/shopwell-cli/internal/admin-api"
	"github.com/shopwell-shop/shopwell-cli/internal/shop"
	"github.com/shopwell-shop/shopwell-cli/logging"
)

var projectClearCacheCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clear a Shopwell project's cache locally or via Admin API",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		projectRoot, err := shop.FindClosestShopwellProject(true)
		if err != nil {
			return err
		}

		cmdExecutor, err := resolveExecutor(cmd, projectRoot)
		if err != nil {
			return err
		}

		cfg := cmdExecutor.ShopConfig()
		if cfg == nil || !shop.HasAdminAPICredentials(cfg) {
			projectRoot, err = shop.FindClosestShopwellProject(false)
			if err != nil {
				return err
			}

			logging.FromContext(cmd.Context()).Infof("Clearing cache locally")

			return os.RemoveAll(projectRoot + "/var/cache")
		}

		logging.FromContext(cmd.Context()).Infof("Clearing cache using the Admin API")

		client, err := cmdExecutor.AdminAPIClient(cmd.Context())
		if err != nil {
			return err
		}

		_, err = client.CacheManager.Clear(adminSdk.NewApiContext(cmd.Context()))

		return err
	},
}

func init() {
	projectRootCmd.AddCommand(projectClearCacheCmd)
}
