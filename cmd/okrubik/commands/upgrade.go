package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
)

func initUpgradeCmd() *cobra.Command {
	var upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the project dependencies or upgrade self",
		Run: func(cmd *cobra.Command, args []string) {
			err := upgrade(args)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	return upgradeCmd
}

// upgrade is ran when `okrubik upgrade` is executed
func upgrade(args []string) error {
	if len(args) == 0 {
		dir, _ := os.Getwd()
		gomod := filepath.Join(".", "go.mod")
		if f, _ := os.Stat(gomod); f != nil {
			runTidyCommand(dir)
		} else {
			pkg.ErrorMsg("not a Go project")
		}
	} else if args[0] == "self" {
		rubikDir := pkg.MakeAndGetCacheDirPath()
		installScriptPath := filepath.Join(rubikDir, "install")
		mainCmd := "curl"
		cmd := exec.Command(mainCmd,
			"https://rubik.ashishshekar.com/install", "-o",
			installScriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}

		shCmd := exec.Command("sh", installScriptPath)
		shCmd.Stdout = os.Stdout
		shCmd.Stderr = os.Stderr
		err = shCmd.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("no such update command")
	}
	return nil
}
