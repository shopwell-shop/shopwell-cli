package project

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/shopwell-shop/shopwell-cli/internal/shop"
	"github.com/shopwell-shop/shopwell-cli/internal/shop/pluginmigrate"
	"github.com/shopwell-shop/shopwell-cli/internal/system"
	pluginmigratetui "github.com/shopwell-shop/shopwell-cli/internal/tui/pluginmigrate"
)

var projectAutofixComposerCmd = &cobra.Command{
	Use:   "composer-plugins",
	Short: "Move custom/ extensions into Composer management",
	Long: "Migrate the extensions living in custom/ under Composer management: Shopwell Store plugins are required from packages.shopwell.cn and their local copy removed, everything else is registered as a Composer path repository.\n" +
		"In a terminal this runs as an interactive wizard. With --no-interaction (or without a terminal) the migration runs headless: set SHOPWELL_PACKAGIST_TOKEN to migrate Store plugins (otherwise everything becomes a path repository) and use --dry-run to preview the plan.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, err := shop.FindClosestShopwellProject(false)
		if err != nil {
			return err
		}

		exec, err := resolveExecutor(cmd, projectRoot)
		if err != nil {
			return err
		}

		if !system.IsInteractionEnabled(cmd.Context()) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			return pluginmigrate.NewPluginMigrator(projectRoot, exec).RunHeadless(cmd.Context(), pluginmigrate.HeadlessOptions{
				Token:  os.Getenv("SHOPWELL_PACKAGIST_TOKEN"),
				DryRun: dryRun,
				Out:    cmd.OutOrStdout(),
			})
		}

		_, err = pluginmigratetui.NewApp(pluginmigratetui.Options{
			ProjectRoot: projectRoot,
			Executor:    exec,
		}).Run()
		return err
	},
}

func init() {
	projectAutofixCmd.AddCommand(projectAutofixComposerCmd)
	projectAutofixComposerCmd.Flags().Bool("dry-run", false, "In non-interactive mode, print the migration plan without modifying the project")
}
