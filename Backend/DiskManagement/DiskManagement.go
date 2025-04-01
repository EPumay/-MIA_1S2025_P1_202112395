package DiskManagement

import (
	"fmt"
	"math/rand"
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
