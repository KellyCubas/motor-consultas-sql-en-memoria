package executor

import (
	"io"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

// NestedLoopJoin compara cada fila izquierda contra las filas derechas.
type NestedLoopJoin struct {
	left      Operator
	rightRows []catalog.Row
	condition query.Expression
	columns   []catalog.Column
	current   catalog.Row
	index     int
}

func NewNestedLoopJoin(left, right Operator, leftName, rightName string, condition query.Expression) (*NestedLoopJoin, error) {
	rightRows := []catalog.Row{}
	for {
		row, err := right.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		rightRows = append(rightRows, row)
	}
	columns := qualify(left.Columns(), leftName)
	columns = append(columns, qualify(right.Columns(), rightName)...)
	return &NestedLoopJoin{left: left, rightRows: rightRows, condition: condition, columns: columns}, nil
}

func (j *NestedLoopJoin) Next() (catalog.Row, error) {
	for {
		if j.current == nil {
			row, err := j.left.Next()
			if err != nil {
				return nil, err
			}
			j.current = row
			j.index = 0
		}
		for j.index < len(j.rightRows) {
			right := j.rightRows[j.index]
			j.index++
			row := append(append(catalog.Row{}, j.current...), right...)
			matches, err := EvaluatePredicate(j.condition, row, j.columns)
			if err != nil {
				return nil, err
			}
			if matches {
				return row, nil
			}
		}
		j.current = nil
	}
}
func (j *NestedLoopJoin) Columns() []catalog.Column { return j.columns }
func (j *NestedLoopJoin) Close() error              { return j.left.Close() }
func qualify(columns []catalog.Column, table string) []catalog.Column {
	result := make([]catalog.Column, len(columns))
	for i, column := range columns {
		column.Name = table + "." + column.Name
		result[i] = column
	}
	return result
}
