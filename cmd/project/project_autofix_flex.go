package project

import (
	"errors"
	"fmt"
	"os"
	"path"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"

	"github.com/shopwell-shop/shopwell-cli/internal/flexmigrator"
	"github.com/shopwell-shop/shopwell-cli/internal/shop"
	"github.com/shopwell-shop/shopwell-cli/internal/system"
	"github.com/shopwell-shop/shopwell-cli/internal/tui"
)

var projectAutofixFlexCmd = &cobra.Command{
	Use:   "flex",
	Short: "Migrate a Shopwell project to Symfony Flex",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := shop.FindClosestShopwellProject(false)
		if err != nil {
			return err
		}

		confirmed := !system.IsInteractionEnabled(cmd.Context())

		if !confirmed {
			if err := huh.NewConfirm().
				Title("Are you sure you want to autofix this project to Symfony Flex?").
				Description("This will modify your composer.json and .env files. Make sure to commit your changes before running this command.").
				Value(&confirmed).
				Run(); err != nil {
				return err
			}
		}

		if !confirmed {
			return errors.New("autofix cancelled")
		}

		if _, err := os.Stat(path.Join(project, "symfony.lock")); err == nil {
			return errors.New("symfony.lock already exists, is that project already migrated to Symfony Flex?")
		}

		if err := flexmigrator.MigrateComposerJson(project); err != nil {
			return err
		}

		if err := flexmigrator.MigrateEnv(project); err != nil {
			return err
		}

		if err := flexmigrator.Cleanup(project); err != nil {
			return err
		}

		fmt.Println("Project migrated to Symfony Flex")
		fmt.Printf("Please run %s to install the new dependencies\n", tui.GreenText.Render("composer update"))
		fmt.Printf("and %s to apply the recipes\n", tui.GreenText.Render("yes | composer recipes:install --reset --force"))

		return nil
	},
}

func init() {
	projectAutofixCmd.AddCommand(projectAutofixFlexCmd)
}
