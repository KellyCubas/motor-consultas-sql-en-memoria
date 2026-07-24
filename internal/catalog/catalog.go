// Package catalog mantiene las tablas cargadas por el motor en memoria.
package catalog

import (
	"fmt"
	"strings"
)

// Catalog almacena tablas por nombre.
type Catalog struct {
	tables map[string]*Table
}

// New crea un catalogo vacio.
func New() *Catalog {
	return &Catalog{tables: make(map[string]*Table)}
}

// Add registra una tabla. Los nombres no distinguen mayusculas de minusculas.
func (c *Catalog) Add(table *Table) error {
	if table == nil {
		return fmt.Errorf("no se puede agregar una tabla nula")
	}

	key := normalizeName(table.Name)
	if key == "" {
		return fmt.Errorf("el nombre de la tabla no puede estar vacio")
	}
	if _, exists := c.tables[key]; exists {
		return fmt.Errorf("la tabla %q ya existe", table.Name)
	}

	c.tables[key] = table
	return nil
}

// Table devuelve una tabla por nombre.
func (c *Catalog) Table(name string) (*Table, bool) {
	table, ok := c.tables[normalizeName(name)]
	return table, ok
}

// Tables devuelve todas las tablas registradas.
func (c *Catalog) Tables() []*Table {
	tables := make([]*Table, 0, len(c.tables))
	for _, table := range c.tables {
		tables = append(tables, table)
	}
	return tables
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
