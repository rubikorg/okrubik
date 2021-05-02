package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/rubikorg/rubik"

	"github.com/rubikorg/okrubik/cmd/okrubik/commands"

	"github.com/rubikorg/rubik/pkg"
)

func main() {
	commands.Execute()

	// TODO: 1.0 REPL code
	// 	// 	pwd, _ := os.Getwd()
	// 	// 	cfg := pkg.GetRubikConfig()
	// 	// 	if cfg.ProjectName == "" {
	// 	// 		pkg.ErrorMsg("Not a rubik project! Are you on the root of your project?")
	// 	// 		return
	// 	// 	}

	// 	// 	// DANGER: this is using hardcoded App[1]
	// 	// 	basePath := strings.Replace(cfg.App[1].Path, "./",
	// 	// 		pwd+"/", 1)
	// 	// 	path := basePath + "/main.go"
	// 	// 	os.Setenv("RUBIK_MODE", "repl")
	// 	// 	os.Chdir(basePath)
	// 	// 	cmd := exec.Command("go", "run", path)
	// 	// 	cmd.Stdin = os.Stdin
	// 	// 	cmd.Stdout = os.Stdout
	// 	// 	cmd.Run()
	// 	// 	os.Unsetenv("RUBIK_MODE")
	// 	// }
}

func downloadCacheFiles() {
	file := pkg.GetErrorHTMLPath()
	if f, _ := os.Stat(file); f == nil {
		rubcl := rubik.NewClient(commands.BaseAssetURL, time.Second*30)
		en := rubik.BlankRequestEntity{}
		en.PointTo = "/boilerplate/error.html"

		resp, _ := rubcl.Get(en)
		if resp.StringBody != "" {
			err := ioutil.WriteFile(file, []byte(resp.StringBody), 0755)
			if err != nil {
				pkg.ErrorMsg("couldn't write error html cache")
			}
		}
	}
}
