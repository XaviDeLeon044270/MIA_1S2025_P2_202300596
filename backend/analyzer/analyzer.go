package analyzer

import (
    commands "backend/commands"
    "errors"  // Importa el paquete "errors" para manejar errores
    "fmt"     // Importa el paquete "fmt" para formatear e imprimir texto
    "strings" // Importa el paquete "strings" para manipulación de cadenas
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (string, error) {
    // Elimina espacios en blanco al inicio y al final de la entrada
    input = strings.TrimSpace(input)

    // Ignorar líneas que comiencen con el token #
    if strings.HasPrefix(input, "#") {
        return "", nil // Retorna vacío y sin error
    }

    // Divide la entrada en tokens usando espacios en blanco como delimitadores
    tokens := strings.Fields(input)

    // Si no se proporcionó ningún comando, devuelve un error
    if len(tokens) == 0 {
        return "", errors.New("no se proporcionó ningún comando")
    }

    // Convertir el comando principal a minúsculas para hacerlo case insensitive
    command := strings.ToLower(tokens[0])

    // Switch para manejar diferentes comandos
    switch command {
    case "mkdisk":
        return commands.ParseMkdisk(tokens[1:])

    case "rmdisk":
        return commands.ParseRmdisk(tokens[1:])

    case "fdisk":
        return commands.ParseFdisk(tokens[1:])
        
    case "mount":
        return commands.ParseMount(tokens[1:])

    case "mkfs":
        return commands.ParseMkfs(tokens[1:])

    case "rep":
        return commands.ParseRep(tokens[1:])

    case "mkdir":
        return commands.ParseMkdir(tokens[1:])

    case "login":
        return commands.ParseLogin(tokens[1:])

    case "logout":
        return commands.ParseLogout(tokens[1:])

    case "mounted":
        return commands.ParseMounted(tokens[1:])
        
    default:
        // Si el comando no es reconocido, devuelve un error
        return "", fmt.Errorf("comando desconocido: %s", tokens[0])
    }
}