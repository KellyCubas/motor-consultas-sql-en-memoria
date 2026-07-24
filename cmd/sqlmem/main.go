package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/executor"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "cargar":
		load(os.Args[2:])
	case "consultar":
		queryCSV(os.Args[2:])
	default:
		printUsage()
	}
}

func load(arguments []string) {
	if len(arguments) != 2 {
		fmt.Fprintln(os.Stderr, "Uso: sqlmem cargar <tabla> <archivo.csv>")
		return
	}
	table, err := catalog.LoadCSVFile(arguments[0], arguments[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	fmt.Printf("Tabla %q cargada: %d filas\n", table.Name, len(table.Rows))
	for _, column := range table.Columns {
		fmt.Printf("- %s: %s\n", column.Name, column.Type)
	}
}

func queryCSV(arguments []string) {
	if len(arguments) != 3 {
		queryMultipleCSV(arguments)
		return
	}
	table, err := catalog.LoadCSVFile(arguments[0], arguments[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	cat := catalog.New()
	if err := cat.Add(table); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	statement, err := query.Parse(arguments[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	operator, err := executor.Build(cat, statement)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer operator.Close()

	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, column := range operator.Columns() {
		fmt.Fprintf(writer, "%s\t", column.Name)
	}
	fmt.Fprintln(writer)
	for {
		row, err := operator.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		for _, value := range row {
			if value.Null {
				fmt.Fprint(writer, "NULL\t")
				continue
			}
			fmt.Fprintf(writer, "%v\t", value.Data)
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()
}

func queryMultipleCSV(arguments []string) {
	separator := -1
	for index, argument := range arguments {
		if argument == "--" {
			separator = index
			break
		}
	}
	if separator < 1 || separator != len(arguments)-2 {
		fmt.Fprintln(os.Stderr, "Uso: sqlmem consultar <tabla=archivo.csv> [tabla=archivo.csv ...] -- <consulta SQL>")
		return
	}

	cat := catalog.New()
	for _, source := range arguments[:separator] {
		parts := strings.SplitN(source, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			fmt.Fprintf(os.Stderr, "Error: fuente invalida %q; use tabla=archivo.csv\n", source)
			return
		}
		table, err := catalog.LoadCSVFile(parts[0], parts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		if err := cat.Add(table); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
	}
	statement, err := query.Parse(arguments[separator+1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	runQuery(cat, statement)
}

func runQuery(catalog *catalog.Catalog, statement *query.Query) {
	operator, err := executor.Build(catalog, statement)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer operator.Close()
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, column := range operator.Columns() {
		fmt.Fprintf(writer, "%s\t", column.Name)
	}
	fmt.Fprintln(writer)
	for {
		row, err := operator.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		for _, value := range row {
			if value.Null {
				fmt.Fprint(writer, "NULL\t")
			} else {
				fmt.Fprintf(writer, "%v\t", value.Data)
			}
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Uso:")
	fmt.Fprintln(os.Stderr, "  sqlmem cargar <tabla> <archivo.csv>")
	fmt.Fprintln(os.Stderr, "  sqlmem consultar <tabla> <archivo.csv> <consulta SQL>")
	fmt.Fprintln(os.Stderr, "  sqlmem consultar <tabla=archivo.csv> [tabla=archivo.csv ...] -- <consulta SQL>")
}
