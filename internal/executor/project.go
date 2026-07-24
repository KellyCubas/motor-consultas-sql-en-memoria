package executor

import (
	"fmt"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
)

// Project conserva solo las columnas solicitadas por SELECT.
type Project struct {
	input   Operator
	indices []int
	columns []catalog.Column
}

// NewProject crea una proyeccion de columnas. Un arreglo vacio conserva todas.
func NewProject(input Operator, names []string) (*Project, error) {
	if len(names) == 0 {
		return &Project{input: input, columns: input.Columns()}, nil
	}

	project := &Project{input: input, indices: make([]int, len(names)), columns: make([]catalog.Column, len(names))}
	for nameIndex, name := range names {
		index, column, err := findColumn(input.Columns(), name)
		if err != nil {
			return nil, err
		}
		project.indices[nameIndex] = index
		project.columns[nameIndex] = column
	}
	return project, nil
}

// Next entrega una fila con las columnas proyectadas.
func (p *Project) Next() (catalog.Row, error) {
	row, err := p.input.Next()
	if err != nil || p.indices == nil {
		return row, err
	}

	projected := make(catalog.Row, len(p.indices))
	for outputIndex, inputIndex := range p.indices {
		projected[outputIndex] = row[inputIndex]
	}
	return projected, nil
}

// Columns devuelve el esquema despues de aplicar SELECT.
func (p *Project) Columns() []catalog.Column {
	return p.columns
}

// Close cierra el operador de entrada.
func (p *Project) Close() error {
	return p.input.Close()
}

func findColumn(columns []catalog.Column, name string) (int, catalog.Column, error) {
	for index, column := range columns {
		if equalName(column.Name, name) {
			return index, column, nil
		}
	}
	return 0, catalog.Column{}, fmt.Errorf("la columna %q no existe", name)
}
