package main

import (
	"archive/zip"
	"io"
	//"net/http"
	"os"
	"path/filepath"

	"github.com/oksketch/sketch/pkg"
)

const (
	// BaseAssetURL is the base url for getting files neede for oksketch
	BaseAssetURL = "http://localhost:7000"
	// GSFile is the getting started file path
	GSFile = "/gs.zip"
)

func getCachePath() string {
	home, _ := os.UserHomeDir()
	return home + string(os.PathSeparator) + ".cherry" + string(os.PathSeparator) + "cache"
}

func getCacheDir() string {
	cachePath := getCachePath()

	// if cache folder is not there then create one
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		_ = os.MkdirAll(cachePath, os.ModePerm)
	}

	return cachePath
}

func unzipFileFromCache(fileName string, projPath string) {
	cachePath := getCachePath()
	// src
	filePath := cachePath + string(os.PathSeparator) + fileName
	// target
	err := unzip(filePath, projPath)

	if err != nil {
		pkg.ErrorMsg("Wasn't able to unzip the template because: " + err.Error())
		return
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	_ = os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, f.Mode())
		} else {
			_ = os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
