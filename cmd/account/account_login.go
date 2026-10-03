package account

import (
	"fmt"

	"github.com/spf13/cobra"

	accountApi "github.com/shopwell-shop/shopwell-cli/internal/account-api"
	"github.com/shopwell-shop/shopwell-cli/internal/tui"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Shopwell Account to manage extensions and credentials",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		tui.PrintBanner()

		_, err := accountApi.NewApi(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(tui.GreenText.Render("  Login successful!"))
		fmt.Println(tui.DimText.Render("  To logout, run: shopwell-cli account logout"))
		fmt.Println()

		return nil
	},
}

func init() {
	accountRootCmd.AddCommand(loginCmd)
}
