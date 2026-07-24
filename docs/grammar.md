# Gramatica soportada

## Hito 2

```ebnf
query       = "SELECT" , select_list , "FROM" , identifier , [ "WHERE" , expression ] , EOF ;
select_list = "*" | identifier , { "," , identifier } ;
expression  = and_expression , { "OR" , and_expression } ;
and_expression = comparison , { "AND" , comparison } ;
comparison  = "(" , expression , ")" | operand , comparison_operator , operand ;
operand     = identifier | number | string | boolean | "NULL" ;
comparison_operator = "=" | "<>" | "<" | ">" | "<=" | ">=" ;
```

Las palabras clave no distinguen mayusculas de minusculas. Los textos se escriben entre comillas simples y una comilla simple interna se escapa duplicandola, por ejemplo: `'O''Brien'`.

Ordenamiento, limite, agregaciones y joins se documentaran al implementarse.
