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
	port     string
	dir      string
	timeOpen time.Duration
)

func init() {

	flag.StringVar(&dir, "d", "", "Directorio que se va servir")
	flag.StringVar(&port, "p", "8081", "Puerto donde se va a servir")
	flag.DurationVar(&timeOpen, "t", 0, "Cuanto timepo estara abierto el servidor (en s/m/h")

	//Convierto los argumentos
	flag.Parse()

	//Si la path esta vacia
	if dir == "" {

		//Identifico la carpeta actual
		localPath, err := os.Getwd()
		if err != nil {
			log.Fatal("Error al buscar el directorio actual: " + err.Error())
		}

		//Y la guardo en "dir"
		dir = localPath
	}

	//Busco la info del directorio
	i, err := os.Stat(dir)

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

	//Si timeOpen es diferente de 0, es decir, que se ingreso algun valor
	if timeOpen != 0 {

		//Le notifico en cuanto se apagara el servidor
		log.Println("El servidor se cerrara automaticamente en:", timeOpen.String())

		/*
			Y en una goroutine espero ese tiempo, con time.Sleep(),
			y cierro el programa
		*/
		go func() {
			time.Sleep(timeOpen)
			os.Exit(0)
		}()
	}
}

func main() {

	//Activar el release mode
	gin.SetMode(gin.ReleaseMode)

	//Creo el Router de Rutas
	router := gin.Default()

	//Cargo los templates
	router.LoadHTMLGlob(PathTempalteHtml)

	//Sirvo los archivos JS, CSS y HTML
	router.StaticFS("/static", http.Dir(RootDirTemplateFiles))

	//Si el usaurio quiero ir a "/" lo rediriga a donde estan los archivos
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "getfiles")
	})

	//Uso "*file" para represntar toda la ruta, ejemplo en "/dir/file1" el parametro "file" sera "dir/file1"
	router.GET("/getfiles/*file", ServirArchivos)

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

// ServirArchivos - GET - /getfiles
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
			c.HTML(http.StatusOK, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "El archivo que ha pedido no existe",
			})
			return
		}

		//Si ocurrio otro error
		if err != nil {
			c.HTML(http.StatusOK, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "Ocurrio un error al cargar el archivo" + err.Error(),
			})
			return
		}

		//Si se esta pidiendo un directorio
		if info.IsDir() {
			c.HTML(http.StatusOK, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "No puede pedir un directorio",
			})
			return
		}

		/*
			Verifico si el archivo existe en la carpeta que se esta sirviendo.
			Sino lo notifico
		*/

		if _, err := filepath.Rel(dir, filePedido); err != nil {
			c.HTML(http.StatusOK, NameTemplateHtml, gin.H{
				"Files": []File{},
				"Error": "No puede pedir una archivo que no exite en la carpeta servida",
			})
			return
		}

		c.File(filePedido)
		return
	}

	//Leo los archivos del directorio que se me pidio(dir)
	files, err := ReturnFiles(dir)
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
