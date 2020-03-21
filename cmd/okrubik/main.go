package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
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
			var rubikConfig pkg.Config
			pwd, _ := os.Getwd()
			configPath := pwd + string(os.PathSeparator) + "rubik.toml"
			_, err := toml.DecodeFile(configPath, &rubikConfig)
			if err != nil {
				pkg.ErrorMsg("Bad config. Raw: " + err.Error())
				return
			}

			fmt.Print("hey")

			//if len(rubikConfig.App) > 1 && len(args) == 1 {
			//	var appSlice []string
			//	var lookup = make(map[string]int)
			//	for i, app := range rubikConfig.App {
			//		appSlice = append(appSlice, app.Name)
			//		lookup[app.Name] = i
			//	}
			//	prompt := promptui.Select{
			//		Label: "Run app:",
			//		Items: appSlice,
			//	}
			//
			//	_, _, _ = prompt.Run()
			//	//in := lookup[result]
			//	//cmd += sketchConfig.App[in].Path
			//	return
			//} else {
			//
			//}

			break
		default:
			pkg.ErrorMsg("No such command")
		}

	} else {
		pkg.ErrorMsg("Nothing to do")
	}
}
