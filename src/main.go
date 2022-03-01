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
	PortSelected         string
	DirToServe           string
	TemplateDirSeleceted string
	DurationTimeOpened   time.Duration
	RecursiveMode        bool
)

func init() {

	flag.StringVar(&DirToServe, "d", "", "Directorio que se va servir")
	flag.StringVar(&PortSelected, "p", "8081", "Puerto donde se va a servir")
	flag.StringVar(&TemplateDirSeleceted, "D", "", "Directorio donde se obtendra los archivos HTML/CCS/JS")
	flag.DurationVar(&DurationTimeOpened, "t", 0, "Cuanto tiempo estara abierto el servidor (en s/m/h)")
	flag.BoolVar(&RecursiveMode, "r", false, "Seriva todos los archivos de todos los directorio dentro de la ruta indicada")

	//Convierto los argumentos
	flag.Parse()

	//Checkeo las Flags
	{
		/*
			Si el usuario no ingresa ningun directorios(osea que DirToServe esta vacia).
			Busco el directorio actual, os.Getwd(), y sirvo ese directorio
		*/
		if DirToServe == "" {

			//Identifico el directorio actual
			localPath, err := os.Getwd()
			if err != nil {
				log.Fatal("Error al buscar el directorio actual: " + err.Error())
			}

			//Y la guardo en "DirToServe"
			DirToServe = localPath
		}

		//Si timeOpen es diferente de 0, es decir, que se ingreso algun valor
		if DurationTimeOpened != 0 {

			//Le notifico en cuanto se apagara el servidor
			log.Println("El servidor se cerrara automaticamente en:", DurationTimeOpened.String())

			/*
				Y en una goroutine espero ese tiempo, con time.Sleep(),
				y cierro el programa
			*/
			go func() {
				time.Sleep(DurationTimeOpened)
				os.Exit(0)
			}()
		}

		/*
			Si el usuario ingresa un directorio de templates compruebo que exista.
			Y cambio el directorio predeterminado(RootDirTemplatesFiles) por ese directorio
			que ingreso el usuario
		*/
		if TemplateDirSeleceted != "" {
			if _, err := os.Stat(TemplateDirSeleceted); err != nil {
				log.Fatal("El direcotrio de template ingresado no existe")
			}

			//Cambio todas las variables referentes al templates
			RootDirTemplateFiles = filepath.Clean(TemplateDirSeleceted)

			log.Printf("Se esta usando el directorio \"%s\" para templates\n", TemplateDirSeleceted)
		}

		log.Printf("Modo recursivo: %v", RecursiveMode)
	}

	//Busco la info del directorio
	i, err := os.Stat(DirToServe)

	//Si la ruta no exite
	if os.IsNotExist(err) {
		log.Fatal("La ruta ingresada no exite: ", err)
	}

	if err != nil {
		log.Fatal("Ocurrio un error con la ruta seleccionada", err)
	}

	//Si no es un directorio
	if !i.IsDir() {
		log.Fatal("La ruta ingresada debe ser un directorio")
	}

	//Creo el Template
	if err := CheckTemplate(); err != nil {
		log.Fatal("Error: ", err.Error())
	}
}

func main() {

	//Activar el release mode
	gin.SetMode(gin.ReleaseMode)

	//Creo el Router de Rutas
	router := gin.Default()

	//Sirvo los archivos JS, CSS y HTML
	router.StaticFS("/static", http.Dir(RootDirTemplateFiles))

	//Agrego todos los handlers
	{
		//Si el usuario quiero ir a "/" lo rediriga a donde estan los archivos
		router.GET("/", RedirectToFiles)

		//Aqui donde se resibe la peticion de borrar archivos
		router.DELETE("/removefiles/*file", BorrarArchivo)

		//Uso "*file" para represntar toda la ruta, ejemplo en "/dir/file1" el
		//parametro "file" sera "dir/file1"
		router.GET("/getfiles/*file", ServirArchivos)

		router.GET("/downloadfiles/*file", DescargarArchivos)

		//Aqui es donde se resiben las peticiones con archivos para subir al servidor
		router.POST("/uploadfiles", SubirArchivo)
	}

	//Personalizo el servidor
	s := &http.Server{
		Addr:           ":" + PortSelected,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 10 << 20,
	}

	log.Println("Iniciando servidor...")

	//Abro el servidor
	log.Println("Servidor abierto en:", PortSelected)
	log.Println("Ruta servida:", DirToServe)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("Error al abri el servidor: ", err)
	}
}
