package commands

import (
	"fmt"

	"github.com/rubikorg/okrubik/pkg/entity"
	"github.com/rubikorg/rubik/pkg"

	//"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/spin"
	//"os"
	//"github.com/rubikorg/rubik"
	//"github.com/rubikorg/rubik/pkg"
	//"os"
	// "time"
)

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

func prompts() error {
	var cbe entity.CreateBoilerplateEntity
	survey.Ask(createQuestions, &cbe)

	s := spin.New("%s Requesting template")
	s.Set(spin.Spin3)
	s.Start()
	defer s.Stop()

	cbe.PointTo = "/boilerplate/create"
	resp, err := rubcl.Get(cbe)
	if err != nil {
		pkg.ErrorMsg("Error while requesting boilerplate for rubik")
		return err
	}

	fmt.Println(resp.ParsedBody)

	return nil
}

// Create command main method of the okrubik cli
func Create() error {
	// ask necessary questions
	err := prompts()
	if err != nil {
		fmt.Println(err)
	}

	//path, _ := os.Getwd()
	//projPath := path + string(os.PathSeparator) + projectName
	//
	//if f, _ := os.Stat(projPath); f != nil {
	//	return errors.New("Folder with name " + projectName + " already exists.")
	//}
	//
	//// create cache dir if not exists
	//cachePath := getCacheDir()
	//gsPath := cachePath + string(os.PathSeparator) + "gs.zip"
	//// check if getting started zip file is present in cache dir
	//if _, err := os.Stat(gsPath); os.IsNotExist(err) {
	//	// if not download it
	//	gsFileEn := rubik.DownloadRequestEntity{
	//		TargetFilePath: gsPath,
	//	}.Route(GSFile)
	//
	//	err := chcl.Download(gsFileEn)
	//
	//	if err != nil {
	//		return err
	//	}
	//}

	// unzip base template to project path
	//unzipFileFromCache(GSFile, projPath)
	//
	//pkg.RubikMsg("Created rubik project " + projectName + ". Happy solving your cube!")

	return nil
}
