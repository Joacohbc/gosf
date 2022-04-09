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
	"context"
	"fmt"
	"gosf/src/myfuncs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:        "server",
	SuggestFor: []string{"serve", "dir", "serve-dir"},
	Short:      "Inicia el servidor y sirve el directorio indicado (predeterminadamente es \"./\")",
	Args:       cobra.MaximumNArgs(1),
	PreRun:     cargarVariables,
	Run:        serverOn,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	//Flags
	serverCmd.Flags().StringVarP(&PortSelected, "port", "p", "8081", "Puerto donde se va a servir")
	serverCmd.Flags().StringVarP(&TemplateDirSeleceted, "template-directory", "D", "", "Directorio donde se obtendra los archivos HTML/CCS/JS")
	serverCmd.Flags().DurationVarP(&DurationTimeOpened, "time-live", "t", 0, "Cuanto tiempo estara abierto el servidor (en s/m/h)")
	serverCmd.Flags().StringSliceP("user", "U", []string{}, "Indica cual sera el usuario y la contraseña del administrador")
}

func serverOn(cmd *cobra.Command, args []string) {
	//Activar el release mode
	gin.SetMode(gin.ReleaseMode)

	//Creo el Router de Rutas
	router := gin.Default()

	//Sirvo los archivos JS, CSS y HTML
	router.StaticFS("/static", http.Dir(TemplateDirSeleceted))

	//Agrego todos los handlers
	//Si el usuario quiero ir a "/" lo rediriga a donde estan los archivos
	router.GET("/", RedirectToFiles)

	//En el grupo de "/api" estaran todos los hanlder referentes al funciones
	api := router.Group("/api")
	{
		//Uso "*file" para represntar toda la ruta, ejemplo en "/dir/file1" el
		//parametro "file" sera "dir/file1"
		api.GET("/getfiles/*path", ServirArchivos)

		//Aqui es donde se solicita descargar
		api.GET("/downloadfiles/*file", DescargarArchivos)

		//En "/api/auth/" estaran los handlers que se refieran a
		//la modificacion de archivos
		auth := api.Group("/auth", gin.BasicAuth(AdminAuth))
		{
			//Aqui donde se recibe la peticion de borrar archivos
			auth.DELETE("/removefiles/*file", BorrarArchivo)

			//Aqui es donde se reciben las peticiones con archivos para subir al servidor
			auth.POST("/uploadfiles/*dir", SubirArchivo)
		}
	}

	//Personalizo el servidor
	s := &http.Server{
		Addr:           ":" + PortSelected,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 10 << 20,
	}

	apagarServer := func() {
		log.Println("Apagando servidor...")

		//Pido el contexto con un DeadLine de 5s
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		//Cuando termine "cierre" el contexto
		defer cancel()

		//Y apago el servidor
		if err := s.Shutdown(ctx); err != nil {
			cobra.CheckErr(fmt.Errorf("no se pudo al apagar el servidor correctamente: %s", err.Error()))
		}
		log.Println("Servidor apagado con exito")
		os.Exit(0)
	}

	if cmd.Flags().Changed("time-live") {
		/*
			Y en una goroutine espero ese tiempo, con time.Sleep(),
			y cierro el programa
		*/
		go func() {
			time.Sleep(DurationTimeOpened)
			apagarServer()
		}()
	}

	//Inicio el servidor en segundo plano
	go func() {
		log.Println("Iniciando servidor...")
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				cobra.CheckErr(fmt.Errorf("no se pudo ecender el servidor: %s", err.Error()))
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
	apagarServer()
}

//Cargo lo que obtuve de las flags en las variabels
func cargarVariables(cmd *cobra.Command, args []string) {

	passed := func(s string) bool {
		return cmd.Flags().Changed(s)
	}

	//Si el usuario ingresa algun directorio(osea que DirToServe no esta vacia) compruebo
	//que el directorio ingresa es una ruta absoluta, sino cierro el programa.
	if len(args) != 0 {
		DirToServe = args[0]
		err := myfuncs.EsAbsolutaYExite(&DirToServe)
		cobra.CheckErr(err)
	} else {

		//En caso de que no haya ingresado busco una ruta local
		localPaths, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(fmt.Errorf("Error al buscar el directorio actual: " + err.Error()))
		}

		DirToServe = filepath.Clean(localPaths)

		//Busco la info del directorio
		i, err := os.Stat(DirToServe)

		//Si la ruta no exite
		if os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("la ruta a servir(\"%s\") no exite: %s", DirToServe, err.Error()))
		}

		if err != nil {
			cobra.CheckErr(fmt.Errorf("ocurrio un error con la ruta a servir (\"%s\"): %s", DirToServe, err.Error()))
		}

		//Si no es un directorio
		if !i.IsDir() {
			cobra.CheckErr(fmt.Errorf("a ruta a servir(\"%s\") debe ser un directorio", DirToServe))
		}
	}

	//Checkeo que la ruta de los templates (haya cambiado o no)
	if passed("template-directory") {
		//Si se ingreso compruebo que se haya ingresado y sea un directorio
		info, err := os.Stat(TemplateDirSeleceted)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("no se pudo accerder al fichero de templates indicado \"%s\": %s", TemplateDirSeleceted, err.Error()))
		}

		if !info.IsDir() {
			cobra.CheckErr(fmt.Errorf("el fichero de templates debe ser un directorio"))
		}
	} else {
		TemplateDirSeleceted = crearConfigDir()
	}

	//Creo el Administrador
	AdminAuth = gin.Accounts{
		"admin": "admin",
	}

	//Si se pasan
	if passed("user") {
		u, err := cmd.Flags().GetStringSlice("user")
		cobra.CheckErr(err)

		if len(u) != 2 {
			cobra.CheckErr(fmt.Errorf("debe ingresar el usuario de administrador en el formato \"-U user,password\""))
		}

		AdminAuth = gin.Accounts{
			u[0]: u[1],
		}
	}

	//Si timeOpen es diferente de 0, es decir, que se ingreso algun valor
	if passed("time-live") {

		//Comrpuebo que el valor no sea un valor negativo
		if DurationTimeOpened <= 0 {
			cobra.CheckErr(fmt.Errorf("se debe ingresar un tiempo de cierre valido (un valor positivo)"))
		}
	}

	//Una vez checkeadas todas las flags
	//Muestro sus estados
	log.Println("Flags:")
	log.Println("- Servidor abierto en:", PortSelected)
	log.Println("- Ruta servida:", DirToServe)
	log.Println("- Tiempo de cierre automático:", DurationTimeOpened.String())
	log.Printf("- Se esta usando el directorio \"%s\" para templates\n", TemplateDirSeleceted)
	log.Printf("- Usuario de administrador: %v", AdminAuth)
}
