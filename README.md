# Motor de consultas SQL en memoria

Proyecto para el Taller de Programacion en Go. Implementa un subconjunto de SQL que consulta tablas cargadas desde archivos CSV y mantenidas en memoria.

## Estado

El proyecto se encuentra en su fase inicial. El primer hito sera la carga de CSV, el catalogo de tablas y el manejo de tipos basicos.

## Estructura

```text
cmd/sqlmem/  Punto de entrada del CLI o REPL.
internal/    Implementacion interna del motor.
data/        Archivos CSV de ejemplo.
docs/        Gramatica, decisiones de diseno y declaracion de IA.
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
go run ./cmd/sqlmem
```

## Alcance previsto

1. Carga de CSV como tablas en memoria con catalogo, esquemas y tipos.
2. Lexer, parser y AST para `SELECT ... FROM ... WHERE ...`.
3. Operadores `Scan`, `Filter` y `Project` mediante el modelo Volcano.
4. `ORDER BY`, `LIMIT`, `GROUP BY` y agregados.
5. `INNER JOIN` con nested-loop y hash join.
