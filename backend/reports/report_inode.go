package reports

import (
	structures "backend/structures"
	utils "backend/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func ReportInode(superblock *structures.SuperBlock, diskPath string, path string) error {
    // Crear las carpetas padre si no existen
    err := utils.CreateParentDirs(path)
    if err != nil {
        return err
    }

    // Obtener el nombre base del archivo sin la extensión
    dotFileName, outputImage := utils.GetFileNames(path)

    // Iniciar el contenido DOT con estilo mejorado
    dotContent := `digraph G {
        rankdir=LR;
        node [shape=plaintext];
        graph [bgcolor="#f5f5f5", fontname="Arial", fontsize=12];
        edge [color="#666666", arrowsize=0.8];
        
    `

    // Paleta de colores profesional
    colors := struct {
        header      string
        evenRow     string
        oddRow      string
        accent1     string
        accent2     string
        accent3     string
        accent4     string
        accent5     string
        accent6     string
    }{
        header:      "#2c3e50",
        evenRow:     "#ecf0f1",
        oddRow:      "#ffffff",
        accent1:     "#3498db",
        accent2:     "#e74c3c",
        accent3:     "#2ecc71",
        accent4:     "#f39c12",
        accent5:     "#9b59b6",
        accent6:     "#1abc9c",
    }

    // Iterar sobre cada inodo
    for i := int32(0); i < superblock.S_inodes_count; i++ {
        inode := &structures.Inode{}
        // Deserializar el inodo
        err := inode.Deserialize(diskPath, int64(superblock.S_inode_start+(i*superblock.S_inode_size)))
        if err != nil {
            return err
        }

        // Convertir tiempos a string
        atime := time.Unix(int64(inode.I_atime), 0).Format("2006-01-02 15:04:05")
        ctime := time.Unix(int64(inode.I_ctime), 0).Format("2006-01-02 15:04:05")
        mtime := time.Unix(int64(inode.I_mtime), 0).Format("2006-01-02 15:04:05")

        // Definir el contenido DOT para el inodo actual con diseño mejorado
        dotContent += fmt.Sprintf(`inode%d [label=<
            <table border="0" cellborder="1" cellspacing="0" cellpadding="8" style="rounded" bgcolor="%s" gradientangle="315">
                <!-- Encabezado -->
                <tr>
                    <td colspan="2" bgcolor="%s" style="rounded" border="0">
                        <font color="white" face="Arial" point-size="14"><b>INODO %d</b></font>
                    </td>
                </tr>
                
                <!-- Datos del inodo -->
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>UID</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>GID</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Tamaño</b></font></td>
                    <td bgcolor="%s" border="0">%d bytes</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Último Acceso</b></font></td>
                    <td bgcolor="%s" border="0">%s</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Creación</b></font></td>
                    <td bgcolor="%s" border="0">%s</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Modificación</b></font></td>
                    <td bgcolor="%s" border="0">%s</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Tipo</b></font></td>
                    <td bgcolor="%s" border="0">%c</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font color="white"><b>Permisos</b></font></td>
                    <td bgcolor="%s" border="0">%s</td>
                </tr>
                
                <!-- Encabezado de bloques -->
                <tr>
                    <td colspan="2" bgcolor="%s" style="rounded" border="0">
                        <font color="white" face="Arial"><b>BLOQUES DIRECTOS</b></font>
                    </td>
                </tr>
            `, 
            i, colors.evenRow, 
            colors.header, i,
            colors.accent1, colors.oddRow, inode.I_uid,
            colors.accent2, colors.oddRow, inode.I_gid,
            colors.accent3, colors.oddRow, inode.I_size,
            colors.accent4, colors.oddRow, atime,
            colors.accent5, colors.oddRow, ctime,
            colors.accent6, colors.oddRow, mtime,
            colors.accent1, colors.oddRow, rune(inode.I_type[0]),
            colors.accent2, colors.oddRow, string(inode.I_perm[:]),
            colors.header)

        // Agregar los bloques directos con estilo alternado
        for j, block := range inode.I_block {
            if j > 11 {
                break
            }
            rowColor := colors.evenRow
            if j%2 == 0 {
                rowColor = colors.oddRow
            }
            dotContent += fmt.Sprintf(`<tr>
                <td bgcolor="%s" border="0"><font><b>Bloque %d</b></font></td>
                <td bgcolor="%s" border="0">%d</td>
            </tr>`, colors.accent3, j+1, rowColor, block)
        }

        // Agregar los bloques indirectos con estilo mejorado
        dotContent += fmt.Sprintf(`
                <!-- Encabezado de bloques indirectos -->
                <tr>
                    <td colspan="2" bgcolor="%s" style="rounded" border="0">
                        <font color="white" face="Arial"><b>BLOQUES INDIRECTOS</b></font>
                    </td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Indirecto</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Indirecto Doble</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>
                <tr>
                    <td bgcolor="%s" border="0"><font><b>Indirecto Triple</b></font></td>
                    <td bgcolor="%s" border="0">%d</td>
                </tr>
            </table>>];
        `, 
        colors.header,
        colors.accent4, colors.oddRow, inode.I_block[12],
        colors.accent5, colors.evenRow, inode.I_block[13],
        colors.accent6, colors.oddRow, inode.I_block[14])

        // Agregar enlace al siguiente inodo con estilo mejorado
        if i < superblock.S_inodes_count-1 {
            dotContent += fmt.Sprintf(`
                inode%d -> inode%d [
                    color="%s",
                    penwidth=1.5,
                    arrowhead=open,
                    arrowsize=0.8,
                    style="dashed"
                ];
            `, i, i+1, colors.accent1)
        }
    }

    // Cerrar el contenido DOT
    dotContent += `
        /* Estilos globales para los nodos */
        node [fontname="Arial", fontsize=10, shape=box, style="rounded,filled", 
              fillcolor="#ffffff", color="#2c3e50", penwidth=1.5];
    }`

    // Crear el archivo DOT
    dotFile, err := os.Create(dotFileName)
    if err != nil {
        return err
    }
    defer dotFile.Close()

    // Escribir el contenido DOT en el archivo
    _, err = dotFile.WriteString(dotContent)
    if err != nil {
        return err
    }

    // Generar la imagen con Graphviz (agregando opciones para mejor calidad)
    cmd := exec.Command("dot", "-Tpng", "-Gdpi=300", "-Nfontname=Arial", "-Efontname=Arial", 
        dotFileName, "-o", outputImage)
    err = cmd.Run()
    if err != nil {
        return err
    }

    fmt.Printf("\x1b[32mReporte de inodos generado exitosamente:\x1b[0m %s\n", outputImage)
    return nil
}