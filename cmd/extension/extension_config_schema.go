package extension

import (
	"github.com/spf13/cobra"

	"github.com/shopwell-shop/shopwell-cli/internal/extension"
)

var extensionConfigSchemaCmd = &cobra.Command{
	Use:   "config-schema",
	Short: "Print the JSON schema for .config/shopwell-extension.yml",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, err := cmd.OutOrStdout().Write(extension.ConfigSchema())
		return err
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionConfigSchemaCmd)
}
