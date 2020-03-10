package main

import (
	"errors"
	"os"
	"time"

	"github.com/oksketch/sketch"

	"github.com/oksketch/sketch/pkg"
)

var chcl = sketch.NewClient(BaseAssetURL, time.Second*30)

func create(projectName string) error {
	path, _ := os.Getwd()
	projPath := path + string(os.PathSeparator) + projectName

	if f, _ := os.Stat(projPath); f != nil {
		return errors.New("Folder with name " + projectName + " already exists.")
	}

	// create cache dir if not exists
	cachePath := getCacheDir()
	gsPath := cachePath + string(os.PathSeparator) + "gs.zip"
	// check if getting started zip file is present in cache dir
	if _, err := os.Stat(gsPath); os.IsNotExist(err) {
		// if not download it
		gsFileEn := sketch.DownloadRequestEntity{
			TargetFilePath: gsPath,
		}.Route(GSFile)

		err := chcl.Download(gsFileEn)

		if err != nil {
			return err
		}
	}

	// unzip base template to project path
	unzipFileFromCache(GSFile, projPath)

	pkg.CherryMsg("Created cherry project " + projectName + ". Happy picking!")

	return nil
}

func run(projectName string) {

}
