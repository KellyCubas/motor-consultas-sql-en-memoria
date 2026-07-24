// Package query analiza el subconjunto SQL soportado por el motor.
package query

// Query representa una consulta SELECT analizada.
type Query struct {
	SelectAll  bool
	Columns    []string
	Table      string
	Where      Expression
	OrderBy    []OrderTerm
	Limit      *int
	GroupBy    []string
	Aggregates []Aggregate
}

type Aggregate struct {
	Function TokenKind
	Column   string
	Star     bool
}

// OrderTerm describe una columna y su direccion de ordenamiento.
type OrderTerm struct {
	Column     string
	Descending bool
}

// Expression es un nodo del arbol de condiciones WHERE.
type Expression interface {
	expression()
}

// Identifier representa el nombre de una columna.
type Identifier struct {
	Name string
}

func (Identifier) expression() {}

// Literal representa un texto, numero, booleano o NULL.
type Literal struct {
	Value string
	Kind  TokenKind
}

func (Literal) expression() {}

// BinaryExpression representa una comparacion o una operacion logica.
type BinaryExpression struct {
	Left     Expression
	Operator TokenKind
	Right    Expression
}

func (BinaryExpression) expression() {}
