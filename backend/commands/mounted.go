package commands

import (
	"backend/stores"
	"fmt"
	"strings"
)

// ParseMount parsea el comando mount y devuelve una instancia de MOUNT
func ParseMounted(tokens []string) (string, error) {
	// Verifica que no se hayan pasado parámetros
	if len(tokens) > 0 {
		return "", fmt.Errorf("el comando 'mounted' no acepta parámetros: %s", strings.Join(tokens, " "))
	}

	// Ejecuta el comando mounted
	return commandMounted()
}

// Ejecuta el comando mounted
func commandMounted() (string, error) {
    // Verifica si hay particiones montadas
    if len(stores.MountedPartitions) == 0 {
        return "No hay particiones montadas.", nil
    }

    // Crea un string para mostrar las particiones montadas
    var result strings.Builder
    result.WriteString("Particiones montadas:\n")
    for id := range stores.MountedPartitions {
        // Obtener la partición montada y su ruta
        partition, path, err := stores.GetMountedPartition(id)
        if err != nil {
            result.WriteString(fmt.Sprintf("ID: %s, Error: %s\n", id, err.Error()))
            continue
        }

        // Agregar la información de la partición al resultado
        result.WriteString(fmt.Sprintf("ID: %s, Path: %s, Nombre: %s\n", id, path, partition.Part_name))
    }

    return result.String(), nil
}