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

func TestBuildGroupsAndAggregatesRows(t *testing.T) {
	cat := catalog.New()
	table, _ := catalog.LoadCSV("ventas", strings.NewReader("zona,monto\nNorte,10\nNorte,20\nSur,5\nSur,\n"))
	_ = cat.Add(table)
	statement, err := query.Parse("SELECT zona, COUNT(*), SUM(monto), AVG(monto), MIN(monto), MAX(monto) FROM ventas GROUP BY zona ORDER BY zona")
	if err != nil {
		t.Fatalf("Parse devolvio error: %v", err)
	}
	operator, err := Build(cat, statement)
	if err != nil {
		t.Fatalf("Build devolvio error: %v", err)
	}
	defer operator.Close()

	rows := []catalog.Row{}
	for {
		row, err := operator.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Next devolvio error: %v", err)
		}
		rows = append(rows, row)
	}
	if len(rows) != 2 {
		t.Fatalf("grupos = %d; se esperaban 2", len(rows))
	}
	if rows[0][0].Data != "Norte" || rows[0][1].Data != int64(2) || rows[0][2].Data != 30.0 || rows[0][3].Data != 15.0 || rows[0][4].Data != int64(10) || rows[0][5].Data != int64(20) {
		t.Errorf("agregados Norte incorrectos: %#v", rows[0])
	}
	if rows[1][0].Data != "Sur" || rows[1][1].Data != int64(2) || rows[1][2].Data != 5.0 || rows[1][3].Data != 5.0 {
		t.Errorf("agregados Sur incorrectos: %#v", rows[1])
	}
}

func TestBuildExecutesInnerJoin(t *testing.T) {
	cat := catalog.New()
	empleados, _ := catalog.LoadCSV("empleados", strings.NewReader("nombre,area_id\nAna,1\nBeto,2\n"))
	areas, _ := catalog.LoadCSV("areas", strings.NewReader("id,nombre\n1,Ventas\n2,Soporte\n"))
	_ = cat.Add(empleados)
	_ = cat.Add(areas)
	statement, err := query.Parse("SELECT empleados.nombre, areas.nombre FROM empleados INNER JOIN areas ON empleados.area_id = areas.id ORDER BY empleados.nombre")
	if err != nil {
		t.Fatalf("Parse devolvio error: %v", err)
	}
	operator, err := Build(cat, statement)
	if err != nil {
		t.Fatalf("Build devolvio error: %v", err)
	}
	defer operator.Close()
	row, err := operator.Next()
	if err != nil || row[0].Data != "Ana" || row[1].Data != "Ventas" {
		t.Fatalf("primera fila incorrecta: %#v, %v", row, err)
	}
	row, err = operator.Next()
	if err != nil || row[0].Data != "Beto" || row[1].Data != "Soporte" {
		t.Fatalf("segunda fila incorrecta: %#v, %v", row, err)
	}
}
