// Package executor construye y recorre el arbol de operadores de una consulta.
package executor

import "github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"

// Operator entrega filas una por una. Retorna io.EOF cuando no hay mas filas.
type Operator interface {
	Next() (catalog.Row, error)
	Columns() []catalog.Column
	Close() error
}
