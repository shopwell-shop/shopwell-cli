package main

import (
	"context"
	"os"

	"github.com/shopwell-shop/shopwell-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute(context.Background()))
}
