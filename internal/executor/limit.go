package executor

import (
	"io"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
)

// Limit detiene la entrega de filas al alcanzar el limite indicado.
type Limit struct {
	input Operator
	max   int
	count int
}

func NewLimit(input Operator, max int) *Limit { return &Limit{input: input, max: max} }

func (l *Limit) Next() (catalog.Row, error) {
	if l.count >= l.max {
		return nil, io.EOF
	}
	row, err := l.input.Next()
	if err == nil {
		l.count++
	}
	return row, err
}

func (l *Limit) Columns() []catalog.Column { return l.input.Columns() }
func (l *Limit) Close() error              { return l.input.Close() }
