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

	"github.com/AlecAivazis/survey/v2"
)

var t = tint.Init()

var createQuestions = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "Project Name?"},
		Validate: survey.Required,
	},
	{
		Name: "modulePath",
		Prompt: &survey.Input{
			Message: "Init go.mod with path?",
			Help: "Keeping this blank will use the go.mod module " +
				"directive same as project name",
		},
		Validate: survey.Required,
	},
	// {
	// 	Name: "type",
	// 	Prompt: &survey.Select{
	// 		Message: "Project Type?",
	// 		Options: []string{"server"},
	// 	},
	// },
	{
		Name: "port",
		Prompt: &survey.Input{
			Message: "Default server port?",
			Default: "7000",
		},
	},
	{
		Name: "done",
		Prompt: &survey.Confirm{
			Default: true,
			Message: "Confirm?",
			Help:    "Start with rubik's development by confirming",
		},
	},
}

func prompts() (entity.CreateBoilerplateEntity, error) {
	var cbe entity.CreateBoilerplateEntity
	err := survey.Ask(createQuestions, &cbe)

	if err != nil {
		return entity.CreateBoilerplateEntity{}, err
	}

	return cbe, nil
}

// Create command main method of the okrubik cli
func Create() error {
	// ask necessary questions
	cbe, err := prompts()
	if err != nil {
		return err
	}

	var files map[string]string
	cbe.PointTo = "/boilerplate/create"
	cbe.Infer = &files

	_, err = rubcl.Get(cbe)
	if err != nil {
		pkg.ErrorMsg("Error while requesting boilerplate for rubik")
		return err
	}

	// check if project dir exists
	basePath := filepath.Join(".", cbe.Name, "cmd", "server")
	if f, _ := os.Stat(basePath); f != nil {
		return errors.New("Folder with same project name exists")
	}

	os.MkdirAll(basePath, 0755)

	for name, content := range files {
		var truePath string
		namePath := strings.Split(name, "-")
		if strings.Contains(name, "rubik.toml") {
			truePath = filepath.Join(".", cbe.Name)
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
	os.Chdir(filepath.Join(".", cbe.Name))
	cmd := exec.Command("go", "mod", "init", cbe.ModulePath)
	tidyCmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	tidyCmd.Stdout = os.Stdout

	cmd.Run()
	creationOutput("create", "go.mod")

	tidyCmd.Run()
	creationOutput("tidy", cbe.Name)

	fmt.Println("Done! Run command: okrubik run")

	return nil
}

func creationOutput(typ, path string) {
	msg := fmt.Sprintf("@(%s) %s", typ, path)
	op := t.Exp(msg, tint.Green)
	fmt.Println(op)
}
