# Bitacora de decisiones

## Plantilla de entrada

- Hito y fecha:
- Decision tomada:
- Alternativas evaluadas:
- Justificacion tecnica:

## Hito inicial - estructura del repositorio

- Hito y fecha: Preparacion inicial, 2026-07-24.
- Decision tomada: Separar el punto de entrada en `cmd/sqlmem`, la implementacion en `internal`, los datos de ejemplo en `data` y la documentacion en `docs`.
- Alternativas evaluadas: Mantener todos los archivos en la raiz del proyecto.
- Justificacion tecnica: La separacion facilita el crecimiento del motor sin mezclar el ejecutable, la logica interna, los datos y los documentos requeridos.

## Hito 1 - carga de CSV y tipos

- Hito y fecha: Hito 1, 2026-07-24.
- Decision tomada: Representar una fila como una lista de valores tipados y conservar explicitamente si un valor es `NULL`.
- Alternativas evaluadas: Guardar todas las celdas como texto y convertirlas durante cada consulta.
- Justificacion tecnica: Convertir los valores al cargar el CSV valida los datos una sola vez y permite que los futuros operadores comparen valores segun su tipo real.

## Hito 2 - analisis de consultas

- Hito y fecha: Hito 2, 2026-07-24.
- Decision tomada: Separar lexer, parser y AST en el paquete `internal/query`.
- Alternativas evaluadas: Interpretar la consulta directamente mientras se recorre la cadena.
- Justificacion tecnica: El AST conserva la estructura y precedencia de la consulta, por lo que los operadores de ejecucion podran construirse despues sin mezclar analisis sintactico y acceso a datos.

## Hito 3 - operadores de ejecucion

- Hito y fecha: Hito 3, 2026-07-24.
- Decision tomada: Usar una interfaz `Operator` con los metodos `Next`, `Columns` y `Close`.
- Alternativas evaluadas: Ejecutar cada consulta construyendo slices intermedios en una unica funcion.
- Justificacion tecnica: Los operadores encadenados conservan la evaluacion perezosa. `Filter` solicita filas a `Scan` y `Project` solicita filas a `Filter`, sin conocer como se implementa el operador inferior.

## Hito 4 - ordenamiento y limite

- Hito y fecha: Hito 4, 2026-07-24.
- Decision tomada: Materializar las filas solo en el operador `Order` y aplicar `Limit` despues de ordenar.
- Alternativas evaluadas: Ordenar directamente dentro de `Scan` o limitar antes del ordenamiento.
- Justificacion tecnica: Ordenar requiere conocer todas las filas candidatas; aplicar primero `LIMIT` produciria resultados incorrectos cuando la consulta incluye `ORDER BY`.
- Comando de prueba: `go run ./cmd/sqlmem consultar empleados data/empleados.csv "SELECT nombre, salario FROM empleados ORDER BY salario DESC LIMIT 2"`.
