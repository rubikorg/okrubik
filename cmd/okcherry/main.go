package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"strings"

	"github.com/okcherry/cherry/pkg"
)

func main() {
	flag.Parse()

	var args = flag.Args()

	mainCmd := args[0]

	switch mainCmd {
	case "create":
		if len(args) == 1 {
			pkg.ErrorMsg("okcherry needs a value for folder name after create command. Example: okcherry create helloworld.")
			return
		}
		err := create(args[1])
		if err != nil {
			pkg.ErrorMsg(err.Error())
		}
		break
	case "run":
		var cherryConfig pkg.Config
		pwd, _ := os.Getwd()
		configPath := pwd + string(os.PathSeparator) + "cherry.toml"
		_, err := toml.DecodeFile(configPath, &cherryConfig)
		if err != nil {
			pkg.ErrorMsg("Bad config. Raw: " + err.Error())
			return
		}

		if len(cherryConfig.App) > 1 && len(args) == 1 {
			var appSlice []string
			for _, app := range cherryConfig.App {
				appSlice = append(appSlice, app.Name)
			}
			apps := strings.Join(appSlice, " or ")
			pkg.ErrorMsg("okcherry needs to know which app to run: " + apps)
			return
		} else {

		}

		break
	default:
		pkg.ErrorMsg("No such command")
	}

}
