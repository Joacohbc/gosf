package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Link     string    `json:"link"`
	Index    int       `json:"index"`
	ModTime  time.Time `json:"modTime"`
	SModTime string    `json:"sModTime"`
	Size     int64     `json:"size"`
}

func ReturnFile(path string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}

	file := File{
		Path:     path,
		Name:     info.Name(),
		Index:    0,
		ModTime:  info.ModTime(),
		SModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		Size:     info.Size(),
	}

	file.saveLink()

	return file, nil
}

//Guarda la Path como un URL en el atributo Link del objeto
func (f *File) saveLink() {
	location := url.URL{Path: filepath.Clean(f.Path)}
	f.Link = location.RequestURI()
}

//Retorna todos los archivos de un directorio
func ReturnFiles(root string) (files []File, err error) {

	if RecursiveMode {
		var i int = 0
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			//Eso para evitar que agregue al directorio padre al array
			if path == root {
				return nil
			}

			//Si es un directorio lo omito
			if info.IsDir() {
				return nil
			}

			//Creo el arcihvo y le asigno algunos valores
			f := File{
				Path:     path,
				Name:     info.Name(),
				Index:    i,
				ModTime:  info.ModTime(),
				SModTime: info.ModTime().Format("2006-01-02 15:04:05"),
				Size:     info.Size(),
			}

			f.saveLink()
			i++

			files = append(files, f)
			return nil
		})
		return files, err

	} else {

		infos, err := ioutil.ReadDir(DirToServe)
		for i, info := range infos {

			if info.IsDir() {
				continue
			}

			f := File{
				Path:     filepath.Join(DirToServe, info.Name()),
				Name:     info.Name(),
				Index:    i,
				ModTime:  info.ModTime(),
				SModTime: info.ModTime().Format("2006-01-02 15:04:05"),
				Size:     info.Size(),
			}
			f.saveLink()
			files = append(files, f)
		}
		return files, err
	}
}
