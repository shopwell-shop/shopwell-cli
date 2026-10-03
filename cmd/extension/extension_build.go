package extension

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shopwell-shop/shopwell-cli/internal/executor"
	"github.com/shopwell-shop/shopwell-cli/internal/extension"
)

var extensionAssetBundleCmd = &cobra.Command{
	Use:   "build path...",
	Short: "Compile extension Administration and Storefront assets",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assetCfg := extension.AssetBuildConfig{
			ShopwellRoot: os.Getenv("SHOPWELL_PROJECT_ROOT"),
		}
		if assetCfg.ShopwellRoot != "" {
			assetCfg.Executor = executor.NewLocal(assetCfg.ShopwellRoot)
		}
		validatedExtensions := make([]extension.Extension, 0)

		for _, arg := range args {
			path, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}

			ext, err := extension.GetExtensionByFolder(cmd.Context(), path)
			if err != nil {
				return fmt.Errorf("cannot open extension: %w", err)
			}

			validatedExtensions = append(validatedExtensions, ext)
		}

		if assetCfg.ShopwellRoot != "" {
			constraint, err := extension.GetShopwellProjectConstraint(assetCfg.ShopwellRoot)
			if err != nil {
				return fmt.Errorf("cannot get shopwell version constraint from project %s: %w", assetCfg.ShopwellRoot, err)
			}
			assetCfg.ShopwellVersion = constraint
		} else {
			constraint, err := extension.GetShopwellVersionConstraintForBuild(validatedExtensions[0])
			if err != nil {
				return fmt.Errorf("cannot get shopwell version constraint: %w", err)
			}

			assetCfg.ShopwellVersion = constraint
		}

		if err := extension.BuildAssetsForExtensions(cmd.Context(), extension.ConvertExtensionsToSources(cmd.Context(), validatedExtensions), assetCfg); err != nil {
			return fmt.Errorf("cannot build assets: %w", err)
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionAssetBundleCmd)
}
