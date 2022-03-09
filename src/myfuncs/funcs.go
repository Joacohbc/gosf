package myfuncs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//Pone el primer caracter de un string en Mayus
func PrimeraMayus(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}

	s = strings.ToUpper(string(s[0])) + s[1:]
	return s
}

//Comprueba si una ruta es absolutas, si no lo es intenta
//encontrarla y guardar en la variable. Tambien valida que fichero exita
//
//Simpre devuelve la ruta "limpia" (usando filepath.Clean)
//sea o no absulta la ruta de entrada
func EsAbsolutaYExite(path *string) error {

	//Si absoluta ya es absoluta no haho nada
	if filepath.IsAbs(*path) {
		*path = filepath.Clean(*path)
		return nil
	}

	//Si no es absoluta inteto encontrar la ruta absoluta
	absPath, err := filepath.Abs(*path)
	if err != nil {
		return fmt.Errorf("fallo al encontrar la ruta absoluta de \"%s\": %s", *path, err.Error())
	}

	//Compruebo si existe
	_, err = os.Stat(absPath)

	//Si no existe el fichero
	if os.IsNotExist(err) {
		return fmt.Errorf("no existe el fichero \"%s\": %s", absPath, err.Error())
	}

	//Si ocurrio un directorio
	if err != nil {
		return fmt.Errorf("no se pudo acceder al fichero \"%s\": %s", absPath, err.Error())
	}

	//Si todo ocurrio con exito escribo el archivo
	*path = filepath.Clean(absPath)
	return nil
}
