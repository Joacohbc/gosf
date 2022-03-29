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
	"errors"
	"fmt"
	"gosf/src/myfuncs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

//Valida que la ruta que se envio que se envio sea valido para
//que el usuario pueda acceder a el:
//- Valida si existe la ruta
//- Valida si se puede acceder a el archivo (Si esta dentro de la carpeta servida, usando rutas relativas)
//- Valida que nosea el directorio de tempalates
func archivoValido(llamado, path string) (string, int, error) {

	var mensajes []string
	defer func() {
		log.Println("\n["+llamado+"]", " - ", path, "\n", strings.Join(mensajes, "\n"))
	}()

	//Agrega mensajes a la varaible que se imprimira
	log := func(args ...interface{}) {
		mensajes = append(mensajes, fmt.Sprint(args...))
	}

	//La ruta del archivo que se pidio
	//Agrego el [1:] para que no tome el primer caracter que es "/"
	/*
		Simpre tiene un "/" porque se usan URls como "/getfile/*file"
		que sin importar si "*file" tiene o no tiene nada simpre se le
		agrega un "/"
	*/

	//Obtengo la ruta "absoluta", uniendo el Directorio a servir y
	//y la ruta pedida
	var pathAbs string = filepath.Join(DirToServe, path[1:])

	log("Se intento acceder a:", pathAbs)

	//Veo la info dela arch
	_, err := os.Stat(pathAbs)

	//Si el archivo no exite
	if os.IsNotExist(err) {
		log("El fichero pedido no existe:", pathAbs)
		return "", http.StatusInternalServerError, errors.New("el fichero al que intenta acceder no existe: " + filepath.Base(pathAbs))
	}

	//Si ocurrio otro error
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("ocurrio un error al buscar el fichero: " + filepath.Base(pathAbs))
	}

	//Verifico que el archivo pedido no pertenesca al directorio que Templates
	if strings.Contains(pathAbs, TemplateDirSeleceted) {
		log("Se intento acceder a un archivo del directorio de templates:", pathAbs)
		return "", http.StatusUnauthorized, errors.New("no puede acceder a un archivo del directorio de templates")
	}

	return pathAbs, http.StatusOK, nil
}

// RedirectToFiles - GET - /
func RedirectToFiles(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "static")
}

// ServirArchivos - GET - /api/getfiles/*path
/*
	Sirve los archivos en el directorio que se pase (/*path).
	Antes de servirlo se encarga de validar que la ruta sea un
	directorio y este dentro de DirToServe
*/
func ServirArchivos(c *gin.Context) {

	/*
		Si "path" tiene mas de un caracter, osea que no es solo "/", entoces
		sigmifica que se pidio algun archivo.

		Entoces sirvo el archivo que pidio y retorno
	*/
	if len(c.Param("path")) > 1 {

		//Obtengo la ruta del archivo validado
		path, code, err := archivoValido(c.FullPath(), c.Param("path"))
		if err != nil {
			c.JSON(code, gin.H{
				"error": myfuncs.PrimeraMayus(err.Error()),
			})
			return
		}

		//Obtengo el File de esa ruta
		file, err := ReturnFile(path)
		if err != nil {
			c.JSON(code, gin.H{
				"error": myfuncs.PrimeraMayus(err.Error()),
			})
			return
		}

		//Si se esta pidiendo un directorio
		if file.IsDir {
			//Leo los archivos del directorio que se me pidio(dir)
			files, err := ReturnFiles(path)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("no se pudo leer los archivos: %s", err.Error()))
			}

			log.Println("Cantidad de archivos cargados:", len(files))

			//Sirvo los archivos
			c.JSON(http.StatusOK, files)
			return
		}

		//Si no es un directorio, osea es un archivo, lo sirvo
		c.File(file.Path)
		return
	}

	//Leo los archivos del directorio que se me pidio(dir)
	files, err := ReturnFiles(DirToServe)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("no se pudo leer los archivos: %s", err.Error()))
	}

	log.Println("Cantidad de archivos cargados:", len(files))

	//Sirvo los archivos
	c.JSON(http.StatusOK, files)
}

// DescargarArchivos - GET - /api/downloadfiles/*files
/*
	Envia el archivo en "forma descargable" que se le pase (/*file).
	Antes de servirlo se encarga de validar que la ruta sea un
	archivo y este dentro de DirToServe
*/
func DescargarArchivos(c *gin.Context) {

	archivo, resp, err := archivoValido(c.FullPath(), c.Param("file"))

	if err != nil {
		c.JSON(resp, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
		})
		return
	}

	file, err := ReturnFile(archivo)
	if err != nil {
		c.JSON(resp, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
		})
		return
	}

	if file.IsDir {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No se puede descargar un directorio",
		})
		return
	}

	c.FileAttachment(archivo, filepath.Base(archivo))

	log.Println("Se descargo el archivo: " + archivo)
}

// BorrarArchivo - DELETE - /api/auth/removefiles/*file
/*
	Borra el archivo que se le pase (/*file).
	Antes de servirlo se encarga de validar que la ruta y
	este dentro de DirToServe
*/
func BorrarArchivo(c *gin.Context) {

	archivo, resp, err := archivoValido(c.FullPath(), c.Param("file"))
	if err != nil {
		c.JSON(resp, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
		})
		return
	}

	err = os.Remove(archivo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Se borro con exito el archivo: " + archivo)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "Se borro el archivo exitosamente",
	})
}

// SubirArchivo - POST - /api/auth/uploadfiles/*dir
/*
	Obtiene todos los archivos que se le envien y los
	guarda en la ruta que se le pase (/*dir). Antes valida
	que la ruta este dentro DirToServe y que sea un directorio
*/
func SubirArchivo(c *gin.Context) {

	//Leo el Form
	form, err := c.MultipartForm()

	//Si ocurrio un error al leer
	if err != nil {
		log.Println("Error al leer multipart-form:" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
		})
		return
	}

	//Valido la ruta donde se guardar el archivo
	path, code, err := archivoValido(c.FullPath(), c.Param("dir"))
	if err != nil {
		log.Println("Error al obtener la ruta en donde descargar el archivo:" + err.Error())
		c.JSON(code, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
		})
		return
	}

	//Obtengo la info de ese fichero
	file, err := ReturnFile(path)
	if err != nil {
		log.Println("Error al obtener la ruta en donde descargar el archivo:" + err.Error())
		c.JSON(code, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
		})
		return
	}

	//Si el fichero no es un directorio, retorno un error
	if !file.IsDir {
		log.Println("Se intento guardar un archivo en una ruta que no es un directorio")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No puede guardar un archivo en una ruta que no es un directorio",
		})
		return
	}

	//Recorro todos los archivos del formulario(form.File) con el nombre "fileToUpload"
	for _, file := range form.File["fileToUpload"] {
		err := c.SaveUploadedFile(file, filepath.Join(path, file.Filename))
		if err != nil {
			log.Println("Error al guardar archivo: " + err.Error())
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		log.Println("Arcihvo guardado en:", filepath.Join(DirToServe, file.Filename))
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
