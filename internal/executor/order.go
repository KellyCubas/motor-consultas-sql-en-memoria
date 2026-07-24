package executor

import (
	"fmt"
	"io"
	"sort"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

// Order materializa y ordena las filas del operador de entrada.
type Order struct {
	input   Operator
	terms   []query.OrderTerm
	indices []int
	rows    []catalog.Row
	index   int
	loaded  bool
}

// NewOrder crea un operador de ordenamiento.
func NewOrder(input Operator, terms []query.OrderTerm) (*Order, error) {
	order := &Order{input: input, terms: terms, indices: make([]int, len(terms))}
	for termIndex, term := range terms {
		index, _, err := findColumn(input.Columns(), term.Column)
		if err != nil {
			return nil, err
		}
		order.indices[termIndex] = index
	}
	return order, nil
}

func (o *Order) Next() (catalog.Row, error) {
	if !o.loaded {
		if err := o.load(); err != nil {
			return nil, err
		}
	}
	if o.index >= len(o.rows) {
		return nil, io.EOF
	}
	row := o.rows[o.index]
	o.index++
	return row, nil
}

func (o *Order) Columns() []catalog.Column { return o.input.Columns() }
func (o *Order) Close() error              { return o.input.Close() }

func (o *Order) load() error {
	for {
		row, err := o.input.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		o.rows = append(o.rows, row)
	}
	var sortErr error
	sort.SliceStable(o.rows, func(leftIndex, rightIndex int) bool {
		if sortErr != nil {
			return false
		}
		left, right := o.rows[leftIndex], o.rows[rightIndex]
		for termIndex, columnIndex := range o.indices {
			comparison, err := compareForOrder(left[columnIndex], right[columnIndex])
			if err != nil {
				sortErr = err
				return false
			}
			if comparison == 0 {
				continue
			}
			if o.terms[termIndex].Descending {
				return comparison > 0
			}
			return comparison < 0
		}
		return false
	})
	if sortErr != nil {
		return sortErr
	}
	o.loaded = true
	return nil
}

func compareForOrder(left, right catalog.Value) (int, error) {
	if left.Null && right.Null {
		return 0, nil
	}
	if left.Null {
		return 1, nil
	}
	if right.Null {
		return -1, nil
	}
	comparison, err := compare(left, right)
	if err != nil {
		return 0, fmt.Errorf("ordenar: %w", err)
	}
	return comparison, nil
}
