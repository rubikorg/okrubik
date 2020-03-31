package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rubikorg/rubik/replc"

	"github.com/rubikorg/okrubik/cmd/okrubik/commands"

	"github.com/rubikorg/rubik/pkg"
)

func main() {
	flag.Parse()

	var args = flag.Args()

	if len(args) > 0 {
		mainCmd := args[0]

		switch mainCmd {
		case "create":
			err := commands.Create()
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
			break
		case "run":
			err := commands.Run()
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
			break
		case "help":
			fmt.Println(replc.HelpCommand([]string{}))
		default:
			pkg.ErrorMsg("No such command")
		}

	} else {
		pwd, _ := os.Getwd()
		cfg := pkg.GetRubikConfig()
		if cfg.ProjectName == "" {
			pkg.ErrorMsg("Not a rubik project! Are you on the root of your project?")
			return
		}

		// DANGER: this is using hardcoded App[1]
		basePath := strings.Replace(cfg.App[1].Path, "./",
			pwd+"/", 1)
		path := basePath + "/main.go"
		os.Setenv("RUBIK_MODE", "repl")
		os.Chdir(basePath)
		cmd := exec.Command("go", "run", path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
		os.Unsetenv("RUBIK_MODE")
	}
}
