package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port string
	dir  string
)

func init() {
	flag.StringVar(&dir, "d", "./", "Directorio que se va servir")
	flag.StringVar(&port, "p", "8081", "Puerto donde se va a servir")

	//Convierto los argumentos
	flag.Parse()

	i, err := os.Stat(dir)

	//Si la ruta no exite
	if os.IsNotExist(err) {
		log.Fatal("La ruta ingresada no exite: ", err)
	}

	//Si no es un directorio
	if !i.IsDir() {
		log.Fatal("La ruta ingresada debe ser un directorio")
	}

	if err := CreateTemplate(); err != nil {
		log.Fatal("Error al crear el template: " + err.Error())
	}
}

func main() {

	//Activar el release mode
	gin.SetMode(gin.ReleaseMode)

	//Creo el Router de Rutas
	router := gin.Default()

	//Cargo los templates
	router.LoadHTMLGlob(PathHtml)

	//Uso "*file" para represntar toda la ruta, ejemplo en "/dir/file1" el parametro "file" sera "dir/file1"
	router.GET("/*file", func(c *gin.Context) {

		/*
			Si "file" tiene mas de un caracter, osea que no es solo "/", entoces
			sigmifica que se pidio algun archivo.

			Entoces sirvo el archivo que pidio y retorno
		*/
		if len(c.Param("file")) > 1 {
			log.Println("Archivo pedido:", c.Param("file"))
			c.File(filepath.Join(c.Param("file")))
			return
		}

		//Leo los archivos del directorio que se me pidio(dir)
		files, err := ReturnFiles(dir)
		if err != nil {
			log.Fatal("Error al leer los archivos", err)
		}

		log.Println("Cantidad de archivos cargados:", len(files))

		//Sirvo los archivos
		c.HTML(http.StatusOK, "template.html", gin.H{
			"Files": files,
		})
	})

	//Personalizo el servidor
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 10 << 20,
	}

	//Abro el servidor
	log.Println("Servidor abierto en:", port)
	log.Println("Ruta servida:", dir)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("Error al abri el servidor: ", err)
	}
}
