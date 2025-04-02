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
	} else if strings.Contains(command, "rmdisk") {
		fmt.Print("Comando: rmdisk\n")
		respuesta = fn_rmdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fmt.Print("Comando: fdisk\n")
		respuesta = fn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		fmt.Print("Comando: mount\n")
		respuesta = fn_mount(params)
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
			return "\n Error: Parametro desconocido"
		}
	}
	//pasar flags a minisculas menos path
	*fit = strings.ToLower(*fit)
	*unit = strings.ToLower(*unit)

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

func fn_rmdisk(input string) (respuesta string) {
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := match[2]
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}
	if *path == "" {
		fmt.Println("Error: Path is required")
		return
	}
	respuesta = DiskManagement.Rmdisk(*path)
	return respuesta
}

func fn_fdisk(input string) (respuesta string) {
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "m", "Unidad")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "", "Ajuste")

	// Parsear los flags
	fs.Parse(os.Args[1:])

	// Encontrar los flags en el input
	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar el input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Convertir el nombre y la unidad a minúsculas
	*name = strings.ToLower(*name)
	*unit = strings.ToLower(*unit)
	*type_ = strings.ToLower(*type_)
	*fit = strings.ToLower(*fit)

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		respuesta = "Error: Size must be greater than 0"
		return respuesta
	}

	if *path == "" {
		fmt.Println("Error: Path is required")
		return "Error: Path is required"
	}

	// Si no se proporcionó un fit, usar el valor predeterminado "w"
	if *fit == "" {
		*fit = "w"
	}

	// Validar fit (b/w/f)
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error: Fit must be 'b', 'f', or 'w'")
		return "Error: Fit must be 'b', 'f', or 'w'"
	}

	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be 'k' or 'm'")
		return "Error: Unit must be 'k' or 'm'"
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Println("Error: Type must be 'p', 'e', or 'l'")
		return "Error: Type must be 'p', 'e', or 'l'"
	}

	// Llamar a la función
	respuesta = DiskManagement.Fdisk(*size, *path, *name, *unit, *type_, *fit)
	return respuesta
}

func fn_mount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre de la partición")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2]) // Convertir todo a minúsculas
		flagValue = strings.Trim(flagValue, "\"")
		fs.Set(flagName, flagValue)
	}

	if *path == "" || *name == "" {
		fmt.Println("Error: Path y Name son obligatorios")
		return
	}

	lowercaseName := strings.ToLower(*name)
	DiskManagement.Mount(*path, lowercaseName)
}
