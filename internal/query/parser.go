package query

import "fmt"

// Parse analiza una consulta SQL y devuelve su AST.
func Parse(input string) (*Query, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, err
	}
	return ParseTokens(tokens)
}

// ParseTokens analiza una lista de tokens terminada en EOF.
func ParseTokens(tokens []Token) (*Query, error) {
	parser := parser{tokens: tokens}
	query, err := parser.parseQuery()
	if err != nil {
		return nil, err
	}
	if token := parser.current(); token.Kind != EOF {
		return nil, parser.unexpected(token, "fin de consulta")
	}
	return query, nil
}

type parser struct {
	tokens []Token
	index  int
}

func (p *parser) parseQuery() (*Query, error) {
	if _, err := p.expect(SelectToken); err != nil {
		return nil, err
	}

	query := &Query{}
	if p.match(StarToken) {
		query.SelectAll = true
	} else {
		columns, err := p.parseColumns()
		if err != nil {
			return nil, err
		}
		query.Columns = columns
	}

	if _, err := p.expect(FromToken); err != nil {
		return nil, err
	}
	table, err := p.expect(IdentifierToken)
	if err != nil {
		return nil, err
	}
	query.Table = table.Lexeme

	if p.match(WhereToken) {
		where, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		query.Where = where
	}

	return query, nil
}

func (p *parser) parseColumns() ([]string, error) {
	columns := make([]string, 0, 1)
	for {
		column, err := p.expect(IdentifierToken)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column.Lexeme)
		if !p.match(CommaToken) {
			return columns, nil
		}
	}
}

func (p *parser) parseOr() (Expression, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.match(OrToken) {
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = BinaryExpression{Left: left, Operator: OrToken, Right: right}
	}
	return left, nil
}

func (p *parser) parseAnd() (Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.match(AndToken) {
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = BinaryExpression{Left: left, Operator: AndToken, Right: right}
	}
	return left, nil
}

func (p *parser) parseComparison() (Expression, error) {
	if p.match(LeftParenToken) {
		expression, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(RightParenToken); err != nil {
			return nil, err
		}
		return expression, nil
	}

	left, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	operator := p.current()
	if !isComparison(operator.Kind) {
		return nil, p.unexpected(operator, "operador de comparacion")
	}
	p.index++
	right, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	return BinaryExpression{Left: left, Operator: operator.Kind, Right: right}, nil
}

func (p *parser) parseOperand() (Expression, error) {
	token := p.current()
	switch token.Kind {
	case IdentifierToken:
		p.index++
		return Identifier{Name: token.Lexeme}, nil
	case NumberToken, StringToken, BooleanToken, NullToken:
		p.index++
		return Literal{Value: token.Lexeme, Kind: token.Kind}, nil
	default:
		return nil, p.unexpected(token, "identificador o literal")
	}
}

func (p *parser) current() Token {
	if p.index >= len(p.tokens) {
		return Token{Kind: EOF}
	}
	return p.tokens[p.index]
}

func (p *parser) match(kind TokenKind) bool {
	if p.current().Kind != kind {
		return false
	}
	p.index++
	return true
}

func (p *parser) expect(kind TokenKind) (Token, error) {
	token := p.current()
	if token.Kind != kind {
		return Token{}, p.unexpected(token, kind.String())
	}
	p.index++
	return token, nil
}

func (p *parser) unexpected(token Token, expected string) error {
	return fmt.Errorf("error de sintaxis en la posicion %d: se esperaba %s y se encontro %s", token.Position, expected, token.Kind)
}

func isComparison(kind TokenKind) bool {
	switch kind {
	case EqualToken, NotEqualToken, LessToken, GreaterToken, LessEqualToken, GreaterEqualToken:
		return true
	default:
		return false
	}
}
