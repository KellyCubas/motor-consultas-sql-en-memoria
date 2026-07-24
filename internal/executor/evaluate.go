package executor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

// EvaluatePredicate evalua una condicion WHERE. Un resultado NULL no selecciona la fila.
func EvaluatePredicate(expression query.Expression, row catalog.Row, columns []catalog.Column) (bool, error) {
	value, err := evaluate(expression, row, columns)
	if err != nil {
		return false, err
	}
	if value.Null {
		return false, nil
	}
	if value.Type != catalog.Boolean {
		return false, fmt.Errorf("la condicion WHERE no produce un booleano")
	}
	return value.Data.(bool), nil
}

func evaluate(expression query.Expression, row catalog.Row, columns []catalog.Column) (catalog.Value, error) {
	switch expression := expression.(type) {
	case query.Identifier:
		index, _, err := findColumn(columns, expression.Name)
		if err != nil {
			return catalog.Value{}, err
		}
		return row[index], nil
	case query.Literal:
		return literalValue(expression)
	case query.BinaryExpression:
		if expression.Operator == query.AndToken || expression.Operator == query.OrToken {
			return logicalValue(expression, row, columns)
		}
		return comparisonValue(expression, row, columns)
	default:
		return catalog.Value{}, fmt.Errorf("expresion no soportada")
	}
}

func literalValue(literal query.Literal) (catalog.Value, error) {
	switch literal.Kind {
	case query.StringToken:
		return catalog.Value{Type: catalog.Text, Data: literal.Value}, nil
	case query.BooleanToken:
		value, err := strconv.ParseBool(strings.ToLower(literal.Value))
		return catalog.Value{Type: catalog.Boolean, Data: value}, err
	case query.NumberToken:
		if integer, err := strconv.ParseInt(literal.Value, 10, 64); err == nil {
			return catalog.Value{Type: catalog.Integer, Data: integer}, nil
		}
		decimal, err := strconv.ParseFloat(literal.Value, 64)
		return catalog.Value{Type: catalog.Decimal, Data: decimal}, err
	case query.NullToken:
		return catalog.Value{Null: true}, nil
	default:
		return catalog.Value{}, fmt.Errorf("literal no soportado")
	}
}

func logicalValue(expression query.BinaryExpression, row catalog.Row, columns []catalog.Column) (catalog.Value, error) {
	left, err := evaluate(expression.Left, row, columns)
	if err != nil {
		return catalog.Value{}, err
	}
	right, err := evaluate(expression.Right, row, columns)
	if err != nil {
		return catalog.Value{}, err
	}
	if (left.Null || left.Type != catalog.Boolean) && !left.Null {
		return catalog.Value{}, fmt.Errorf("el operando izquierdo de %s no es booleano", expression.Operator)
	}
	if (right.Null || right.Type != catalog.Boolean) && !right.Null {
		return catalog.Value{}, fmt.Errorf("el operando derecho de %s no es booleano", expression.Operator)
	}

	if expression.Operator == query.AndToken {
		if (!left.Null && !left.Data.(bool)) || (!right.Null && !right.Data.(bool)) {
			return catalog.Value{Type: catalog.Boolean, Data: false}, nil
		}
	} else if (!left.Null && left.Data.(bool)) || (!right.Null && right.Data.(bool)) {
		return catalog.Value{Type: catalog.Boolean, Data: true}, nil
	}
	if left.Null || right.Null {
		return catalog.Value{Type: catalog.Boolean, Null: true}, nil
	}
	return catalog.Value{Type: catalog.Boolean, Data: expression.Operator == query.AndToken}, nil
}

func comparisonValue(expression query.BinaryExpression, row catalog.Row, columns []catalog.Column) (catalog.Value, error) {
	left, err := evaluate(expression.Left, row, columns)
	if err != nil {
		return catalog.Value{}, err
	}
	right, err := evaluate(expression.Right, row, columns)
	if err != nil {
		return catalog.Value{}, err
	}
	if left.Null || right.Null {
		return catalog.Value{Type: catalog.Boolean, Null: true}, nil
	}

	comparison, err := compare(left, right)
	if err != nil {
		return catalog.Value{}, err
	}
	result := false
	switch expression.Operator {
	case query.EqualToken:
		result = comparison == 0
	case query.NotEqualToken:
		result = comparison != 0
	case query.LessToken:
		result = comparison < 0
	case query.GreaterToken:
		result = comparison > 0
	case query.LessEqualToken:
		result = comparison <= 0
	case query.GreaterEqualToken:
		result = comparison >= 0
	}
	return catalog.Value{Type: catalog.Boolean, Data: result}, nil
}

func compare(left, right catalog.Value) (int, error) {
	if isNumber(left.Type) && isNumber(right.Type) {
		leftNumber := asFloat(left)
		rightNumber := asFloat(right)
		if leftNumber < rightNumber {
			return -1, nil
		}
		if leftNumber > rightNumber {
			return 1, nil
		}
		return 0, nil
	}
	if left.Type != right.Type {
		return 0, fmt.Errorf("no se puede comparar %s con %s", left.Type, right.Type)
	}

	switch left.Type {
	case catalog.Text:
		return strings.Compare(left.Data.(string), right.Data.(string)), nil
	case catalog.Boolean:
		leftBoolean, rightBoolean := left.Data.(bool), right.Data.(bool)
		if leftBoolean == rightBoolean {
			return 0, nil
		}
		if !leftBoolean {
			return -1, nil
		}
		return 1, nil
	default:
		return 0, fmt.Errorf("tipo no soportado: %s", left.Type)
	}
}

func isNumber(dataType catalog.DataType) bool {
	return dataType == catalog.Integer || dataType == catalog.Decimal
}

func asFloat(value catalog.Value) float64 {
	if value.Type == catalog.Integer {
		return float64(value.Data.(int64))
	}
	return value.Data.(float64)
}
