package main

import (
	"fmt"
	"os"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
)

func main() {
	if len(os.Args) != 4 || os.Args[1] != "cargar" {
		fmt.Fprintln(os.Stderr, "Uso: sqlmem cargar <tabla> <archivo.csv>")
		return
	}

	table, err := catalog.LoadCSVFile(os.Args[2], os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	fmt.Printf("Tabla %q cargada: %d filas\n", table.Name, len(table.Rows))
	for _, column := range table.Columns {
		fmt.Printf("- %s: %s\n", column.Name, column.Type)
	}
}
