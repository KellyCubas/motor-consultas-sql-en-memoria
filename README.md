# Motor de consultas SQL en memoria

Proyecto para el Taller de Programacion en Go. Implementa un subconjunto de SQL que consulta tablas cargadas desde archivos CSV y mantenidas en memoria.

## Estado

Hito 1 iniciado: carga de CSV, catalogo de tablas, inferencia de tipos y manejo de valores `NULL`.
Hitos 1 a 4 completados: ordenamiento, limite, agrupacion y agregados.

## Estructura

```text
cmd/sqlmem/  Punto de entrada del CLI o REPL.
internal/    Implementacion interna del motor.
data/        Archivos CSV de ejemplo.
docs/        Gramatica y decisiones de diseno.
```

## Requisitos

- Go 1.24 o superior.

## Comandos de desarrollo

```bash
go build ./...
go test ./...
go vet ./...
```

## Ejecucion actual

```bash
go run ./cmd/sqlmem cargar empleados data/empleados.csv
```

Salida esperada:

```text
Tabla "empleados" cargada: 3 filas
- id: entero
- nombre: texto
- edad: entero
- salario: decimal
- activo: booleano
```

Consultar los datos cargados desde un CSV:

```bash
go run ./cmd/sqlmem consultar empleados data/empleados.csv "SELECT nombre, salario FROM empleados WHERE activo = true AND edad >= 25"
```

Ordenar y limitar resultados:

```bash
go run ./cmd/sqlmem consultar empleados data/empleados.csv "SELECT nombre, salario FROM empleados ORDER BY salario DESC LIMIT 2"
```

Agrupar resultados y usar agregados:

```bash
go run ./cmd/sqlmem consultar empleados data/empleados.csv "SELECT activo, COUNT(*), AVG(salario) FROM empleados GROUP BY activo ORDER BY activo"
```

## Alcance previsto

1. Carga de CSV como tablas en memoria con catalogo, esquemas y tipos.
2. Lexer, parser y AST para `SELECT ... FROM ... WHERE ...`.
3. Operadores `Scan`, `Filter` y `Project` mediante el modelo Volcano. Completado.
4. `ORDER BY`, `LIMIT`, `GROUP BY` y agregados. Completado.
5. `INNER JOIN` con nested-loop y hash join. En progreso: las estrategias estan implementadas en el motor; falta exponer la carga de multiples CSV desde el CLI.
