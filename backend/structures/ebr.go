package structures

import "fmt"

type EBR struct {
	Part_mount [1]byte 	// Bandera de montaje
	Part_fit [1]byte 	// Tipo de ajuste
	Part_start int32 	// Byte de inicio de la partición
	Part_s int32 		// Tamaño de la partición
	Part_next int32 	// Siguiente EBR
	Part_name [16]byte 	// Nombre de la partición
}

// SerializeEBR escribe la estructura EBR al inicio de una partición extendida
func (ebr *EBR) CreateEBR(ebrStart, ebrSize int, ebrFit, ebrName string) {
	// Asignar el estado de montaje del EBR
	ebr.Part_mount[0] = 'N' // El valor 'N' indica que el EBR no está montado
	
	// Asignar el byte de inicio del EBR
	ebr.Part_start = int32(ebrStart)

	// Asignar el tamaño del EBR
	ebr.Part_s = int32(ebrSize)

	// Asignar el siguiente EBR
	ebr.Part_next = -1 // Inicialmente, no hay siguiente EBR

	// Asignar el tipo de ajuste del EBR
	if len(ebrFit) > 0 {
		ebr.Part_fit[0] = ebrFit[0]
	}

	// Asignar el nombre del EBR
	copy(ebr.Part_name[:], ebrName)
}

func (ebr *EBR) PrintEBR() {
	fmt.Printf("EBR:\n")
	fmt.Printf("  Part_mount: %s\n", string(ebr.Part_mount[:]))
	fmt.Printf("  Part_fit: %s\n", string(ebr.Part_fit[:]))
	fmt.Printf("  Part_start: %d\n", ebr.Part_start)
	fmt.Printf("  Part_s: %d\n", ebr.Part_s)
	fmt.Printf("  Part_next: %d\n", ebr.Part_next)
	fmt.Printf("  Part_name: %s\n", string(ebr.Part_name[:]))
}