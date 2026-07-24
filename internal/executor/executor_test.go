package executor

import (
	"io"
	"strings"
	"testing"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

func TestBuildExecutesFilterAndProjection(t *testing.T) {
	cat := catalog.New()
	table, err := catalog.LoadCSV("empleados", strings.NewReader("id,nombre,edad,activo\n1,Ana,28,true\n2,Beto,17,true\n3,Carla,35,false\n"))
	if err != nil {
		t.Fatalf("LoadCSV devolvio error: %v", err)
	}
	if err := cat.Add(table); err != nil {
		t.Fatalf("Add devolvio error: %v", err)
	}

	statement, err := query.Parse("SELECT nombre FROM empleados WHERE edad >= 18 AND activo = true")
	if err != nil {
		t.Fatalf("Parse devolvio error: %v", err)
	}
	operator, err := Build(cat, statement)
	if err != nil {
		t.Fatalf("Build devolvio error: %v", err)
	}
	defer operator.Close()

	row, err := operator.Next()
	if err != nil {
		t.Fatalf("Next devolvio error: %v", err)
	}
	if got, want := row[0].Data.(string), "Ana"; got != want {
		t.Errorf("nombre = %q; se esperaba %q", got, want)
	}
	if _, err := operator.Next(); err != io.EOF {
		t.Errorf("el operador debia terminar con io.EOF; recibio %v", err)
	}
}

func TestBuildRejectsUnknownColumns(t *testing.T) {
	cat := catalog.New()
	table, _ := catalog.LoadCSV("empleados", strings.NewReader("id\n1\n"))
	_ = cat.Add(table)
	statement, _ := query.Parse("SELECT nombre FROM empleados")

	if _, err := Build(cat, statement); err == nil {
		t.Fatal("Build permitio una columna inexistente")
	}
}

func TestBuildOrdersAndLimitsRows(t *testing.T) {
	cat := catalog.New()
	table, _ := catalog.LoadCSV("empleados", strings.NewReader("nombre,edad\nAna,28\nBeto,17\nCarla,35\n"))
	_ = cat.Add(table)
	statement, _ := query.Parse("SELECT nombre FROM empleados ORDER BY nombre DESC LIMIT 2")
	operator, err := Build(cat, statement)
	if err != nil {
		t.Fatalf("Build devolvio error: %v", err)
	}
	defer operator.Close()

	for _, want := range []string{"Carla", "Beto"} {
		row, err := operator.Next()
		if err != nil || row[0].Data != want {
			t.Fatalf("fila = %#v, %v; se esperaba %q", row, err, want)
		}
	}
	if _, err := operator.Next(); err != io.EOF {
		t.Errorf("se esperaba io.EOF; se recibio %v", err)
	}
}
