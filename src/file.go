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
	ModTime  time.Time `json:"modTime"`
	SModTime string    `json:"sModTime"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"isDir"`
}

//Guarda la Path como un URL en el atributo Link del objeto
func (f *File) saveLink() {
	location := url.URL{Path: filepath.Clean(f.Path)}
	f.Link = location.RequestURI()
}

func ReturnFile(path string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}

	file := File{
		Path:     path,
		Name:     info.Name(),
		ModTime:  info.ModTime(),
		SModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		Size:     info.Size(),
		IsDir:    info.IsDir(),
	}

	file.saveLink()

	return file, nil
}

//Lista todos los archivos/directorios directos de un directorio, es decir,
//lo que estan dentro de el no los que estan dentro de los subdiretorios
func ReturnFiles(root string) ([]File, error) {

	var files []File

	//Su Root no es DirToServe
	if root != DirToServe {

		//Pongo como primer File a el directorio padre (para que el usuario
		//pueda volver al directorio anterior)
		f := File{
			Path:  filepath.Dir(root),
			Name:  "...",
			IsDir: true,
		}
		f.saveLink()

		files = append(files, f)
	}

	infos, err := ioutil.ReadDir(root)
	for _, info := range infos {

		path := filepath.Join(root, info.Name())
		//Eso para evitar que agregue al directorio padre al array
		if path == root {
			break
		}

		//Si el directorio padre del fichero es diferente de root, es decir, que
		//que esta dentro de un subdiretorio de root que no lo liste
		//Solo listara los archivo y directorio que tiene directamnete root
		if filepath.Dir(path) != filepath.Clean(root) {
			break
		}

		//Creo el arcihvo y le asigno algunos valores
		f := File{
			Path:     path,
			Name:     info.Name(),
			ModTime:  info.ModTime(),
			SModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			Size:     info.Size(),
			IsDir:    info.IsDir(),
		}

		f.saveLink()

		files = append(files, f)
	}

	return files, err
}
