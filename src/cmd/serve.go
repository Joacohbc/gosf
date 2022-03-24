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
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

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

	log.Println("Iniciando servidor...")

	//Si timeOpen es diferente de 0, es decir, que se ingreso algun valor
	if cmd.Flags().Changed("time-live") {

		//Comrpuebo que el valor no sea un valor negativo
		if DurationTimeOpened <= 0 {
			cobra.CheckErr(fmt.Errorf("se debe ingresar un tiempo de cierre valido (un valor positivo)"))
		}

		/*
			Y en una goroutine espero ese tiempo, con time.Sleep(),
			y cierro el programa
		*/
		go func() {
			time.Sleep(DurationTimeOpened)

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
			os.Exit(0)
		}()
	}

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
