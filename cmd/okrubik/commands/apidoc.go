package commands

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func initApiDocCommand() *cobra.Command {
	var docCmd = &cobra.Command{
		Use:   "apidoc",
		Short: "Open the API Documentation for rubik",
		Run: func(cmd *cobra.Command, args []string) {
			err := browser.OpenURL("https://pkg.go.dev/github.com/rubikorg/rubik?tab=doc")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			}
		},
	}

	return docCmd
}
