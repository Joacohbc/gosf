/*
Copyright © 2022 Joacohbc <joacog48@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

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

//Retorna un File apartir de una ruta. Error en caso de que no exista
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

	//Si el directorio a listar no es DirToServe
	/*
		Debido que no se puede acceder a un directorio por detras
		de DirToServer, solo los que estan dentro de el
	*/
	if root != DirToServe {

		//Obtengo la ruta relativa
		pathRelative, err := filepath.Rel(DirToServe, root)
		if err != nil {
			return nil, err
		}

		//Y agrego el Directorio anterior al que se pide
		anteriorDir := File{
			//La ruta sera el padre del directorio que se pidio
			Path: filepath.Dir(pathRelative),
			//Y el nombre sera el nombre del padre del directorio que se pidio
			Name:  "...",
			IsDir: true,
		}
		anteriorDir.saveLink()

		files = append(files, anteriorDir)
	}

	//Leo todos los archivos y subdirectorio del Directorio pedido
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		return []File{}, err
	}

	for _, info := range infos {

		//Creo la ruta del archivo/subdirectorio apartir de la ruta
		//de su directorio padre y su nombre
		path := filepath.Join(root, info.Name())

		//Para que no se agrege al el mismo
		if root == path {
			break
		}

		//Busco la ruta relativa con respecto al directorio que se esta sirviendo
		pathRelative, err := filepath.Rel(DirToServe, path)
		if err != nil {
			return nil, err
		}

		//Creo el arcihvo y le asigno algunos valores
		f := File{
			Path:     pathRelative,
			Name:     info.Name(),
			ModTime:  info.ModTime(),
			SModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			Size:     info.Size(),
			IsDir:    info.IsDir(),
		}

		f.saveLink()

		files = append(files, f)
	}

	//Y como ultimo File agrego la ruta del directorio que se
	//del cual se esta listando

	pathRelative, err := filepath.Rel(DirToServe, root)
	if err != nil {
		return nil, err
	}

	//Y agrego el Directorio anterior al que se pide
	anteriorDir := File{
		Path: pathRelative,
	}
	anteriorDir.saveLink()

	files = append(files, anteriorDir)

	return files, err
}
