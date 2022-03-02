package main

import (
	"ServerFile/src/myfuncs"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
	Valida que la ruta que se envio que se envio sea valido para
	que el usuario pueda acceder a el:

	- "Limpia" la ruta
	- Valida si existe
	- Valida si se puede acceder a el
	- Valida que no sea un directorio
	- Valida que el directorio pedido este en el directorio servido
	- Valida el modo recursivo
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
	//Agrego el [1:] para que no tome el primer caracter que es "/"
	/*
		Simpre tiene un "/" porque se usan URls como "/getfile/*file"
		que sin importar si "*file" tiene o no tiene nada simpre se le
		agrega un "/"
	*/
	filePedido := filepath.Clean(path[1:])

	log("Se intento acceder a:", filePedido)

	/*
		Si el archivo que se pide no tiene directorio
		el servidor lo tomara como que se esta buscado
		en el directorio servido
	*/
	//if s := filepath.Dir(filePedido); s == "/" {
	//	filePedido = filepath.Join(DirToServe, filePedido)
	//	log("Se le agrego el directorio al fichero:", filePedido)
	//}

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
		Verifico si el archivo existe en el directorio que se esta sirviendo.
		Sino lo notifico.

		Compruebo checkeando que si en la ruta que se esta pidiendo
		contiene el directorio servido, si no lo contiene sigmifica
		que no esta dentro de el. (Al estar trabajando con rutas absolutas
		y no relativas esto funciona)
	*/
	if !strings.Contains(filePedido, DirToServe) {
		log("Se intento acceder a un archivo fuera del directorio:", filePedido)
		return filePedido, http.StatusUnauthorized, errors.New("no puede acceder a un archivo que no exite en el directorio servida")
	}

	//Verifico que el archivo pedido no pertenesca al directorio que Templates
	if strings.Contains(filePedido, RootDirTemplateFiles) {
		log("Se intento acceder a un archivo del directorio de templates:", filePedido)
		return filePedido, http.StatusUnauthorized, errors.New("no puede acceder a un archivo del directorio de templates")
	}

	//Si el modo recursivo no esta activado
	if !RecursiveMode {
		//Los archivos que se pidan deben tener como padre
		//estrictamente al directorio servido
		if filepath.Dir(filePedido) != DirToServe {
			log("Se intento acceder a un archivo de un directorio no permitido (Modo Recursivo off):", filePedido)
			return filePedido, http.StatusUnauthorized, errors.New("no puede acceder a un archivo dentro de este directorio sin el modo recursivo activado")
		}
	}

	return filePedido, http.StatusOK, nil
}

// RedirectToFiles - GET - /
func RedirectToFiles(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "static")
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
			c.JSON(resp, gin.H{
				"error": myfuncs.PrimeraMayus(err.Error()),
			})
			return
		}

		//blob, err := ioutil.ReadFile(archivo)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{
		//		"error": myfuncs.PrimeraMayus(err.Error()),
		//	})
		//}

		//c.JSON(http.StatusOK, gin.H{
		//	"file": blob,
		//	"type": http.DetectContentType(blob),
		//	"name": filepath.Base(archivo),
		//})

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
	c.JSON(http.StatusOK, files)
}

// DescargarArchivos - GET - /downloadfiles/*files
func DescargarArchivos(c *gin.Context) {

	archivo, resp, err := archivoValido(c.FullPath(), c.Param("file"))

	if err != nil {
		c.JSON(resp, gin.H{
			"error": myfuncs.PrimeraMayus(err.Error()),
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

// BorrarArchivo - DELETE - /removefiles/*files
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

// SubirArchivo - POST - /uploadfiles
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
