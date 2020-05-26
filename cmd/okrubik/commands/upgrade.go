package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rubikorg/rubik/pkg"
)

// Update is ran when `okrubik update` is executed
func Upgrade(args []string) error {
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
			"https://raw.githubusercontent.com/rubikorg/okrubik/master/install", "-o",
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
