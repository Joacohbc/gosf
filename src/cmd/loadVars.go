package cmd

import (
	"fmt"
	"gosf/src/myfuncs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
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

	//Admin accounts
	AdminAuth gin.Accounts
)

//Crear el directorio de configuracion (retornar la ruta del directorio creado)
func crearConfigDir() string {
	//Si no se ingreso nada obtengo el ruta predeterminada
	configPath, err := os.UserConfigDir()
	if err != nil {
		cobra.CheckErr(fmt.Errorf("no se pudo obtener la ruta de configuracion: %s", err.Error()))
	}

	//La Ruta seria -> /home/user/.config/gosf/static/
	templateDir := filepath.Join(configPath, "gosf", "static")

	//Si el directorio no existe, lo creo
	if _, err := os.Stat(templateDir); err != nil {
		err = os.MkdirAll(templateDir, 0755)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("no se pudo ruta de configuracion: %s", err.Error()))
		}
	}

	return templateDir
}

//Cargo lo que obtuve de las flags en las variabels
func cargarVariables(cmd *cobra.Command, args []string) {

	passed := func(s string) bool {
		return cmd.Flags().Changed(s)
	}

	//Si el usuario ingresa algun directorio(osea que DirToServe no esta vacia) compruebo
	//que el directorio ingresa es una ruta absoluta, sino cierro el programa.
	if passed("directory") {
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

	//Una vez checkeadas todas las flags
	//Muestro sus estados
	log.Println("Flags:")
	log.Println("- Servidor abierto en:", PortSelected)
	log.Println("- Ruta servida:", DirToServe)
	log.Println("- Tiempo de cierre automatico:", DurationTimeOpened.String())
	log.Printf("- Se esta usando el directorio \"%s\" para templates\n", TemplateDirSeleceted)
	log.Printf("- Usuario de adminstrador: %v", AdminAuth)
}
