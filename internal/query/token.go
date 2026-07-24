package query

import "fmt"

// TokenKind identifica una unidad lexica de SQL.
type TokenKind uint8

const (
	Invalid TokenKind = iota
	EOF
	IdentifierToken
	NumberToken
	StringToken
	BooleanToken
	NullToken
	SelectToken
	FromToken
	WhereToken
	AndToken
	OrToken
	OrderToken
	ByToken
	AscToken
	DescToken
	LimitToken
	GroupToken
	CountToken
	SumToken
	AvgToken
	MinToken
	MaxToken
	InnerToken
	JoinToken
	OnToken
	CommaToken
	StarToken
	LeftParenToken
	RightParenToken
	EqualToken
	NotEqualToken
	LessToken
	GreaterToken
	LessEqualToken
	GreaterEqualToken
)

func (k TokenKind) String() string {
	switch k {
	case EOF:
		return "fin de consulta"
	case IdentifierToken:
		return "identificador"
	case NumberToken:
		return "numero"
	case StringToken:
		return "texto"
	case BooleanToken:
		return "booleano"
	case NullToken:
		return "NULL"
	case SelectToken:
		return "SELECT"
	case FromToken:
		return "FROM"
	case WhereToken:
		return "WHERE"
	case AndToken:
		return "AND"
	case OrToken:
		return "OR"
	case OrderToken:
		return "ORDER"
	case ByToken:
		return "BY"
	case AscToken:
		return "ASC"
	case DescToken:
		return "DESC"
	case LimitToken:
		return "LIMIT"
	case GroupToken:
		return "GROUP"
	case CountToken:
		return "COUNT"
	case SumToken:
		return "SUM"
	case AvgToken:
		return "AVG"
	case MinToken:
		return "MIN"
	case MaxToken:
		return "MAX"
	case InnerToken:
		return "INNER"
	case JoinToken:
		return "JOIN"
	case OnToken:
		return "ON"
	case CommaToken:
		return ","
	case StarToken:
		return "*"
	case LeftParenToken:
		return "("
	case RightParenToken:
		return ")"
	case EqualToken:
		return "="
	case NotEqualToken:
		return "<>"
	case LessToken:
		return "<"
	case GreaterToken:
		return ">"
	case LessEqualToken:
		return "<="
	case GreaterEqualToken:
		return ">="
	default:
		return "invalido"
	}
}

// Token conserva el lexema y su posicion inicial, basada en cero.
type Token struct {
	Kind     TokenKind
	Lexeme   string
	Position int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q) en %d", t.Kind, t.Lexeme, t.Position)
}
