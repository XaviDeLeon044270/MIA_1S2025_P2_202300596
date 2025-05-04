package structures

import (
	"encoding/binary"
	"os"
)

type PointerBlock struct {
	P_pointers [16]int32 // 16 * 4 = 64 bytes
	// Total: 64 bytes
}

func (pb *PointerBlock) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición deseada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Leer el bloque desde el archivo
	err = binary.Read(file, binary.LittleEndian, pb)
	if err != nil {
		return err
	}

	return nil
}