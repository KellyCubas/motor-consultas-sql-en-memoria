package executor

import (
	"fmt"
	"strings"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

// Build crea el arbol de operadores para una consulta ya analizada.
func Build(catalog *catalog.Catalog, statement *query.Query) (Operator, error) {
	table, ok := catalog.Table(statement.Table)
	if !ok {
		return nil, fmt.Errorf("la tabla %q no existe", statement.Table)
	}

	var operator Operator = NewScan(table)
	var err error
	if statement.Where != nil {
		if err := validateExpression(statement.Where, table.Columns); err != nil {
			return nil, err
		}
		operator = NewFilter(operator, statement.Where)
	}

	if statement.SelectAll {
		operator, err = NewProject(operator, nil)
	} else {
		operator, err = NewProject(operator, statement.Columns)
	}
	if err != nil {
		return nil, err
	}
	if len(statement.OrderBy) > 0 {
		operator, err = NewOrder(operator, statement.OrderBy)
		if err != nil {
			return nil, err
		}
	}
	if statement.Limit != nil {
		operator = NewLimit(operator, *statement.Limit)
	}
	return operator, nil
}

func validateExpression(expression query.Expression, columns []catalog.Column) error {
	switch expression := expression.(type) {
	case query.Identifier:
		_, _, err := findColumn(columns, expression.Name)
		return err
	case query.Literal:
		return nil
	case query.BinaryExpression:
		if err := validateExpression(expression.Left, columns); err != nil {
			return err
		}
		return validateExpression(expression.Right, columns)
	default:
		return fmt.Errorf("expresion no soportada")
	}
}

func equalName(left, right string) bool {
	return strings.EqualFold(strings.TrimSpace(left), strings.TrimSpace(right))
}
