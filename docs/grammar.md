# Gramatica soportada

## Hito 2

```ebnf
query       = "SELECT" , select_list , "FROM" , identifier , [ "WHERE" , expression ] , [ order_by ] , [ limit ] , EOF ;
select_list = "*" | identifier , { "," , identifier } ;
expression  = and_expression , { "OR" , and_expression } ;
and_expression = comparison , { "AND" , comparison } ;
comparison  = "(" , expression , ")" | operand , comparison_operator , operand ;
operand     = identifier | number | string | boolean | "NULL" ;
comparison_operator = "=" | "<>" | "<" | ">" | "<=" | ">=" ;
order_by    = "ORDER" , "BY" , order_term , { "," , order_term } ;
order_term  = identifier , [ "ASC" | "DESC" ] ;
limit       = "LIMIT" , integer ;
```

Las palabras clave no distinguen mayusculas de minusculas. Los textos se escriben entre comillas simples y una comilla simple interna se escapa duplicandola, por ejemplo: `'O''Brien'`.

Ordenamiento, limite, agregaciones y joins se documentaran al implementarse.

## Ejecucion de Hito 4

```bash
go run ./cmd/sqlmem consultar empleados data/empleados.csv "SELECT nombre, salario FROM empleados ORDER BY salario DESC LIMIT 2"
```
