package commands

import (
	"fmt"

	"github.com/printzero/tint"
	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
)

var outBundlePath string
var runService string
var dockerFlag bool

func initBundleCommand() *cobra.Command {
	bundleCommand := &cobra.Command{
		Use:   "bundle",
		Short: "Create/Manage release bundle of your Rubik service",
		Run: func(cmd *cobra.Command, args []string) {
			err := bundle(args)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	// bundleCommand.Flags().StringVarP(&genService, "gen", "g", "",
	// 	"use this flag to generate a docker files for your Rubik server")
	bundleCommand.Flags().StringVarP(&runService, "run", "r", "",
		"use this flag to start a Rubik server")
	bundleCommand.Flags().StringVarP(&outBundlePath, "out", "o", "",
		"use this flag to specify the out directory of final Rubik service deployment")
	bundleCommand.Flags().BoolVarP(&dockerFlag, "docker", "d", false,
		"use this flag to generate docker files for Rubik service")
	return bundleCommand
}

func bundle(args []string) error {
	if len(args) == 0 && outBundlePath == "" && !dockerFlag {
		msg := t.Exp("bundle command needs a service name or flags take a look at"+
			" @(okrubik bundle --help) for more information", tint.BgMagenta.Add(tint.White.Bold()))
		fmt.Println(msg)
		return nil
	}

	if len(args) > 0 && dockerFlag {
		return generateDockerFiles()
	} else if len(args) > 0 && outBundlePath != "" && !dockerFlag {
		return buildService(true)
	} else {
		return buildService(false)
	}
}

func generateDockerFiles() error {
	return nil
}

func buildService(isCustomPath bool) error {
	return nil
}
