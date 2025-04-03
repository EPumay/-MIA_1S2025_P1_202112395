package Structs

import (
	"fmt"
)

type MBR struct {
	MbrSize    int32        //tamaño del disco
	FechaC     [16]byte     //fecha de creacion
	Id         int32        //mbr_dsk_signature (random de forma unica)
	Fit        [1]byte      // B, F, W
	Partitions [4]Partition //mbr_partitions
}

func PrintMBR(data MBR) {
	fmt.Println("\n\t\tDisco")
	fmt.Printf("CreationDate: %s, fit: %s, size: %d, id: %d\n", string(data.FechaC[:]), string(data.Fit[:]), data.MbrSize, data.Id)
	for i := 0; i < 4; i++ {
		fmt.Printf("Partition %d: %s, %s, %d, %d, %s, %d\n", i, string(data.Partitions[i].Name[:]), string(data.Partitions[i].Type[:]), data.Partitions[i].Start, data.Partitions[i].Size, string(data.Partitions[i].Fit[:]), data.Partitions[i].Correlative)
	}
}

func GetIdMBR(m MBR) int32 {
	return m.Id
}

type Partition struct {
	Status      [1]byte //
	Type        [1]byte // P o E
	Fit         [1]byte // B, F o W
	Start       int32   // byte donde inicia la partición
	Size        int32   //
	Name        [16]byte
	Correlative int32 //desde -1
	Id          [4]byte
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}

type EBR struct {
	Status [1]byte //part_mount (si esta montada)
	Type   [1]byte
	Fit    byte     //part_fit
	Start  int32    //part_start
	Size   int32    //part_s
	Name   [16]byte //part_name
	Next   int32    //part_next
}

func PrintEBR(data EBR) {
	fmt.Println(fmt.Sprintf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c",
		string(data.Name[:]),
		data.Fit,
		data.Start,
		data.Size,
		data.Next,
		data.Status))
}
