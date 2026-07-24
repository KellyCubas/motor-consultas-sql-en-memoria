package executor

import (
	"fmt"
	"io"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

// HashJoin indexa las filas derechas para joins de igualdad.
type HashJoin struct {
	left       Operator
	rows       map[string][]catalog.Row
	leftIndex  int
	rightIndex int
	current    catalog.Row
	matches    []catalog.Row
	columns    []catalog.Column
}

func NewHashJoin(left, right Operator, leftName, rightName string, condition query.Expression) (*HashJoin, error) {
	expr, ok := condition.(query.BinaryExpression)
	if !ok || expr.Operator != query.EqualToken {
		return nil, fmt.Errorf("hash join requiere una condicion de igualdad")
	}
	l, ok := expr.Left.(query.Identifier)
	if !ok {
		return nil, fmt.Errorf("hash join requiere columnas")
	}
	r, ok := expr.Right.(query.Identifier)
	if !ok {
		return nil, fmt.Errorf("hash join requiere columnas")
	}
	cols := qualify(left.Columns(), leftName)
	rightCols := qualify(right.Columns(), rightName)
	li, _, err := findColumn(cols, l.Name)
	if err != nil {
		return nil, err
	}
	ri, _, err := findColumn(rightCols, r.Name)
	if err != nil {
		return nil, err
	}
	index := map[string][]catalog.Row{}
	for {
		row, e := right.Next()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		if !row[ri].Null {
			key := fmt.Sprintf("%d:%v", row[ri].Type, row[ri].Data)
			index[key] = append(index[key], row)
		}
	}
	return &HashJoin{left: left, rows: index, leftIndex: li, rightIndex: 0, columns: append(cols, rightCols...)}, nil
}
func (j *HashJoin) Next() (catalog.Row, error) {
	for {
		if j.current == nil {
			row, err := j.left.Next()
			if err != nil {
				return nil, err
			}
			j.current = row
			j.rightIndex = 0
			j.matches = nil
			if !row[j.leftIndex].Null {
				j.matches = j.rows[fmt.Sprintf("%d:%v", row[j.leftIndex].Type, row[j.leftIndex].Data)]
			}
		}
		if j.rightIndex < len(j.matches) {
			right := j.matches[j.rightIndex]
			j.rightIndex++
			return append(append(catalog.Row{}, j.current...), right...), nil
		}
		j.current = nil
	}
}
func (j *HashJoin) Columns() []catalog.Column { return j.columns }
func (j *HashJoin) Close() error              { return j.left.Close() }
