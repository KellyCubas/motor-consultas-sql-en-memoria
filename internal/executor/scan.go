package executor

import (
	"io"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
)

// Scan recorre secuencialmente las filas de una tabla en memoria.
type Scan struct {
	table *catalog.Table
	index int
}

// NewScan crea un escaneo para una tabla.
func NewScan(table *catalog.Table) *Scan {
	return &Scan{table: table}
}

// Next entrega la siguiente fila de la tabla.
func (s *Scan) Next() (catalog.Row, error) {
	if s.index >= len(s.table.Rows) {
		return nil, io.EOF
	}
	row := s.table.Rows[s.index]
	s.index++
	return row, nil
}

// Columns devuelve el esquema de la tabla.
func (s *Scan) Columns() []catalog.Column {
	return s.table.Columns
}

// Close no requiere liberar recursos para una tabla en memoria.
func (s *Scan) Close() error {
	return nil
}
