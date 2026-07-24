package executor

import (
	"fmt"
	"io"
	"strings"

	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/catalog"
	"github.com/KellyCubas/motor-consultas-sql-en-memoria/internal/query"
)

type Aggregate struct {
	input   Operator
	groups  []int
	specs   []query.Aggregate
	indices []int
	columns []catalog.Column
	rows    []catalog.Row
	loaded  bool
	index   int
}
type aggregateState struct {
	group catalog.Row
	count []int
	sum   []float64
	first []catalog.Value
	last  []catalog.Value
}

func NewAggregate(input Operator, groupNames []string, specs []query.Aggregate) (*Aggregate, error) {
	a := &Aggregate{input: input, specs: specs, groups: make([]int, len(groupNames)), indices: make([]int, len(specs))}
	for i, name := range groupNames {
		index, column, err := findColumn(input.Columns(), name)
		if err != nil {
			return nil, err
		}
		a.groups[i] = index
		a.columns = append(a.columns, column)
	}
	for i, spec := range specs {
		if spec.Star {
			if spec.Function != query.CountToken {
				return nil, fmt.Errorf("%s requiere una columna", spec.Function)
			}
			a.indices[i] = -1
		} else {
			index, _, err := findColumn(input.Columns(), spec.Column)
			if err != nil {
				return nil, err
			}
			a.indices[i] = index
		}
		name := spec.Function.String() + "(" + spec.Column + ")"
		if spec.Star {
			name = "COUNT(*)"
		}
		a.columns = append(a.columns, catalog.Column{Name: name, Type: aggregateType(spec.Function)})
	}
	return a, nil
}

func (a *Aggregate) Next() (catalog.Row, error) {
	if !a.loaded {
		if err := a.load(); err != nil {
			return nil, err
		}
	}
	if a.index >= len(a.rows) {
		return nil, io.EOF
	}
	row := a.rows[a.index]
	a.index++
	return row, nil
}
func (a *Aggregate) Columns() []catalog.Column { return a.columns }
func (a *Aggregate) Close() error              { return a.input.Close() }

func (a *Aggregate) load() error {
	states := map[string]*aggregateState{}
	order := []string{}
	for {
		row, err := a.input.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		key, group := a.group(row)
		state := states[key]
		if state == nil {
			state = &aggregateState{group: group, count: make([]int, len(a.specs)), sum: make([]float64, len(a.specs)), first: make([]catalog.Value, len(a.specs)), last: make([]catalog.Value, len(a.specs))}
			states[key] = state
			order = append(order, key)
		}
		for i, spec := range a.specs {
			value := catalog.Value{}
			if a.indices[i] >= 0 {
				value = row[a.indices[i]]
			}
			if spec.Function == query.CountToken {
				if spec.Star || !value.Null {
					state.count[i]++
				}
				continue
			}
			if value.Null {
				continue
			}
			state.count[i]++
			if state.count[i] == 1 {
				state.first[i] = value
				state.last[i] = value
			} else if spec.Function == query.MinToken {
				if c, _ := compare(value, state.first[i]); c < 0 {
					state.first[i] = value
				}
			} else if spec.Function == query.MaxToken {
				if c, _ := compare(value, state.last[i]); c > 0 {
					state.last[i] = value
				}
			}
			if spec.Function == query.SumToken || spec.Function == query.AvgToken {
				if !isNumber(value.Type) {
					return fmt.Errorf("%s requiere una columna numerica", spec.Function)
				}
				state.sum[i] += asFloat(value)
			}
		}
	}
	if len(order) == 0 && len(a.groups) == 0 {
		states[""] = &aggregateState{count: make([]int, len(a.specs)), sum: make([]float64, len(a.specs)), first: make([]catalog.Value, len(a.specs)), last: make([]catalog.Value, len(a.specs))}
		order = []string{""}
	}
	for _, key := range order {
		state := states[key]
		row := append(catalog.Row{}, state.group...)
		for i, spec := range a.specs {
			row = append(row, a.result(state, i, spec))
		}
		a.rows = append(a.rows, row)
	}
	a.loaded = true
	return nil
}
func (a *Aggregate) group(row catalog.Row) (string, catalog.Row) {
	values := make(catalog.Row, len(a.groups))
	parts := make([]string, len(a.groups))
	for i, index := range a.groups {
		values[i] = row[index]
		parts[i] = fmt.Sprintf("%d:%v:%t", values[i].Type, values[i].Data, values[i].Null)
	}
	return strings.Join(parts, "|"), values
}
func (a *Aggregate) result(s *aggregateState, i int, spec query.Aggregate) catalog.Value {
	if spec.Function == query.CountToken {
		return catalog.Value{Type: catalog.Integer, Data: int64(s.count[i])}
	}
	if s.count[i] == 0 {
		return catalog.Value{Type: aggregateType(spec.Function), Null: true}
	}
	switch spec.Function {
	case query.SumToken:
		return catalog.Value{Type: catalog.Decimal, Data: s.sum[i]}
	case query.AvgToken:
		return catalog.Value{Type: catalog.Decimal, Data: s.sum[i] / float64(s.count[i])}
	case query.MinToken:
		return s.first[i]
	default:
		return s.last[i]
	}
}
func aggregateType(kind query.TokenKind) catalog.DataType {
	if kind == query.CountToken {
		return catalog.Integer
	}
	return catalog.Decimal
}
