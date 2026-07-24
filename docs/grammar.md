# Gramatica soportada

Este documento se completara de forma incremental conforme se implementen los hitos. La gramatica inicial prevista es:

```ebnf
query       = "SELECT" , select_list , "FROM" , identifier , [ "WHERE" , expression ] ;
select_list = "*" | identifier , { "," , identifier } ;
```

Las expresiones, ordenamiento, limite, agregaciones y joins se documentaran al implementarse.
