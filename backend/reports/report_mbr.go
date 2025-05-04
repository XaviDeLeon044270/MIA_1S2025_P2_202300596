package reports

import (
	structures "backend/structures"
	utils "backend/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ReportMBR(mbr *structures.MBR, path string) error {
    // Crear las carpetas padre si no existen
    err := utils.CreateParentDirs(path)
    if err != nil {
        return fmt.Errorf("error al crear directorios padres: %v", err)
    }

    // Obtener el nombre base del archivo sin la extensión
    dotFileName, outputImage := utils.GetFileNames(path)

    // Definir paleta de colores
    colors := struct {
        primary    string
        secondary  string
        accent1    string
        accent2    string
        accent3    string
        evenRow    string
        oddRow     string
        partition  string
    }{
        primary:   "#2c3e50",
        secondary: "#34495e",
        accent1:   "#3498db",
        accent2:   "#e74c3c",
        accent3:   "#2ecc71",
        evenRow:   "#ecf0f1",
        oddRow:    "#ffffff",
        partition: "#9b59b6",
    }

    // Crear el archivo DOT primero para manejar errores temprano
    file, err := os.Create(dotFileName)
    if err != nil {
        return fmt.Errorf("error al crear archivo DOT: %v", err)
    }
    defer file.Close()

    // Escribir encabezado del gráfico
    _, err = file.WriteString(fmt.Sprintf(`digraph G {
        node [shape=plaintext];
        graph [bgcolor="#f5f5f5", fontname="Arial", fontsize=12];
        edge [color="#666666", arrowsize=0.8];
        
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0" cellpadding="8" style="rounded" bgcolor="%s">
                <tr>
                    <td colspan="2" bgcolor="%s" style="rounded" border="0">
                        <font color="white" face="Arial" point-size="16"><b>REPORTE MBR</b></font>
                    </td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Tamaño del MBR</b></font></td>
                    <td bgcolor="%s" border="0">%d bytes</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Fecha de Creación</b></font></td>
                    <td bgcolor="%s" border="0">%s</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Firma del Disco</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>`,
        colors.evenRow,
        colors.primary,
        colors.accent1, colors.oddRow, mbr.Mbr_size,
        colors.accent2, colors.oddRow, time.Unix(int64(mbr.Mbr_creation_date), 0).Format("2006-01-02 15:04:05"),
        colors.accent3, colors.oddRow, mbr.Mbr_disk_signature))

    if err != nil {
        return fmt.Errorf("error al escribir encabezado: %v", err)
    }

    // Procesar particiones
    for i, part := range mbr.Mbr_partitions {
        if part.Part_size == -1 {
            continue
        }

        partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
        partStatus := rune(part.Part_status[0])
        partType := rune(part.Part_type[0])
        partFit := rune(part.Part_fit[0])

        rowColor := colors.evenRow
        if i%2 == 0 {
            rowColor = colors.oddRow
        }

        _, err = file.WriteString(fmt.Sprintf(`
                <tr>
                    <td colspan="2" bgcolor="%s" style="rounded" border="0">
                        <font color="white" face="Arial" point-size="14"><b>PARTICIÓN %d</b></font>
                    </td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Estado</b></font></td>
                    <td bgcolor="%s" border="0">%c</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Tipo</b></font></td>
                    <td bgcolor="%s" border="0">%c</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Ajuste</b></font></td>
                    <td bgcolor="%s" border="0">%c</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Inicio</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Tamaño</b></font></td>
                    <td bgcolor="%s" border="0">%d bytes</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Nombre</b></font></td>
                    <td bgcolor="%s" border="0">%s</td>
                </tr>`,
            colors.partition,
            i+1,
            colors.accent1, rowColor, partStatus,
            colors.accent2, rowColor, partType,
            colors.accent3, rowColor, partFit,
            colors.accent1, rowColor, part.Part_start,
            colors.accent2, rowColor, part.Part_size,
            colors.accent3, rowColor, partName))

        if err != nil {
            return fmt.Errorf("error al escribir datos de partición: %v", err)
        }
    }

    // Cerrar la tabla y el gráfico
    _, err = file.WriteString(`</table>>];
        
        node [fontname="Arial", fontsize=10, shape=box, style="rounded,filled", 
              fillcolor="#ffffff", color="#2c3e50", penwidth=1.5];
    }`)
    if err != nil {
        return fmt.Errorf("error al cerrar archivo DOT: %v", err)
    }

    // Generar la imagen
    cmd := exec.Command("dot", "-Tpng", "-Gdpi=300", "-Nfontname=Arial", "-Efontname=Arial", 
        dotFileName, "-o", outputImage)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("error al generar imagen: %v", err)
    }

    fmt.Printf("\x1b[32m✓ Reporte MBR generado:\x1b[0m %s\n", outputImage)
    return nil
}