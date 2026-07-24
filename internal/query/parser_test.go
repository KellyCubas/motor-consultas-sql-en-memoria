package query

import (
	"strings"
	"testing"
)

func TestParseSelectWithWhere(t *testing.T) {
	query, err := Parse("SELECT nombre, edad FROM empleados WHERE (edad >= 18 AND activo = true) OR nombre = 'Ana'")
	if err != nil {
		t.Fatalf("Parse devolvio error: %v", err)
	}

	if query.SelectAll || query.Table != "empleados" {
		t.Fatalf("consulta analizada incorrectamente: %#v", query)
	}
	if got, want := strings.Join(query.Columns, ","), "nombre,edad"; got != want {
		t.Errorf("columnas = %q; se esperaban %q", got, want)
	}

	or, ok := query.Where.(BinaryExpression)
	if !ok || or.Operator != OrToken {
		t.Fatalf("WHERE no conserva OR: %#v", query.Where)
	}
	and, ok := or.Left.(BinaryExpression)
	if !ok || and.Operator != AndToken {
		t.Fatalf("WHERE no conserva AND: %#v", or.Left)
	}
}

func TestParseRejectsInvalidQueries(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "sin FROM", input: "SELECT nombre empleados"},
		{name: "sin operador", input: "SELECT * FROM empleados WHERE edad 18"},
		{name: "parentesis sin cerrar", input: "SELECT * FROM empleados WHERE (edad = 18"},
		{name: "texto sin cerrar", input: "SELECT * FROM empleados WHERE nombre = 'Ana"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse(test.input)
			if err == nil {
				t.Fatal("Parse no devolvio error")
			}
			if !strings.Contains(err.Error(), "posicion") {
				t.Errorf("el error no incluye posicion: %v", err)
			}
		})
	}
}

func TestParseOrderByAndLimit(t *testing.T) {
	query, err := Parse("SELECT nombre FROM empleados ORDER BY nombre DESC LIMIT 2")
	if err != nil {
		t.Fatalf("Parse devolvio error: %v", err)
	}
	if len(query.OrderBy) != 1 || !query.OrderBy[0].Descending {
		t.Fatalf("ORDER BY no se analizo correctamente: %#v", query.OrderBy)
	}
	if query.Limit == nil || *query.Limit != 2 {
		t.Fatalf("LIMIT no se analizo correctamente: %#v", query.Limit)
	}
}

func TestParseInnerJoin(t *testing.T) {
	query, err := Parse("SELECT empleados.nombre FROM empleados INNER JOIN areas ON empleados.area_id = areas.id")
	if err != nil {
		t.Fatalf("Parse devolvio error: %v", err)
	}
	if query.Join == nil || query.Join.Table != "areas" {
		t.Fatalf("JOIN no se analizo correctamente: %#v", query.Join)
	}
}
