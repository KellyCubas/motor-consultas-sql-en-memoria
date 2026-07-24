package query

import "testing"

func TestLexRecognizesKeywordsAndOperators(t *testing.T) {
	tokens, err := Lex("select * from empleados where edad >= 18 and activo <> false")
	if err != nil {
		t.Fatalf("Lex devolvio error: %v", err)
	}

	want := []TokenKind{SelectToken, StarToken, FromToken, IdentifierToken, WhereToken, IdentifierToken, GreaterEqualToken, NumberToken, AndToken, IdentifierToken, NotEqualToken, BooleanToken, EOF}
	if len(tokens) != len(want) {
		t.Fatalf("tokens = %d; se esperaban %d", len(tokens), len(want))
	}
	for index, kind := range want {
		if tokens[index].Kind != kind {
			t.Errorf("token %d = %s; se esperaba %s", index, tokens[index].Kind, kind)
		}
	}
}
