package main

import (
	"ServerFile/src/myfuncs"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

//Ruta predeterminada del directorio de Templates
const (
	defaultTemplateDir string = "./static"
	defaultDirToServer string = "./"
)

var (
	//Puerto en el cual sirve el servidor
	PortSelected string

	//Ruta que sirve el servidor
	DirToServe string

	//Directorio de templates usara brindar el HTML/CSS/JS
	TemplateDirSeleceted string

	//Tiempo que el servidor estara abierto
	DurationTimeOpened time.Duration

	//Si el servidor servira las carpetas dentro del directorio servido o solo los archivos
	//RecursiveMode bool

	//Si se pidio el mensaje de ayuda
	HelpMessage bool
)

func init() {

	flag.BoolVar(&HelpMessage, "help", false, "Muestra el mensaje de ayuda")
	flag.StringVar(&DirToServe, "d", defaultDirToServer, "Directorio que se va servir")
	flag.StringVar(&PortSelected, "p", "8081", "Puerto donde se va a servir")
	flag.StringVar(&TemplateDirSeleceted, "D", defaultTemplateDir, "Directorio donde se obtendra los archivos HTML/CCS/JS")
	flag.DurationVar(&DurationTimeOpened, "t", 0, "Cuanto tiempo estara abierto el servidor (en s/m/h)")
	//flag.BoolVar(&RecursiveMode, "r", false, "Serivira todos los archivos y directorios de todos los directorio dentro de la ruta indicada")

	//Convierto los argumentos
	flag.Parse()

	//Si se ingreso la flag de ayuda, mostrar el mensaje y cerrar
	if HelpMessage {
		flag.Usage()
		os.Exit(0)
	}

	//Si el usuario ingresa algun directorio(osea que DirToServe no esta vacia) compruebo
	//que el directorio ingresa es una ruta absoluta, sino cierro el programa.
	if DirToServe != defaultDirToServer {
		if err := myfuncs.EsAbsolutaYExite(&DirToServe); err != nil {
			log.Fatalln(myfuncs.PrimeraMayus(err.Error()))
		}
	} else {
		//En caso de que no haya ingresado busco una ruta local
		localPaths, err := os.Getwd()
		if err != nil {
			log.Fatalln("Error al buscar el directorio actual: " + err.Error())
		}
		DirToServe = filepath.Clean(localPaths)

		//Busco la info del directorio
		i, err := os.Stat(DirToServe)

		//Si la ruta no exite
		if os.IsNotExist(err) {
			log.Fatalf("La ruta a servir(\"%s\") no exite: %s", DirToServe, err.Error())
		}

		if err != nil {
			log.Fatalf("Ocurrio un error con la ruta a servir (\"%s\"): %s", DirToServe, err.Error())
		}

		//Si no es un directorio
		if !i.IsDir() {
			log.Fatalf("La ruta a servir(\"%s\") debe ser un directorio", DirToServe)
		}
	}

	//Si timeOpen es diferente de 0, es decir, que se ingreso algun valor
	if DurationTimeOpened != 0 {

		//Comrpuebo que el valor no sea un valor negativo
		if DurationTimeOpened <= 0 {
			log.Fatalln("Se debe ingresar un tiempo de cierre valido (un valor positivo)")
		}

		/*
			Y en una goroutine espero ese tiempo, con time.Sleep(),
			y cierro el programa
		*/
		go func() {
			time.Sleep(DurationTimeOpened)
			os.Exit(0)
		}()
	}

	//Checkeo que la ruta de los templates (haya cambiado o no)
	//exista y sea una ruta absoluta
	if err := myfuncs.EsAbsolutaYExite(&TemplateDirSeleceted); err != nil {
		log.Fatalln(myfuncs.PrimeraMayus(err.Error()))
	}

	//Una vez checkeadas todas las flags
	//Muestro sus estados
	log.Println("Flags:")
	log.Println("- Servidor abierto en:", PortSelected)
	log.Println("- Ruta servida:", DirToServe)
	//log.Println("- Modo recursivo:", RecursiveMode)
	log.Println("- Tiempo de cierre automatico:", DurationTimeOpened.String())
	log.Printf("- Se esta usando el directorio \"%s\" para templates\n", TemplateDirSeleceted)
}

func main() {

	//Activar el release mode
	gin.SetMode(gin.ReleaseMode)

	//Creo el Router de Rutas
	router := gin.Default()

	//Sirvo los archivos JS, CSS y HTML
	router.StaticFS("/static", http.Dir(TemplateDirSeleceted))

	//Agrego todos los handlers
	//Si el usuario quiero ir a "/" lo rediriga a donde estan los archivos
	router.GET("/", RedirectToFiles)

	//Creo los

	api := router.Group("/api", gin.BasicAuth(gin.Accounts{
		"joaco": "joaco",
	}))
	{
		//Uso "*file" para represntar toda la ruta, ejemplo en "/dir/file1" el
		//parametro "file" sera "dir/file1"
		api.GET("/getfiles/*path", ServirArchivos)

		//Aqui donde se resibe la peticion de borrar archivos
		api.DELETE("/removefiles/*file", BorrarArchivo)

		//Aqui es donde se solicita descargar
		api.GET("/downloadfiles/*file", DescargarArchivos)

		//Aqui es donde se resiben las peticiones con archivos para subir al servidor
		api.POST("/uploadfiles", SubirArchivo)
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

	//Inicio el servidor en segundo plano
	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalln("Error al abri el servidor: ", err)
			}
		}
	}()

	//Creo un canal de signal
	c := make(chan os.Signal, 1)

	//Luego le "pido" que me notifique cuando se intente
	//cerrar un programa
	signal.Notify(c, os.Interrupt)

	//Si espero  que llegue la senial de cierre
	<-c
	log.Println("Apagando servidor...")

	//Pido el contexto con un DeadLine de 5s
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//Cuando termine "cierre" el contexto
	defer cancel()

	//Y apago el servidor
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Error al apagar el servidor:", err)
	}
	log.Println("Servidor apagado con exito")

}
