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

	info, err := os.Stat(RootDirTemplateFiles)

	if os.IsNotExist(err) {
		return errors.New("no existe el directorio de templates: " + RootDirTemplateFiles)
	}

	if err != nil {
		return errors.New("ocurrio un error al validar el archivo de templates: " + err.Error())
	}

	if !info.IsDir() {
		return errors.New("el fichero de templates debe ser un directorio")
	}
	return nil
}
