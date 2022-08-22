package sq

import (
	"bytes"
	"errors"

	"github.com/lann/builder"
)

func init() {
	builder.Register(CaseBuilder{}, caseData{})
}

// sqlizerBuffer is a helper that allows to write many SQLizers one by one
// without constant checks for errors that may come from SQLizer.
type sqlizerBuffer struct {
	bytes.Buffer
	args []interface{}
	err  error
}

// WriteSQL converts SQLizer to SQL strings and writes it to buffer.
func (b *sqlizerBuffer) WriteSQL(item SQLizer) {
	if b.err != nil {
		return
	}

	var str string
	var args []interface{}
	str, args, b.err = nestedToSQL(item)

	if b.err != nil {
		return
	}

	b.WriteString(str)
	b.WriteByte(' ')
	b.args = append(b.args, args...)
}

func (b *sqlizerBuffer) ToSQL() (string, []interface{}, error) {
	return b.String(), b.args, b.err
}

// whenPart is a helper structure to describe SQLs "WHEN ... THEN ..." expression.
type whenPart struct {
	when SQLizer
	then SQLizer
}

func newWhenPart(when interface{}, then interface{}) whenPart {
	return whenPart{newPart(when), newPart(then)}
}

// caseData holds all the data required to build a CASE SQL construct.
type caseData struct {
	What      SQLizer
	WhenParts []whenPart
	Else      SQLizer
}

// ToSQL implements SQLizer.
func (d *caseData) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(d.WhenParts) == 0 {
		err = errors.New("case expression must contain at lease one WHEN clause")

		return
	}

	sql := sqlizerBuffer{}

	sql.WriteString("CASE ")
	if d.What != nil {
		sql.WriteSQL(d.What)
	}

	for _, p := range d.WhenParts {
		sql.WriteString("WHEN ")
		sql.WriteSQL(p.when)
		sql.WriteString("THEN ")
		sql.WriteSQL(p.then)
	}

	if d.Else != nil {
		sql.WriteString("ELSE ")
		sql.WriteSQL(d.Else)
	}

	sql.WriteString("END")

	return sql.ToSQL()
}

// CaseBuilder builds SQL CASE construct which could be used as parts of queries.
type CaseBuilder builder.Builder

// ToSQL builds the query into a SQL string and bound args.
func (b CaseBuilder) ToSQL() (string, []interface{}, error) {
	data := builder.GetStruct(b).(caseData)
	return data.ToSQL()
}

// MustSQL builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b CaseBuilder) MustSQL() (string, []interface{}) {
	sql, args, err := b.ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// what sets optional value for CASE construct "CASE [value] ...".
func (b CaseBuilder) what(expr interface{}) CaseBuilder {
	return builder.Set(b, "What", newPart(expr)).(CaseBuilder)
}

// When adds "WHEN ... THEN ..." part to CASE construct.
func (b CaseBuilder) When(when interface{}, then interface{}) CaseBuilder {
	// TODO: performance hint: replace slice of WhenPart with just slice of parts
	// where even indices of the slice belong to "when"s and odd indices belong to "then"s
	return builder.Append(b, "WhenParts", newWhenPart(when, then)).(CaseBuilder)
}

// Else sets optional "ELSE ..." part for CASE construct.
func (b CaseBuilder) Else(expr interface{}) CaseBuilder {
	return builder.Set(b, "Else", newPart(expr)).(CaseBuilder)
}
