package reports

import (
	"backend/structures"
	"backend/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// GenerateDiskReport crea un reporte gráfico de la estructura del disco con estilo profesional
func ReportDisk(mbr *structures.MBR, diskPath, outputPath string) error {
	fmt.Println("Generando reporte del disco...")

	// Crear directorios si no existen
	if err := utils.CreateParentDirs(outputPath); err != nil {
		return fmt.Errorf("error al crear directorios: %v", err)
	}

	// Paleta de colores profesional
	colors := struct {
		Primary     string
		Secondary   string
		PrimaryPart string
		ExtendedPart string
		FreeSpace   string
		Header      string
		EvenRow     string
		OddRow      string
		Text        string
	}{
		Primary:     "#3498db",
		Secondary:   "#2c3e50",
		PrimaryPart: "#3498db",
		ExtendedPart: "#e74c3c",
		FreeSpace:   "#2ecc71",
		Header:      "#2c3e50",
		EvenRow:     "#ecf0f1",
		OddRow:      "#ffffff",
		Text:        "#333333",
	}

	totalDiskSize := mbr.Mbr_size
	usedSpace := 0

	// Obtener nombres de archivo para DOT e imagen de salida
	dotFile, imgFile := utils.GetFileNames(outputPath)

	// Iniciar contenido DOT con estilo mejorado
	dotData := fmt.Sprintf(`digraph DiskReport {
		rankdir=LR;
		bgcolor="#f5f5f5";
		node [shape=plaintext, fontname="Arial"];
		edge [color="#666666", arrowsize=0.8];
		
		disk [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="8" style="rounded" bgcolor="%s">
				<!-- Encabezado -->
				<tr>
					<td colspan="3" bgcolor="%s" style="rounded" border="0">
						<font color="white" face="Arial" point-size="16"><b>ESTRUCTURA DEL DISCO</b></font>
					</td>
				</tr>
				<tr>
					<td bgcolor="%s" border="0"><font color="white"><b>Tipo</b></font></td>
					<td bgcolor="%s" border="0"><font color="white"><b>Tamaño</b></font></td>
					<td bgcolor="%s" border="0"><font color="white"><b>Porcentaje</b></font></td>
				</tr>
	`, colors.EvenRow, colors.Header, colors.Secondary, colors.Secondary, colors.Secondary)

	// Información general del disco
	creationDate := time.Unix(int64(mbr.Mbr_creation_date), 0).Format("2006-01-02 15:04:05")
	dotData += fmt.Sprintf(`
		<tr>
			<td bgcolor="%s" border="0"><font><b>Tamaño Total</b></font></td>
			<td bgcolor="%s" border="0">%d bytes</td>
			<td bgcolor="%s" border="0">100.00%%</td>
		</tr>
		<tr>
			<td bgcolor="%s" border="0"><font><b>Fecha Creación</b></font></td>
			<td colspan="2" bgcolor="%s" border="0">%s</td>
		</tr>
	`, colors.OddRow, colors.OddRow, totalDiskSize, colors.OddRow, colors.EvenRow, colors.EvenRow, creationDate)

	// Recorrer las particiones del MBR
	for index, partition := range mbr.Mbr_partitions {
		partType := rune(partition.Part_type[0])
		partSize := partition.Part_size

		if partType == 'P' || partType == 'E' {
			usedSpace += int(partSize)
			partLabel := map[rune]string{'P': "Primaria", 'E': "Extendida"}[partType]
			partColor := colors.PrimaryPart
			if partType == 'E' {
				partColor = colors.ExtendedPart
			}

			// Agregar partición al reporte DOT
			dotData += fmt.Sprintf(`
				<!-- Partición %d -->
				<tr>
					<td colspan="3" bgcolor="%s" style="rounded" border="0">
						<font color="white" face="Arial" point-size="14"><b>Partición %d (%s)</b></font>
					</td>
				</tr>
				<tr>
					<td bgcolor="%s" border="0">%s</td>
					<td bgcolor="%s" border="0">%d bytes</td>
					<td bgcolor="%s" border="0">%.2f%%</td>
				</tr>
			`, index+1, partColor, index+1, partLabel, colors.EvenRow, partLabel, colors.EvenRow, partSize, colors.EvenRow, (float32(partSize)/float32(totalDiskSize))*100)
		}
	}

	// Calcular y agregar el espacio libre general del disco
	freeSpace := totalDiskSize - int32(usedSpace)
	dotData += fmt.Sprintf(`
		<!-- Espacio libre -->
		<tr>
			<td colspan="3" bgcolor="%s" style="rounded" border="0">
				<font color="white" face="Arial" point-size="14"><b>Espacio Libre</b></font>
			</td>
		</tr>
		<tr>
			<td bgcolor="%s" border="0">Libre</td>
			<td bgcolor="%s" border="0">%d bytes</td>
			<td bgcolor="%s" border="0">%.2f%%</td>
		</tr>
	`, colors.FreeSpace, colors.OddRow, colors.OddRow, freeSpace, colors.OddRow, (float32(freeSpace)/float32(totalDiskSize))*100)

	dotData += `</table>>]; 
		
		/* Estilos globales */
		node [fontname="Arial", fontsize=10, shape=box, style="rounded,filled", 
			  fillcolor="#ffffff", color="#2c3e50", penwidth=1.5];
	}`

	// Escribir archivo DOT
	if err := os.WriteFile(dotFile, []byte(dotData), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo DOT: %v", err)
	}

	// Generar imagen con Graphviz (alta calidad)
	cmd := exec.Command("dot", "-Tpng", "-Gdpi=300", "-Nfontname=Arial", "-Efontname=Arial", 
		dotFile, "-o", imgFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando Graphviz: %v", err)
	}

	fmt.Printf("\x1b[32m✓ Reporte del disco generado:\x1b[0m %s\n", imgFile)
	return nil
}