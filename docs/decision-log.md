# Bitacora de decisiones

## Plantilla de entrada

- Hito y fecha:
- Decision tomada:
- Alternativas evaluadas:
- Justificacion tecnica:
- Uso de IA, si aplica:

## Hito inicial - estructura del repositorio

- Hito y fecha: Preparacion inicial, 2026-07-24.
- Decision tomada: Separar el punto de entrada en `cmd/sqlmem`, la implementacion en `internal`, los datos de ejemplo en `data` y la documentacion en `docs`.
- Alternativas evaluadas: Mantener todos los archivos en la raiz del proyecto.
- Justificacion tecnica: La separacion facilita el crecimiento del motor sin mezclar el ejecutable, la logica interna, los datos y los documentos requeridos.
- Uso de IA, si aplica: Pendiente de completar por los integrantes segun el uso real.

## Hito 1 - carga de CSV y tipos

- Hito y fecha: Hito 1, 2026-07-24.
- Decision tomada: Representar una fila como una lista de valores tipados y conservar explicitamente si un valor es `NULL`.
- Alternativas evaluadas: Guardar todas las celdas como texto y convertirlas durante cada consulta.
- Justificacion tecnica: Convertir los valores al cargar el CSV valida los datos una sola vez y permite que los futuros operadores comparen valores segun su tipo real.
- Uso de IA, si aplica: Pendiente de completar por los integrantes segun el uso real.
