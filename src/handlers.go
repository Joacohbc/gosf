package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ServirArchivos - GET - /getfiles/*files
func ServirArchivos(c *gin.Context) {

	/*
		Si "file" tiene mas de un caracter, osea que no es solo "/", entoces
		sigmifica que se pidio algun archivo.

		Entoces sirvo el archivo que pidio y retorno
	*/
	if len(c.Param("file")) > 1 {

		//La ruta del archivo que se pidio
		filePedido := filepath.Clean(c.Param("file"))

		/*
			Si el archivo que se pide no tiene carpeta
			el servidor lo tomara como que se esta buscado en "./"
		*/
		if s := filepath.Dir(filePedido); s == "/" {
			filePedido = "./" + filepath.Clean(c.Param("file"))
		}

		log.Println("Archivo pedido:", filePedido)

		//Veo la info dela arch
		info, err := os.Stat(filePedido)

		//Si el archivo no exite
		if os.IsNotExist(err) {
			c.HTML(http.StatusInternalServerError, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "El archivo que ha pedido no existe",
			})
			return
		}

		//Si ocurrio otro error
		if err != nil {
			c.HTML(http.StatusInternalServerError, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "Ocurrio un error al cargar el archivo" + err.Error(),
			})
			return
		}

		//Si se esta pidiendo un directorio
		if info.IsDir() {
			c.HTML(http.StatusUnauthorized, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "No puede pedir un directorio",
			})
			return
		}

		/*
			Verifico si el archivo existe en la carpeta que se esta sirviendo.
			Sino lo notifico
		*/

		if _, err := filepath.Rel(DirToServe, filePedido); err != nil {
			c.HTML(http.StatusUnauthorized, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "No puede pedir una archivo que no exite en la carpeta servida",
			})
			return
		}

		c.File(filePedido)
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

	//La ruta del archivo que se pidio
	filePedido := filepath.Clean(c.Param("file"))

	/*
		Si el archivo que se pide no tiene carpeta
		el servidor lo tomara como que se esta buscado en "./"
	*/
	if s := filepath.Dir(filePedido); s == "/" {
		filePedido = "./" + filepath.Clean(c.Param("file"))
	}

	//Veo la info dela arch
	info, err := os.Stat(filePedido)

	//Si el archivo no exite
	if os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "El archivo que ha pedido no existe",
		})
		return
	}

	//Si ocurrio otro error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ocurrio un error al buscar el archivo" + err.Error(),
		})
		return
	}

	//Si se esta pidiendo un directorio
	if info.IsDir() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No puede descargar un directorio",
		})
		return
	}

	/*
		Verifico si el archivo existe en la carpeta que se esta sirviendo.
		Sino lo notifico
	*/

	if _, err := filepath.Rel(DirToServe, filePedido); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No puede descargar un archivo que no exite en la carpeta servida",
		})
		return
	}

	c.FileAttachment(filePedido, filepath.Base(filePedido))

	log.Println("Se descargo el archivo: " + filePedido)
}

// BorrarArchivo - POST - /removefiles/*files
func BorrarArchivo(c *gin.Context) {

	//La ruta del archivo que se pidio
	filePedido := filepath.Clean(c.Param("file"))

	/*
		Si el archivo que se pide no tiene carpeta
		el servidor lo tomara como que se esta buscado en "./"
	*/
	if s := filepath.Dir(filePedido); s == "/" {
		filePedido = "./" + filepath.Clean(c.Param("file"))
	}

	//Veo la info dela arch
	info, err := os.Stat(filePedido)

	//Si el archivo no exite
	if os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "El archivo que ha pedido no existe",
		})
		return
	}

	//Si ocurrio otro error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ocurrio un error al buscar el archivo" + err.Error(),
		})
		return
	}

	//Si se esta pidiendo un directorio
	if info.IsDir() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No puede eliminar un directorio",
		})
		return
	}

	/*
		Verifico si el archivo existe en la carpeta que se esta sirviendo.
		Sino lo notifico
	*/

	if _, err := filepath.Rel(DirToServe, filePedido); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No puede elimnar un archivo que no exite en la carpeta servida",
		})
		return
	}

	err = os.Remove(c.Param("file"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Se borro con exito el archivo: " + filePedido)

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
		log.Println("Error al ller multipart-form:" + err.Error())
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
