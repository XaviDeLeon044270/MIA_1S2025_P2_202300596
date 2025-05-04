package reports

import (
	"fmt"
	"html"
	"os"
	"os/exec"
	"sort"
	"strings"

	structures "backend/structures"
	utils "backend/utils"
)

type blockInfo struct {
	index       int32
	btype       string // "folder", "file" o "pointer"
	folderData  *structures.FolderBlock
	fileData    *structures.FileBlock
	pointerData *structures.PointerBlock
}

// ColorPalette define la paleta de colores profesional
type ColorPalette struct {
	Background string
	Primary    string
	Folder     string
	File       string
	Pointer    string
	Text       string
	Accent     string
	EvenRow    string
	OddRow     string
}

func ReportBlock(sb *structures.SuperBlock, diskPath, outPath string) error {
	// Crear la carpeta de salida si no existe
	if err := utils.CreateParentDirs(outPath); err != nil {
		return fmt.Errorf("error al crear directorios: %v", err)
	}

	dotFileName, outputImage := utils.GetFileNames(outPath)

	// Paleta de colores profesional
	colors := ColorPalette{
		Background: "#f5f5f5",
		Primary:    "#2c3e50",
		Folder:     "#3498db",
		File:       "#2ecc71",
		Pointer:    "#e74c3c",
		Text:       "#333333",
		Accent:     "#9b59b6",
		EvenRow:    "#ecf0f1",
		OddRow:     "#ffffff",
	}

	// Se almacenan los bloques descubiertos
	visited := make(map[int32]*blockInfo)

	// Recorrer todos los inodos
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &structures.Inode{}
		offset := sb.S_inode_start + i*sb.S_inode_size
		if err := inode.Deserialize(diskPath, int64(offset)); err != nil {
			continue
		}

		// Procesar bloques directos
		for b := 0; b < 12; b++ {
			if inode.I_block[b] == -1 {
				continue
			}
			processBlock(sb, diskPath, inode.I_block[b], inode.I_type[0], visited)
		}

		// Procesar bloques indirectos
		if inode.I_block[12] != -1 {
			handlePointerBlock(sb, diskPath, inode.I_block[12], visited)
		}
	}

	// Construir el contenido DOT
	dotContent := fmt.Sprintf(`digraph G {
		rankdir=LR;
		bgcolor="%s";
		node [shape=plaintext, fontname="Arial", fontsize=12];
		edge [color="%s", arrowsize=0.8];
		
		/* Estilo global */
		graph [fontname="Arial", fontsize=12];
		node [fontname="Arial", fontsize=10];
	`, colors.Background, colors.Primary)

	// Ordenar los índices de los bloques
	var sortedIndices []int32
	for k := range visited {
		sortedIndices = append(sortedIndices, k)
	}
	sort.Slice(sortedIndices, func(i, j int) bool {
		return sortedIndices[i] < sortedIndices[j]
	})

	// Generar nodos y enlaces
	var prev string
	for _, idx := range sortedIndices {
		info := visited[idx]
		nodeName := fmt.Sprintf("block%d", idx)

		switch info.btype {
		case "folder":
			dotContent += renderFolderBlock(nodeName, idx, info.folderData, colors)
		case "file":
			dotContent += renderFileBlock(nodeName, idx, info.fileData, colors)
		case "pointer":
			dotContent += renderPointerBlock(nodeName, idx, info.pointerData, colors)
		}

		// Enlazar bloques
		if prev != "" {
			dotContent += fmt.Sprintf("%s -> %s [color=\"%s\", penwidth=1.5, style=\"dashed\"];\n", 
				prev, nodeName, colors.Accent)
		}
		prev = nodeName
	}

	dotContent += "\n}\n"

	// Escribir el archivo .dot
	if err := writeDotFile(dotFileName, dotContent); err != nil {
		return fmt.Errorf("error al escribir archivo DOT: %v", err)
	}

	// Generar imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", "-Gdpi=300", "-Nfontname=Arial", 
		"-Efontname=Arial", dotFileName, "-o", outputImage)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al generar imagen: %v", err)
	}

	fmt.Printf("\x1b[32m✓ Reporte de bloques generado:\x1b[0m %s\n", outputImage)
	return nil
}

