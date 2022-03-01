package main

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	RootDirTemplateFiles string = filepath.Clean("./static")
)

func CheckTemplate() error {

	fileExist := func(path string) bool {
		if _, err := os.Stat(path); err != nil {
			return false
		}
		return true
	}

	if !fileExist(RootDirTemplateFiles) {
		return errors.New("no existe el directorio de templates: " + RootDirTemplateFiles)
	}

	return nil
}
