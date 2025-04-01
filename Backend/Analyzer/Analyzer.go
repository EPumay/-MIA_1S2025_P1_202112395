package Analyzer

import (
	"flag"
	"fmt"
	"os"
	"proyecto1/DiskManagement"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input

	/*Después de procesar la entrada:
	command será "mkdisk".
	params será "-size=3000 -unit=K -fit=BF -path=/home/bang/Disks/disk1.bin".*/
}

func Analyze(input string) string {
	command, params := getCommandAndParams(input)

	fmt.Println("Comando: ", command, " - ", "Parametro: ", params)

	respuesta := AnalyzeCommnad(command, params)

	return respuesta

}

func AnalyzeCommnad(command string, params string) string {
	var respuesta string
	if strings.Contains(command, "mkdisk") {
		fmt.Print("Comando: mkdisk\n")
		respuesta = fn_mkdisk(params)
	}

	return respuesta
}

func fn_mkdisk(params string) string {
	// Definir flag
	var respuesta string
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño") //nombre, valor por defecto, descripcion
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:]) //parsea los argumentos de la línea de comandos

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1) //encuentra todas las coincidencias de la expresión regular en el input

	// Process the input
	for _, match := range matches {
		flagName := strings.ToLower(match[1]) //guarda el nombre de la flag
		flagValue := match[2]                 //guarda el valor de la flag

		flagValue = strings.Trim(flagValue, "\"") //elimina las comillas del valor de la flag

		switch flagName {
		case "size", "fit", "unit", "path": //compara el nombre de la flag
			fs.Set(flagName, flagValue) //almacena el valor de la flag
		default:
			fmt.Println("Error: Parametro desconocido")
		}
	}

	/*
			Primera Iteración :
		    flagName es "size".
		    flagValue es "3000".
		    El switch encuentra que "size" es un flag reconocido, por lo que se ejecuta fs.Set("size", "3000").
		    Esto asigna el valor 3000 al flag size.

	*/

	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		respuesta = "Error: Size must be greater than 0"
		return respuesta
	}

	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		respuesta = "Error: Fit must be 'bf', 'ff', or 'wf'"
		return respuesta
	}

	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be 'k' or 'm'")
		respuesta = "Error: Unit must be 'k' or 'm'"
		return respuesta
	}

	if *path == "" {
		fmt.Println("Error: Path is required")
		respuesta = "Error: Path is required"
		return respuesta
	}

	respuesta = DiskManagement.Mkdisk(*size, *fit, *unit, *path)
	return respuesta
}