// renderFolderBlock genera el DOT para bloques de carpeta
func renderFolderBlock(nodeName string, idx int32, fb *structures.FolderBlock, colors ColorPalette) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`%s [label=<
	<table border="0" cellborder="1" cellspacing="0" cellpadding="8" style="rounded" bgcolor="%s">
		<tr>
			<td colspan="2" bgcolor="%s" style="rounded" border="0">
				<font color="white" face="Arial" point-size="14"><b>Bloque Carpeta %d</b></font>
			</td>
		</tr>
		<tr>
			<td bgcolor="%s" border="0"><font color="white"><b>Nombre</b></font></td>
			<td bgcolor="%s" border="0"><font color="white"><b>Inodo</b></font></td>
		</tr>
	`, nodeName, colors.EvenRow, colors.Folder, idx, colors.Primary, colors.Primary))

	for i, c := range fb.B_content {
		name := strings.Trim(string(c.B_name[:]), "\x00 ")
		if name == "" {
			continue
		}
		
		rowColor := colors.EvenRow
		if i%2 == 0 {
			rowColor = colors.OddRow
		}
		
		sb.WriteString(fmt.Sprintf(`
		<tr>
			<td bgcolor="%s" border="0">%s</td>
			<td bgcolor="%s" border="0">%d</td>
		</tr>`, rowColor, name, rowColor, c.B_inodo))
	}

	sb.WriteString(`</table>>];`)
	return sb.String()
}

// renderFileBlock genera el DOT para bloques de archivo
func renderFileBlock(nodeName string, idx int32, fb *structures.FileBlock, colors ColorPalette) string {
	fileContent := strings.Trim(string(fb.B_content[:]), "\x00 ")
	lines := strings.Split(fileContent, "\n")
	formattedContent := ""
	
	for _, line := range lines {
		if line != "" {
			formattedContent += fmt.Sprintf("%s<br align='left'/>", html.EscapeString(strings.TrimSpace(line)))
		}
	}

	return fmt.Sprintf(`%s [label=<
	<table border="0" cellborder="1" cellspacing="0" cellpadding="8" style="rounded" bgcolor="%s">
		<tr>
			<td bgcolor="%s" style="rounded" border="0">
				<font color="white" face="Arial" point-size="14"><b>Bloque Archivo %d</b></font>
			</td>
		</tr>
		<tr>
			<td bgcolor="%s" border="0" align="left">%s</td>
		</tr>
	</table>>];
`, nodeName, colors.EvenRow, colors.File, idx, colors.OddRow, formattedContent)
}

// renderPointerBlock genera el DOT para bloques de punteros
func renderPointerBlock(nodeName string, idx int32, pb *structures.PointerBlock, colors ColorPalette) string {
	pointers := make([]string, 0)
	for _, ptr := range pb.P_pointers {
		if ptr != -1 {
			pointers = append(pointers, fmt.Sprintf("%d", ptr))
		}
	}

	return fmt.Sprintf(`%s [label=<
	<table border="0" cellborder="1" cellspacing="0" cellpadding="8" style="rounded" bgcolor="%s">
		<tr>
			<td bgcolor="%s" style="rounded" border="0">
				<font color="white" face="Arial" point-size="14"><b>Bloque Punteros %d</b></font>
			</td>
		</tr>
		<tr>
			<td bgcolor="%s" border="0">%s</td>
		</tr>
	</table>>];
`, nodeName, colors.EvenRow, colors.Pointer, idx, colors.OddRow, strings.Join(pointers, ", "))
}

// processBlock procesa un bloque individual
func processBlock(sb *structures.SuperBlock, diskPath string, blkIndex int32, inodeType byte, visited map[int32]*blockInfo) {
	if _, ok := visited[blkIndex]; ok {
		return
	}

	blockOffset := sb.S_block_start + blkIndex*sb.S_block_size
	
	if inodeType == '0' { // Carpeta
		fb := &structures.FolderBlock{}
		if err := fb.Deserialize(diskPath, int64(blockOffset)); err == nil {
			visited[blkIndex] = &blockInfo{
				index:      blkIndex,
				btype:      "folder",
				folderData: fb,
			}
		}
	} else if inodeType == '1' { // Archivo
		fb := &structures.FileBlock{}
		if err := fb.Deserialize(diskPath, int64(blockOffset)); err == nil {
			visited[blkIndex] = &blockInfo{
				index:    blkIndex,
				btype:    "file",
				fileData: fb,
			}
		}
	}
}

// handlePointerBlock procesa bloques de punteros recursivamente
func handlePointerBlock(sb *structures.SuperBlock, diskPath string, blockIndex int32, visited map[int32]*blockInfo) {
	if _, ok := visited[blockIndex]; ok {
		return
	}

	pb := &structures.PointerBlock{}
	blockOffset := sb.S_block_start + blockIndex*sb.S_block_size
	if err := pb.Deserialize(diskPath, int64(blockOffset)); err != nil {
		return
	}

	visited[blockIndex] = &blockInfo{
		index:       blockIndex,
		btype:       "pointer",
		pointerData: pb,
	}

	// Procesar bloques apuntados
	for _, ptr := range pb.P_pointers {
		if ptr != -1 {
			processBlock(sb, diskPath, ptr, '0', visited) // Asumimos que apuntan a carpetas
		}
	}
}

// writeDotFile escribe el contenido DOT en un archivo
func writeDotFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}