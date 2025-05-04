package commands

import (
    "backend/stores"
    "errors"
    "fmt"
    "strings"
)

// ParseLogout parsea el comando logout y devuelve una instancia de LOGOUT
func ParseLogout(tokens []string) (string, error) {
    // Verificar que no se hayan pasado parámetros
    if len(tokens) > 0 {
        return "", fmt.Errorf("el comando 'logout' no acepta parámetros: %s", strings.Join(tokens, " "))
    }

    // Ejecutar el comando logout
    err := commandLogout()
    if err != nil {
        return "", err
    }

    return "Sesión cerrada con éxito.", nil
}

// Execute ejecuta el comando logout
func commandLogout() error {
    // Verifica si hay usuario logueado
    if stores.Auth.IsAuthenticated() {
        // Cierra la sesión del usuario
        stores.Auth.Logout()
        fmt.Println("Sesión cerrada con éxito.")
    } else {
        return errors.New("no hay usuario logueado")
    }

    return nil
}