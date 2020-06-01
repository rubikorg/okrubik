package commands

import (
	"fmt"
	"time"

	"github.com/printzero/tint"
	"github.com/rubikorg/rubik"

	//"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rubcl = rubik.NewClient(BaseAssetURL, time.Second*30)

const (
	// BaseAssetURL is the base url for getting files neede for okrubik
	// BaseAssetURL = "https://rubik.ashishshekar.com"
	BaseAssetURL = "http://localhost:7000"

	// GSFile is the getting started file path
	GSFile = "/gs.zip"
)

func getCachePath() string {
	home, _ := os.UserHomeDir()
	return home + string(os.PathSeparator) + ".rubik" + string(os.PathSeparator) + "cache"
}

func getCacheDir() string {
	cachePath := getCachePath()

	// if cache folder is not there then create one
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		_ = os.MkdirAll(cachePath, os.ModePerm)
	}

	return cachePath
}

var rootCmd = &cobra.Command{
	Use:   "okrubik",
	Short: "Okrubik is a command-line manager for your Rubik server",
	Long: `

Rubik is an efficient web framework for Go that encapsulates
common tasks and functions and provides ease of REST API development.
	
Complete documentation is available at https://rubikorg.github.io`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println(t.Exp("@(Welcome to Rubik Command-Line Manager) use --help for help text", tint.Magenta))
	},
}

// Execute cobra root command
func Execute() error {
	rootCmd.AddCommand(initCreateCmd())
	rootCmd.AddCommand(initRunCmd())
	rootCmd.AddCommand(initGenCmd())
	rootCmd.AddCommand(initUpgradeCmd())

	return rootCmd.Execute()
}
