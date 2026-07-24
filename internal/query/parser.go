package query

import (
	"fmt"
	"strconv"
)

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
		columns, aggregates, err := p.parseColumns()
		if err != nil {
			return nil, err
		}
		query.Columns = columns
		query.Aggregates = aggregates
	}

	if _, err := p.expect(FromToken); err != nil {
		return nil, err
	}
	table, err := p.expect(IdentifierToken)
	if err != nil {
		return nil, err
	}
	query.Table = table.Lexeme
	if p.match(InnerToken) {
		if _, err := p.expect(JoinToken); err != nil {
			return nil, err
		}
		right, err := p.expect(IdentifierToken)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(OnToken); err != nil {
			return nil, err
		}
		condition, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		query.Join = &Join{Table: right.Lexeme, Condition: condition}
	}

	if p.match(WhereToken) {
		where, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		query.Where = where
	}
	if p.match(GroupToken) {
		if _, err := p.expect(ByToken); err != nil {
			return nil, err
		}
		columns, _, err := p.parseColumns()
		if err != nil {
			return nil, err
		}
		query.GroupBy = columns
	}
	if p.match(OrderToken) {
		if _, err := p.expect(ByToken); err != nil {
			return nil, err
		}
		terms, err := p.parseOrderTerms()
		if err != nil {
			return nil, err
		}
		query.OrderBy = terms
	}
	if p.match(LimitToken) {
		token, err := p.expect(NumberToken)
		if err != nil {
			return nil, err
		}
		limit, err := strconv.Atoi(token.Lexeme)
		if err != nil || limit < 0 || token.Lexeme != strconv.Itoa(limit) {
			return nil, fmt.Errorf("LIMIT invalido en la posicion %d", token.Position)
		}
		query.Limit = &limit
	}

	return query, nil
}

func (p *parser) parseOrderTerms() ([]OrderTerm, error) {
	terms := make([]OrderTerm, 0, 1)
	for {
		column, err := p.expect(IdentifierToken)
		if err != nil {
			return nil, err
		}
		term := OrderTerm{Column: column.Lexeme}
		if p.match(DescToken) {
			term.Descending = true
		} else {
			p.match(AscToken)
		}
		terms = append(terms, term)
		if !p.match(CommaToken) {
			return terms, nil
		}
	}
}

func (p *parser) parseColumns() ([]string, []Aggregate, error) {
	columns := make([]string, 0, 1)
	aggregates := make([]Aggregate, 0)
	for {
		if isAggregate(p.current().Kind) {
			function := p.current().Kind
			p.index++
			if _, err := p.expect(LeftParenToken); err != nil {
				return nil, nil, err
			}
			aggregate := Aggregate{Function: function}
			if p.match(StarToken) {
				aggregate.Star = true
			} else {
				column, err := p.expect(IdentifierToken)
				if err != nil {
					return nil, nil, err
				}
				aggregate.Column = column.Lexeme
			}
			if _, err := p.expect(RightParenToken); err != nil {
				return nil, nil, err
			}
			aggregates = append(aggregates, aggregate)
		} else {
			column, err := p.expect(IdentifierToken)
			if err != nil {
				return nil, nil, err
			}
			columns = append(columns, column.Lexeme)
		}
		if !p.match(CommaToken) {
			return columns, aggregates, nil
		}
	}
}

func isAggregate(kind TokenKind) bool {
	return kind == CountToken || kind == SumToken || kind == AvgToken || kind == MinToken || kind == MaxToken
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
