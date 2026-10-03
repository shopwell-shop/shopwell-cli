package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/shyim/go-version"

	"github.com/shopwell-shop/shopwell-cli/internal/asset"
	"github.com/shopwell-shop/shopwell-cli/internal/ci"
	"github.com/shopwell-shop/shopwell-cli/internal/esbuild"
	"github.com/shopwell-shop/shopwell-cli/internal/npm"
	"github.com/shopwell-shop/shopwell-cli/logging"
)

func BuildAssetsForExtensions(ctx context.Context, sources []asset.Source, assetConfig AssetBuildConfig) error { // nolint:gocyclo
	cfgs := BuildAssetConfigFromExtensions(ctx, sources, assetConfig)

	if len(cfgs) == 0 {
		return nil
	}

	if err := restoreAssetCaches(ctx, cfgs, assetConfig); err != nil {
		return err
	}

	if !cfgs.RequiresAdminBuild() && !cfgs.RequiresStorefrontBuild() {
		logging.FromContext(ctx).Infof("Building assets has been skipped as not required")
		return nil
	}

	var minVersion string
	getMinVersion := func() (string, error) {
		if minVersion != "" {
			return minVersion, nil
		}

		var err error
		minVersion, err = lookupForMinMatchingVersion(ctx, assetConfig.ShopwellVersion)
		return minVersion, err
	}

	requiresShopwellSources := cfgs.RequiresShopwellRepository()

	shopwellRoot := assetConfig.ShopwellRoot
	if shopwellRoot == "" && requiresShopwellSources {
		mv, err := getMinVersion()
		if err != nil {
			return err
		}

		shopwellRoot, err = setupShopwellInTemp(ctx, mv)
		if err != nil {
			return err
		}

		assetConfig.ShopwellRoot = shopwellRoot
		defer deletePaths(ctx, shopwellRoot)
	}

	nodeInstallSection := ci.Start("Installing node_modules for extensions")

	paths, err := InstallNodeModulesOfConfigs(ctx, cfgs, assetConfig)
	if err != nil {
		return err
	}

	nodeInstallSection.End()

	if shopwellRoot != "" && len(assetConfig.KeepNodeModules) > 0 {
		paths = slices.DeleteFunc(paths, func(path string) bool {
			rel, err := filepath.Rel(shopwellRoot, path)
			if err != nil {
				return false
			}

			return slices.Contains(assetConfig.KeepNodeModules, rel)
		})
	}

	defer deletePaths(ctx, paths...)

	if !assetConfig.DisableAdminBuild && cfgs.RequiresAdminBuild() {
		administrationSection := ci.Start("Building administration assets")

		// Build all extensions compatible with esbuild first
		for name, entry := range cfgs.FilterByAdminAndEsBuild(true) {
			options := esbuild.NewAssetCompileOptionsAdmin(name, entry.BasePath)
			options.DisableSass = entry.DisableSass

			res, err := esbuild.CompileExtensionAsset(ctx, options)
			if err != nil {
				return err
			}

			if err := esbuild.DumpViteConfig(options, res); err != nil {
				return err
			}

			logging.FromContext(ctx).Infof("Building administration assets for %s using ESBuild", name)
		}

		nonCompatibleExtensions := cfgs.FilterByAdminAndEsBuild(false)

		if len(nonCompatibleExtensions) != 0 {
			if projectRequiresBuild(shopwellRoot) {
				// add the storefront itself as plugin into json
				var basePath string
				if shopwellRoot == "" {
					basePath = "src/Storefront/"
				} else {
					basePath = strings.TrimLeft(
						strings.Replace(PlatformPath(shopwellRoot, "Storefront", ""), shopwellRoot, "", 1),
						"/",
					) + "/"
				}

				storefrontEntryPath := "Resources/app/storefront/src/main.js"
				adminEntryPath := "Resources/app/administration/src/main.js"
				nonCompatibleExtensions["Storefront"] = &ExtensionAssetConfigEntry{
					BasePath:      basePath,
					Views:         []string{"Resources/views"},
					TechnicalName: "storefront",
					Storefront: ExtensionAssetConfigStorefront{
						Path:          "Resources/app/storefront/src",
						EntryFilePath: &storefrontEntryPath,
						StyleFiles:    []string{},
					},
					Administration: ExtensionAssetConfigAdmin{
						Path:          "Resources/app/administration/src",
						EntryFilePath: &adminEntryPath,
					},
				}
			}

			if err := prepareShopwellForAsset(shopwellRoot, nonCompatibleExtensions, assetConfig); err != nil {
				return err
			}

			administrationRoot := PlatformPath(shopwellRoot, "Administration", "Resources/app/administration")
			adminRelPath := PlatformRelPath(shopwellRoot, "Administration", "Resources/app/administration")

			if assetConfig.NPMForceInstall || !npm.NodeModulesExists(administrationRoot) {
				var additionalNpmParameters []string

				npmPackage, err := npm.ReadPackage(administrationRoot)
				if err != nil {
					return err
				}

				if npmPackage.HasDevDependency("puppeteer") {
					additionalNpmParameters = []string{"--production"}
				}

				if err := npm.InstallDependencies(ctx, assetConfig.ExecutorWithRelDir(adminRelPath), npmPackage, additionalNpmParameters...); err != nil {
					return err
				}
			}

			envMap := map[string]string{
				"PROJECT_ROOT": assetConfig.NormalizePath(shopwellRoot),
				"ADMIN_ROOT":   assetConfig.NormalizePath(PlatformPath(shopwellRoot, "Administration", "")),
			}

			if !projectRequiresBuild(shopwellRoot) && !assetConfig.ForceAdminBuild {
				logging.FromContext(ctx).Debugf("Building only administration assets for plugins")
				envMap["SHOPWELL_ADMIN_BUILD_ONLY_EXTENSIONS"] = "1"
				envMap["SHOPWELL_ADMIN_SKIP_SOURCEMAP_GENERATION"] = "1"
			} else {
				logging.FromContext(ctx).Debugf("Building also the administration itself")
			}

			adminExec := assetConfig.ExecutorWithRelDir(adminRelPath).WithEnv(envMap)
			npmBuild := adminExec.NPMCommand(ctx, "run", "build")
			npmBuild.Cmd.Stdout = os.Stdout
			npmBuild.Cmd.Stderr = os.Stderr
			err = npmBuild.Run()

			if assetConfig.CleanupNodeModules {
				defer deletePaths(ctx, path.Join(administrationRoot, "node_modules"), path.Join(administrationRoot, "twigVuePlugin"))
			}

			if err != nil {
				return err
			}

			for name, entry := range nonCompatibleExtensions {
				options := esbuild.NewAssetCompileOptionsAdmin(name, entry.BasePath)
				if err := esbuild.DumpViteConfig(options); err != nil {
					return err
				}
			}
		}

		administrationSection.End()
	}

	if !assetConfig.DisableStorefrontBuild && cfgs.RequiresStorefrontBuild() {
		storefrontSection := ci.Start("Building storefront assets")
		// Build all extensions compatible with esbuild first
		for name, entry := range cfgs.FilterByStorefrontAndEsBuild(true) {
			isNewLayout := false

			mv, err := getMinVersion()
			if err != nil {
				return err
			}

			if mv == DevVersionNumber || version.Must(version.NewVersion(mv)).GreaterThanOrEqual(version.Must(version.NewVersion("6.6.0.0"))) {
				isNewLayout = true
			}

			options := esbuild.NewAssetCompileOptionsStorefront(name, entry.BasePath, isNewLayout)

			if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
				return err
			}
			logging.FromContext(ctx).Infof("Building storefront assets for %s using ESBuild", name)
		}

		nonCompatibleExtensions := cfgs.FilterByStorefrontAndEsBuild(false)

		if len(nonCompatibleExtensions) != 0 {
			// add the storefront itself as plugin into json
			var basePath string
			if shopwellRoot == "" {
				basePath = "src/Storefront/"
			} else {
				basePath = strings.TrimLeft(
					strings.Replace(PlatformPath(shopwellRoot, "Storefront", ""), shopwellRoot, "", 1),
					"/",
				) + "/"
			}

			entryPath := "Resources/app/storefront/src/main.js"
			nonCompatibleExtensions["Storefront"] = &ExtensionAssetConfigEntry{
				BasePath:      basePath,
				Views:         []string{"Resources/views"},
				TechnicalName: "storefront",
				Storefront: ExtensionAssetConfigStorefront{
					Path:          "Resources/app/storefront/src",
					EntryFilePath: &entryPath,
					StyleFiles:    []string{},
				},
				Administration: ExtensionAssetConfigAdmin{
					Path: "Resources/app/administration/src",
				},
			}

			if err := prepareShopwellForAsset(shopwellRoot, nonCompatibleExtensions, assetConfig); err != nil {
				return err
			}

			storefrontRoot := PlatformPath(shopwellRoot, "Storefront", "Resources/app/storefront")
			storefrontRelPath := PlatformRelPath(shopwellRoot, "Storefront", "Resources/app/storefront")
			sfExec := assetConfig.ExecutorWithRelDir(storefrontRelPath)

			npmPackage, err := npm.ReadPackage(storefrontRoot)
			if err != nil {
				return err
			}

			if assetConfig.NPMForceInstall || !npm.NodeModulesExists(storefrontRoot) {
				if err := npm.PatchPackageLockToRemoveCanIUse(path.Join(storefrontRoot, "package-lock.json")); err != nil {
					return err
				}

				additionalNpmParameters := []string{"caniuse-lite"}

				if npmPackage.HasDevDependency("puppeteer") {
					additionalNpmParameters = append(additionalNpmParameters, "--production")
				}

				if err := npm.InstallDependencies(ctx, sfExec, npmPackage, additionalNpmParameters...); err != nil {
					return err
				}

				// As we call npm install caniuse-lite, we need to run the postinstall script manually.
				if npmPackage.HasScript("postinstall") {
					npmRunPostInstall := sfExec.NPMCommand(ctx, "run", "postinstall")
					npmRunPostInstall.Cmd.Stdout = os.Stdout
					npmRunPostInstall.Cmd.Stderr = os.Stderr

					if err := npmRunPostInstall.Run(); err != nil {
						return err
					}
				}

				if _, err := os.Stat(path.Join(storefrontRoot, "vendor/bootstrap")); os.IsNotExist(err) {
					npmVendor := sfExec.NPMCommand(ctx, "exec", "--", "node", "copy-to-vendor.js")
					npmVendor.Cmd.Stdout = os.Stdout
					npmVendor.Cmd.Stderr = os.Stderr
					if err := npmVendor.Run(); err != nil {
						return err
					}
				}
			}

			sfEnvMap := map[string]string{
				"NODE_ENV":        "production",
				"PROJECT_ROOT":    assetConfig.NormalizePath(shopwellRoot),
				"STOREFRONT_ROOT": assetConfig.NormalizePath(storefrontRoot),
			}

			if assetConfig.Browserslist != "" {
				sfEnvMap["BROWSERSLIST"] = assetConfig.Browserslist
			}

			storefrontBuildExec := sfExec.WithEnv(sfEnvMap)
			if npmPackage.HasScript("production") {
				npmProduction := storefrontBuildExec.NPMCommand(ctx, "run", "production")
				npmProduction.Cmd.Stdout = os.Stdout
				npmProduction.Cmd.Stderr = os.Stderr

				if err := npmProduction.Run(); err != nil {
					return err
				}
			} else {
				nodeWebpackCmd := storefrontBuildExec.NPMCommand(ctx, "exec", "--", "webpack", "--config", "webpack.config.js")
				nodeWebpackCmd.Cmd.Stdout = os.Stdout
				nodeWebpackCmd.Cmd.Stderr = os.Stderr

				if err := nodeWebpackCmd.Run(); err != nil {
					return err
				}
			}

			if assetConfig.CleanupNodeModules {
				defer deletePaths(ctx, path.Join(storefrontRoot, "node_modules"))
			}
		}

		storefrontSection.End()
	}

	if err := storeAssetCaches(ctx, cfgs, assetConfig); err != nil {
		return err
	}

	return nil
}

