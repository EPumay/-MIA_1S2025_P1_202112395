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
	Status      [1]byte  //part_status 	Activa/Inactiva
	Type        [1]byte  //part_type 	Primaria/Extendida
	Fit         [1]byte  //part_fit 	Best/FIst/Wors
	Start       int32    //part_start
	Size        int32    //part_s
	Name        [16]byte //part_name
	Correlative int32    //part_correlative
	Id          [4]byte  //part_id
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}
