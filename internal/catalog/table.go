package catalog

import "fmt"

// DataType representa los tipos soportados por el motor.
type DataType uint8

const (
	Text DataType = iota
	Integer
	Decimal
	Boolean
)

func (t DataType) String() string {
	switch t {
	case Integer:
		return "entero"
	case Decimal:
		return "decimal"
	case Boolean:
		return "booleano"
	default:
		return "texto"
	}
}

// Value guarda un valor tipado. Null indica que no tiene valor.
type Value struct {
	Type DataType
	Data any
	Null bool
}

// Column describe una columna de una tabla.
type Column struct {
	Name string
	Type DataType
}

// Row contiene los valores de una fila, en el orden del esquema.
type Row []Value

// Table es una tabla cargada por completo en memoria.
type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
}

// ColumnIndex devuelve la posicion de una columna.
func (t *Table) ColumnIndex(name string) (int, error) {
	for index, column := range t.Columns {
		if normalizeName(column.Name) == normalizeName(name) {
			return index, nil
		}
	}
	return 0, fmt.Errorf("la columna %q no existe en la tabla %q", name, t.Name)
}
