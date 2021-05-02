package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/printzero/tint"
	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
)

var t = tint.Init()
var postCreationText = `
Done! To start the server, run:

cd %s
okrubik run
`
var (
	cProjName   string
	cProjModule string
	cProjPort   string
)

func initNewCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new Rubik project",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 0 {
				pkg.ErrorMsg("new command requires a project name")
				return
			}
			err = create(entity.CreateBoilerplateEntity{
				Name:       args[0],
				ModulePath: args[0],
				Port:       "7000"})
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	createCmd.Flags().StringVarP(&cProjName,
		"name", "n", "", "use this flag to set name of the project")
	createCmd.Flags().StringVarP(&cProjModule,
		"module", "m", "", "use this flag to set module path for creation of go.mod")
	createCmd.Flags().StringVarP(&cProjPort,
		"port", "p", "", "use this flag to set the port in which rubik will run")

	return createCmd
}

// create command main method of the okrubik cli
func create(inp entity.CreateBoilerplateEntity) error {
	// ask necessary questions
	var err error
	inp.Bin = "server"
	if err != nil {
		return err
	}

	var files map[string]string
	inp.PointTo = "/boilerplate/create"
	inp.Infer = &files

	_, err = rubcl.Get(inp)
	if err != nil {
		pkg.ErrorMsg("Error while requesting boilerplate for rubik")
		return err
	}

	// check if project dir exists
	basePath := filepath.Join(".", inp.Name, "cmd", "server")
	if f, _ := os.Stat(basePath); f != nil {
		return errors.New("Folder with same project name exists")
	}

	os.MkdirAll(basePath, 0755)
	rootFiles := []string{"rubik.toml", "README.md"}
	for name, content := range files {
		var truePath string
		namePath := strings.Split(name, "-")
		if in(name, rootFiles) {
			truePath = filepath.Join(".", inp.Name)
		} else {
			truePath = basePath
		}

		for _, p := range namePath {
			// ignore file name
			if !strings.Contains(p, ".tpl") {
				truePath = filepath.Join(truePath, p)
			}
		}

		os.MkdirAll(truePath, 0755)

		file := namePath[len(namePath)-1]
		// remove .tpl suffix
		file = strings.ReplaceAll(file, ".tpl", "")
		filePath := filepath.Join(truePath, file)

		creationOutput("create", filePath)
		err := ioutil.WriteFile(filePath, []byte(content), 0755)
		if err != nil {
			return err
		}
	}

	// init go.mod file
	os.Chdir(filepath.Join(".", inp.Name))
	cmd := exec.Command("go", "mod", "init", inp.ModulePath)
	cmd.Stdout = os.Stdout

	cmd.Run()
	creationOutput("create", "go.mod")

	runTidyCommand(inp.Name)

	fmt.Printf(postCreationText+"\n", inp.Name)

	return nil
}

func runTidyCommand(name string) {
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout

	tidyCmd.Run()
	creationOutput("tidy", name)
}

func creationOutput(typ, path string) {
	msg := fmt.Sprintf("@( %s ) %s", typ, path)
	op := t.Exp(msg, tint.BgGreen.Add(tint.White.Bold()))
	fmt.Println(op)
}

func in(name string, collection []string) bool {
	for _, n := range collection {
		if n == name || strings.Contains(name, n) {
			return true
		}
	}
	return false
}