func prepareShopwellForAsset(shopwellRoot string, cfgs ExtensionAssetConfig, assetConfig AssetBuildConfig) error {
	varFolder := shopwellRoot + "/var"
	if _, err := os.Stat(varFolder); os.IsNotExist(err) {
		err := os.Mkdir(varFolder, 0o755)
		if err != nil {
			return fmt.Errorf("cannot create %s: %w", varFolder, err)
		}
	}

	normalized := make(map[string]*ExtensionAssetConfigEntry, len(cfgs))
	for name, cfg := range cfgs {
		entry := new(ExtensionAssetConfigEntry)
		*entry = ExtensionAssetConfigEntry{
			BasePath:       assetConfig.NormalizePath(cfg.BasePath),
			TechnicalName:  cfg.TechnicalName,
			Administration: cfg.Administration,
			Storefront:     cfg.Storefront,
		}
		entry.Views = make([]string, len(cfg.Views))
		for i, v := range cfg.Views {
			entry.Views[i] = assetConfig.NormalizePath(v)
		}
		normalized[name] = entry
	}

	pluginJson, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("cannot encode plugins.json: %w", err)
	}

	if err = os.WriteFile(shopwellRoot+"/var/plugins.json", pluginJson, os.ModePerm); err != nil {
		return fmt.Errorf("cannot write %s: %w", shopwellRoot+"/var/plugins.json", err)
	}

	err = os.WriteFile(shopwellRoot+"/var/features.json", []byte("{}"), 0o644)
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", shopwellRoot+"/var/features.json", err)
	}

	return nil
}
