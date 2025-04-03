package DiskManagement

import (
	"encoding/binary"
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
		fmt.Println("Error: El Disco no existe")
		return "Error: El disco no existe"
	}

	// Abrir el archivo
	file, err := Utilities.OpenFile(path)
	if err != nil {
		return err.Error()
	}

	// Leer el MBR
	var TempMBR Structs.MBR
	err = Utilities.ReadObject(file, &TempMBR, 0)
	if err != nil {
		return "Error al leer el MBR"

	}
	var primaryCount, extendedCount, totalPartitions int
	var usedSpace int32 = 0

	Structs.PrintMBR(TempMBR)
	for i := 0; i < 4; i++ { //4 son las particiones primarias permitidas
		if TempMBR.Partitions[i].Size != 0 { //si la particion esta es uso, size es distinto de 0
			totalPartitions++                       //contador de particiones existentes
			usedSpace += TempMBR.Partitions[i].Size //suma el espacio usado

			if TempMBR.Partitions[i].Type[0] == 'p' {
				primaryCount++ //contador de particiones primarias
			} else if TempMBR.Partitions[i].Type[0] == 'e' {
				extendedCount++ //contador de particiones extendidas
			}
		}
		//no estan las logicas, porque solo pueden existir dentro de una extendida
	}

	// Validar que no se exceda el número máximo de particiones primarias y extendidas
	if totalPartitions >= 4 {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return "Error: No se pueden crear más de 4 particiones primarias o extendidas en total."
	}

	// Validar que no se pueda crear una partición lógica sin una extendida
	if type_ == "l" && extendedCount == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return "Error: No se puede crear una partición lógica sin una partición extendida."
	}

	// Validar que el tamaño de la nueva partición no exceda el tamaño del disco
	if usedSpace+int32(size) > TempMBR.MbrSize {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return "Error: No hay suficiente espacio en el disco para crear esta partición."
	}
	// Determinar la posición de inicio de la nueva partición
	var gap int32 = int32(binary.Size(TempMBR))
	if totalPartitions > 0 {
		gap = TempMBR.Partitions[totalPartitions-1].Start + TempMBR.Partitions[totalPartitions-1].Size
	}

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size == 0 {
			if type_ == "p" || type_ == "e" {
				// Crear partición primaria o extendida
				TempMBR.Partitions[i].Size = int32(size)
				TempMBR.Partitions[i].Start = gap
				copy(TempMBR.Partitions[i].Name[:], name)
				copy(TempMBR.Partitions[i].Fit[:], fit)
				copy(TempMBR.Partitions[i].Status[:], "0")
				copy(TempMBR.Partitions[i].Type[:], type_)
				TempMBR.Partitions[i].Correlative = int32(totalPartitions + 1)

				if type_ == "e" {
					// Inicializar el primer EBR en la partición extendida
					ebrStart := gap // El primer EBR se coloca al inicio de la partición extendida
					ebr := Structs.EBR{
						Fit:   fit[0],
						Start: ebrStart,
						Size:  0,
						Next:  -1,
					}
					copy(ebr.Name[:], "")
					Utilities.WriteObject(file, ebr, int64(ebrStart))
				}

				break
			}
		}
	}
	// Manejar la creación de particiones lógicas dentro de una partición extendida
	if type_ == "l" {
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Type[0] == 'e' {
				ebrPos := TempMBR.Partitions[i].Start
				var ebr Structs.EBR
				for {
					Utilities.ReadObject(file, &ebr, int64(ebrPos))
					if ebr.Next == -1 {
						break
					}
					ebrPos = ebr.Next
				}

				// Calcular la posición de inicio de la nueva partición lógica
				newEBRPos := ebr.Start + ebr.Size                            // El nuevo EBR se coloca después de la partición lógica anterior
				logicalPartitionStart := newEBRPos + int32(binary.Size(ebr)) // El inicio de la partición lógica es justo después del EBR

				// Ajustar el siguiente EBR
				ebr.Next = newEBRPos
				Utilities.WriteObject(file, ebr, int64(ebrPos))

				// Crear y escribir el nuevo EBR
				newEBR := Structs.EBR{
					Fit:   fit[0],
					Start: logicalPartitionStart,
					Size:  int32(size),
					Next:  -1,
				}
				copy(newEBR.Name[:], name)
				Utilities.WriteObject(file, newEBR, int64(newEBRPos))

				// Imprimir el nuevo EBR creado
				fmt.Println("Nuevo EBR creado:")
				Structs.PrintEBR(newEBR)
				fmt.Println("")

				// Imprimir todos los EBRs en la partición extendida
				fmt.Println("Imprimiendo todos los EBRs en la partición extendida:")
				ebrPos = TempMBR.Partitions[i].Start
				for {
					err := Utilities.ReadObject(file, &ebr, int64(ebrPos))
					if err != nil {
						fmt.Println("Error al leer EBR:", err)
						break
					}
					Structs.PrintEBR(ebr)
					if ebr.Next == -1 {
						break
					}
					ebrPos = ebr.Next
				}

				break
			}
		}
		fmt.Println("")
	}

	// Sobrescribir el MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return "Error: Could not write MBR to file"
	}

	var TempMBR2 Structs.MBR
	// Leer el objeto nuevamente para verificar
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file after writing")
		return "Error: Could not read MBR from file after writing"
	}

	// Imprimir el objeto MBR actualizado
	Structs.PrintMBR(TempMBR2)

	// Cerrar el archivo binario
	defer file.Close()

	fmt.Println("======FIN FDISK======")
	fmt.Println("")

	respuesta = "Comando: FDISK\n"
	respuesta += "Tamaño: " + fmt.Sprint(size) + "\n"
	respuesta += "Unidad: " + unit + "\n"
	respuesta += "Ruta: " + path + "\n"
	respuesta += "Nombre: " + name + "\n"
	respuesta += "Ajuste: " + fit + "\n"
	respuesta += "-------------------------------------\n"
	respuesta += "Partición creada con éxito\n"
	respuesta += "-------------------------------------\n"
	respuesta += "Fin FDISK\n"
	return respuesta
}
