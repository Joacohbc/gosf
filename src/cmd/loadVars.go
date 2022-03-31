package cmd

import (
	"fmt"
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
