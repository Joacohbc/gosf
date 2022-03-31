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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var setConfigCmd = &cobra.Command{
	Use:     "set-template-dir",
	Aliases: []string{"stc"},
	Short:   "Mueve (no copia), el directorio indicado a la carpeta de preterminada de templates",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//Obtengo el Directorio predeterminado
		config := crearConfigDir()

		//Obtengo lo que paso el usuario de la config
		dir, err := cmd.Flags().GetString("set-config")
		cobra.CheckErr(err)

		//Comprubo que exista el fichero
		info, err := os.Stat(dir)
		cobra.CheckErr(err)

		//Comprubo que sea un directorio
		if !info.IsDir() {
			cobra.CheckErr(fmt.Errorf("el fichero debe ser un directorio"))
		}

		//Vacio el directorio predeterminado
		err = os.RemoveAll(config)
		cobra.CheckErr(err)

		//Y muevo los archivos (rename es lo mismo que move mv)
		err = os.Rename(dir, config)
		cobra.CheckErr(err)
	},
}

//Agrego las flags
func init() {
	rootCmd.AddCommand(setConfigCmd)
}
