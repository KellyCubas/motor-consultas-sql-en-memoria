package executor

import (
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

// Filter entrega solo las filas que cumplen una condicion WHERE.
type Filter struct {
	input     Operator
	condition query.Expression
}

// NewFilter crea un filtro sobre un operador existente.
func NewFilter(input Operator, condition query.Expression) *Filter {
	return &Filter{input: input, condition: condition}
}

// Next busca y entrega la siguiente fila que cumple la condicion.
func (f *Filter) Next() (catalog.Row, error) {
	for {
		row, err := f.input.Next()
		if err != nil {
			return nil, err
		}

		matches, err := EvaluatePredicate(f.condition, row, f.input.Columns())
		if err != nil {
			return nil, err
		}
		if matches {
			return row, nil
		}
	}
}

// Columns conserva el esquema del operador de entrada.
func (f *Filter) Columns() []catalog.Column {
	return f.input.Columns()
}

// Close cierra el operador de entrada.
func (f *Filter) Close() error {
	return f.input.Close()
}

var _ Operator = (*Filter)(nil)
