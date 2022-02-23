package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

/*
	Valida que la ruta que se envio que se envio sea valido para
	que el usuario pueda acceder a el:

	- "Limpia" la ruta
	- Valida si existe
	- Valida si se puede acceder a el
	- Valida que no sea un directorio
	- Valida que el directorio este en la carpeta servido
*/
func archivoValido(llamado, path string) (string, int, error) {

	log := func(args ...interface{}) {
		//Creo un nuevo array donde agrego [llamado]
		a := append(make([]interface{}, 0), "["+llamado+"]")

		//Y a ese array le agrego
		a = append(a, args...)
		log.Println(a...)
	}

	//La ruta del archivo que se pidio
	filePedido := filepath.Clean(path)

	log("Se intento acceder a:", filePedido)
	/*
		Si el archivo que se pide no tiene carpeta
		el servidor lo tomara como que se esta buscado
		en la carpeta servida
	*/
	if s := filepath.Dir(filePedido); s == "/" {
		filePedido = filepath.Join(DirToServe, filePedido)
		log("Se le agrego el directorio al fichero:", filePedido)
	}

	//Veo la info dela arch
	info, err := os.Stat(filePedido)

	//Si el archivo no exite
	if os.IsNotExist(err) {
		log("El fichero pedido no existe:", filePedido)
		return filePedido, http.StatusInternalServerError, errors.New("el fichero al que intenta acceder no existe: " + err.Error())
	}

	//Si ocurrio otro error
	if err != nil {
		return filePedido, http.StatusInternalServerError, errors.New("ocurrio un error al buscar el fichero: " + err.Error())
	}

	//Si se esta pidiendo un directorio
	if info.IsDir() {
		log("Se intento acceder a al directorio:", filePedido)
		return filePedido, http.StatusUnauthorized, errors.New("no puede acceder a un directorio")
	}

	/*
		Verifico si el archivo existe en la carpeta que se esta sirviendo.
		Sino lo notifico
	*/
	if _, err := filepath.Rel(DirToServe, filePedido); err != nil {
		log("Se intento acceder a un archivo fuera de la carpeta:", filePedido)
		return filePedido, http.StatusUnauthorized, errors.New("no puede acceder a un archivo que no exite en la carpeta servida")
	}

	return filePedido, http.StatusOK, nil
}

// ServirArchivos - GET - /getfiles/*files
func ServirArchivos(c *gin.Context) {

	/*
		Si "file" tiene mas de un caracter, osea que no es solo "/", entoces
		sigmifica que se pidio algun archivo.

		Entoces sirvo el archivo que pidio y retorno
	*/
	if len(c.Param("file")) > 1 {

		archivo, resp, err := archivoValido(c.FullPath(), c.Param("file"))
		//Valido la ruta
		if err != nil {
			c.HTML(resp, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": err.Error(),
			})
		}

		c.File(archivo)
		return
	}

	//Leo los archivos del directorio que se me pidio(dir)
	files, err := ReturnFiles(DirToServe)
	if err != nil {
		log.Fatal("Error al leer los archivos", err)
	}

	log.Println("Cantidad de archivos cargados:", len(files))

	//Sirvo los archivos
	c.HTML(http.StatusOK, NameTemplateHtml, gin.H{
		"Files": files,
		"Error": "",
	})
}

// DescargarArchivos - GET - /downloadfiles/*files
func DescargarArchivos(c *gin.Context) {

	archivo, resp, err := archivoValido(c.FullPath(), c.Param("file"))

	if err != nil {
		c.JSON(resp, gin.H{
			"error": err.Error(),
		})
		return
	}

	/*
		b, err := ioutil.ReadFile(filePedido)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ocurrio un error al intentar leer el archivo",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"stauts": "ok",
			"data":   b,
			"name":   filepath.Base(filePedido),
		})
	*/

	c.FileAttachment(archivo, filepath.Base(archivo))

	log.Println("Se descargo el archivo: " + archivo)
}

// BorrarArchivo - POST - /removefiles/*files
func BorrarArchivo(c *gin.Context) {

	archivo, resp, err := archivoValido(c.FullPath(), c.Param("file"))
	if err != nil {
		c.JSON(resp, gin.H{
			"error": err.Error(),
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

// RedirectToFiles - POST - /
func RedirectToFiles(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "getfiles")
}

// SubirArchivo - POST - /uploadfiles
func SubirArchivo(c *gin.Context) {

	//Leo el Form
	form, err := c.MultipartForm()

	//Si ocurrio un error al leer
	if err != nil {
		log.Println("Error al leer multipart-form:" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Recorro todos los archivos del formulario(form.File) con el nombre "fileToUpload"
	for _, file := range form.File["fileToUpload"] {
		err := c.SaveUploadedFile(file, filepath.Join(DirToServe, file.Filename))
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
