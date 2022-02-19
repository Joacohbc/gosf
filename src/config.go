package main

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	RootDirTemplateFiles string = filepath.Clean("./static")
	PathTempalteHtml     string = filepath.Join(RootDirTemplateFiles, "template.html")
	NameTemplateHtml     string = filepath.Base(PathTempalteHtml)
)

func CheckTemplate() error {

	fileExist := func(path string) bool {
		if _, err := os.Stat(path); err != nil {
			return false
		}
		return true
	}

	if !fileExist(RootDirTemplateFiles) {
		return errors.New("no existe la carpeta de templates")
	}

	if !fileExist(PathTempalteHtml) {
		return errors.New("no existe el archivo template")
	}

	return nil
}
