package DiskManagement

import (
	"fmt"
	"math/rand"
	"os"
	"proyecto1/Structs"
	"proyecto1/Utilities"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) string {
	var respuesta string

	respuesta = "Comando: mkdisk\n"
	respuesta += "Tamaño: " + fmt.Sprint(size) + "\n"
	respuesta += "Ajuste: " + fit + "\n"
	respuesta += "Unidad: " + unit + "\n"
	respuesta += "Ruta: " + path + "\n"
	respuesta += "-------------------------------------\n"

	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser bf, wf or ff")
		respuesta = "Error: Fit debe ser bf, wf or ff"
		return respuesta
	}
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayo a  0")
		respuesta = "Error: Size debe ser mayo a  0"
		return respuesta
	}
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades validas son k o m")
		respuesta = "Error: Las unidades validas son k o m"
		return respuesta
	}

	/*
		Si el usuario especifica unit = "k" (Kilobytes), el tamaño se multiplica por 1024 para convertirlo a bytes.
		Si el usuario especifica unit = "m" (Megabytes), el tamaño se multiplica por 1024 * 1024 para convertirlo a MEGA bytes.
	*/
	// Asignar tamanio
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Crear el archivo
	err := Utilities.CreateFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		respuesta = "Error: " + err.Error()
		return respuesta
	}

	// Abrir el archivo
	file, err := Utilities.OpenFile(path)
	if err != nil {
		return err.Error()
	}

	//llenar el archivo con ceros
	datos := make([]byte, size)
	newErr := Utilities.WriteObject(file, datos, 0)
	if newErr != nil {
		fmt.Println("MKDISK Error: ", newErr)
		return "MKDISK Error: " + newErr.Error()
	}
	var newMBR Structs.MBR
	newMBR.MbrSize = int32(size)
	newMBR.Id = rand.Int31() // Numero random rand.Int31() genera solo números no negativos
	copy(newMBR.Fit[:], fit)
	ahora := time.Now()
	copy(newMBR.FechaC[:], ahora.Format("02/01/2006 15:04"))
	// Escribir el MBR en el archivo
	if err := Utilities.WriteObject(file, newMBR, 0); err != nil {
		return "ERROR"
	}
	// Cerrar el archivo
	defer file.Close()

	return respuesta
}

func Rmdisk(path string) (respuesta string) {
	fmt.Println("*************Inicio RMDISK*************")
	fmt.Println("Path: ", path)

	// Verificar si el archivo existe primero
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Error: El archivo no existe")
		return "Error: El disco no existe"
	}

	// Eliminar el archivo
	err := os.Remove(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return "Error: " + err.Error()
	}
	fmt.Println("*************Fin RMDISK*************")
	return "El disco ha sido eliminado"

}

func Fdisk(size int, path string, name string, unit string, type_ string, fit string) (respuesta string) {
	fmt.Println("*************Inicio FDISK*************")
	fmt.Println("Tamaño: ", size)
	fmt.Println("Unidad: ", unit)
	fmt.Println("Ruta: ", path)
	fmt.Println("Nombre: ", name)
	fmt.Println("Ajuste: ", fit)

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}
	// Verificar si el archivo existe primero
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Error: El archivo no existe")
		return "Error: El disco no existe"
	}

	// Abrir el archivo
	file, err := Utilities.OpenFile(path)
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	// Leer el MBR
	var mbr Structs.MBR
	err = Utilities.ReadObject(file, &mbr, 0)
	if err != nil {
		return "Error al leer el MBR"

	}

}
