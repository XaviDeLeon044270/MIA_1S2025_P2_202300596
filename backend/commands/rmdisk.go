package commands

import (
	"errors"        // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"           // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"os"            // Paquete para interactuar con el sistema operativo
	"regexp"        // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas
	"strings"       // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
	"backend/utils"
)

// RMDISK estructura que representa el comando rmdisk con sus parámetros
type RMDISK struct {
	path string // Ruta del archivo del disco
}
//
// ParseRmdisk parsea el comando rmdisk y devuelve una instancia de RMDISK
func ParseRmdisk(tokens []string) (string, error) {
	cmd := &RMDISK{} // Crea una nueva instancia de RMDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando rmdisk
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-path":
			cmd.path = value // Asigna el valor a la ruta del disco

			// Verifica si la ruta es válida
			if !strings.HasSuffix(cmd.path, ".mia") {
				return "", errors.New("la ruta debe terminar con .mia")
			}

			// Verifica si el disco existe
			if !utils.FileExists(cmd.path) {
				return "", errors.New("el disco no existe")
			}
			
			break
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.path == "" {
		return "", errors.New("se requiere el parámetro -path")
	}

	// Ejecuta el comando rmdisk
	err := cmd.Execute()
	if err != nil {
		return "", fmt.Errorf("error al ejecutar el comando rmdisk: %v", err)
	}

	return fmt.Sprintf("Disco eliminado correctamente"), nil

}

// Execute ejecuta el comando rmdisk
func (cmd *RMDISK) Execute() error {
	// Elimina el disco
	err := os.Remove(cmd.path)
	if err != nil {
		return fmt.Errorf("error al eliminar el disco: %v", err)
	}

	fmt.Printf("Disco eliminado: %s\n", cmd.path)

	// Actualiza la lista de discos montados
	/*for i, disk := range structures.MountedDisks {
		if disk.Path == cmd.path {
			structures.MountedDisks = append(structures.MountedDisks[:i], structures.MountedDisks[i+1:]...)
			break
		}
	}*/

	return nil
}