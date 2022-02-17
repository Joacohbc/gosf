package main

import (
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path    string    `json:"path"`
	Name    string    `json:"name"`
	Link    string    `json:"link"`
	Index   int       `json:"index"`
	ModTime time.Time `json:"modTime"`
}

//Guarda la Path como un URL en el atributo Link del objeto
func (f *File) saveLink() {
	location := url.URL{Path: filepath.Clean(f.Path)}
	f.Link = location.RequestURI()
}

//Retorna todos los archivos de un directorio
func ReturnFiles(root string) ([]File, error) {

	var files []File
	var i int = 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

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

		//Si es el archivo tamplate lo omito ya que no lo quiero mostrar
		if info.Name() == NameTemplateHtml {
			return nil
		}

		//Creo el arcihvo y le asigno algunos valores
		f := File{
			Path:    path,
			Name:    info.Name(),
			Index:   i,
			ModTime: info.ModTime(),
		}

		f.saveLink()
		i++

		files = append(files, f)
		return nil
	})

	return files, err
}
