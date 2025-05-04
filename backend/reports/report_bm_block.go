package reports

import (
	structures "backend/structures"
	utils "backend/utils"
	"fmt"
	"os"
	"strings"
)


func ReportBMBlock(superblock *structures.SuperBlock, diskPath string, path string) error {
    // Crear las carpetas padre si no existen
    err := utils.CreateParentDirs(path)
    if err != nil {
        return err
    }

    // Abrir el archivo de disco
    file, err := os.Open(diskPath)
    if err != nil {
        return fmt.Errorf("error al abrir el archivo de disco: %v", err)
    }
    defer file.Close()

    // Obtener la cantidad total de bloques
    totalBlocks := superblock.S_blocks_count

    // Construir el contenido del bitmap
    var bitmapContent strings.Builder

    for i := int32(0); i < totalBlocks; i++ {
        // Mover el puntero al bitmap de bloques
        _, err := file.Seek(int64(superblock.S_bm_block_start+i), 0)
        if err != nil {
            return fmt.Errorf("error al mover el puntero en el archivo: %v", err)
        }

        // Leer el bit (debe ser '0' o '1')
        char := make([]byte, 1)
        _, err = file.Read(char)
        if err != nil {
            return fmt.Errorf("error al leer el byte del archivo: %v", err)
        }

        // Validar el carácter leído
        if char[0] != '0' && char[0] != '1' {
            char[0] = '0' // Corregir a '0' si no es válido
        }

        // Escribir el bit en el contenido
        bitmapContent.WriteString(string(char))

        // Agregar un espacio entre bits para mejor visualización
        bitmapContent.WriteString(" ")

        // Agregar un salto de línea cada 20 bits
        if (i+1)%20 == 0 {
            bitmapContent.WriteString("\n")
        }
    }

    // Crear el archivo de salida
    txtFile, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("error al crear el archivo TXT: %v", err)
    }
    defer txtFile.Close()

    // Escribir el contenido en el archivo
    _, err = txtFile.WriteString(bitmapContent.String())
    if err != nil {
        return fmt.Errorf("error al escribir en el archivo TXT: %v", err)
    }

    fmt.Println("Archivo del bitmap de bloques generado:", path)
    return nil
}