package sq

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/lann/builder"
)

type updateData struct {
	PlaceholderFormat PlaceholderFormat
	Prefixes          []SQLizer
	Table             string
	SetClauses        []setClause
	WhereParts        []SQLizer
	OrderBys          []string
	Limit             string
	Offset            string
	Suffixes          []SQLizer
}

type setClause struct {
	column string
	value  interface{}
}

func (d *updateData) ToSQL() (sqlStr string, args []interface{}, err error) {
	if len(d.Table) == 0 {
		err = fmt.Errorf("update statements must specify a table")
		return
	}
	if len(d.SetClauses) == 0 {
		err = fmt.Errorf("update statements must have at least one Set clause")
		return
	}

	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, err = appendToSQL(d.Prefixes, sql, " ", args)
		if err != nil {
			return
		}

		sql.WriteString(" ")
	}

	sql.WriteString("UPDATE ")
	sql.WriteString(d.Table)

	sql.WriteString(" SET ")
	setSQLs := make([]string, len(d.SetClauses))
	for i, setClause := range d.SetClauses {
		var valSQL string
		if vs, ok := setClause.value.(SQLizer); ok {
			vsql, vargs, err := vs.ToSQL()
			if err != nil {
				return "", nil, err
			}
			if _, ok := vs.(SelectBuilder); ok {
				valSQL = fmt.Sprintf("(%s)", vsql)
			} else {
				valSQL = vsql
			}
			args = append(args, vargs...)
		} else {
			valSQL = "?"
			args = append(args, setClause.value)
		}
		setSQLs[i] = fmt.Sprintf("%s = %s", setClause.column, valSQL)
	}
	sql.WriteString(strings.Join(setSQLs, ", "))

	if len(d.WhereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSQL(d.WhereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(d.OrderBys, ", "))
	}

	if len(d.Limit) > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(d.Limit)
	}

	if len(d.Offset) > 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(d.Offset)
	}

	if len(d.Suffixes) > 0 {
		sql.WriteString(" ")
		args, err = appendToSQL(d.Suffixes, sql, " ", args)
		if err != nil {
			return
		}
	}

	sqlStr, err = d.PlaceholderFormat.ReplacePlaceholders(sql.String())
	return
}

// Builder

// UpdateBuilder builds SQL UPDATE statements.
type UpdateBuilder builder.Builder

func init() {
	builder.Register(UpdateBuilder{}, updateData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the query.
func (b UpdateBuilder) PlaceholderFormat(f PlaceholderFormat) UpdateBuilder {
	return builder.Set(b, "PlaceholderFormat", f).(UpdateBuilder)
}

// SQL methods

// ToSQL builds the query into a SQL string and bound args.
func (b UpdateBuilder) ToSQL() (string, []interface{}, error) {
	data := builder.GetStruct(b).(updateData)
	return data.ToSQL()
}

// MustSQL builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b UpdateBuilder) MustSQL() (string, []interface{}) {
	sql, args, err := b.ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// Prefix adds an expression to the beginning of the query.
func (b UpdateBuilder) Prefix(sql string, args ...interface{}) UpdateBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query.
func (b UpdateBuilder) PrefixExpr(expr SQLizer) UpdateBuilder {
	return builder.Append(b, "Prefixes", expr).(UpdateBuilder)
}

// Table sets the table to be updated.
func (b UpdateBuilder) Table(table string) UpdateBuilder {
	return builder.Set(b, "Table", table).(UpdateBuilder)
}

// Set adds SET clauses to the query.
func (b UpdateBuilder) Set(column string, value interface{}) UpdateBuilder {
	return builder.Append(b, "SetClauses", setClause{column: column, value: value}).(UpdateBuilder)
}

// SetMap is a convenience method which calls .Set for each key/value pair in clauses.
func (b UpdateBuilder) SetMap(clauses map[string]interface{}) UpdateBuilder {
	keys := make([]string, len(clauses))
	i := 0
	for key := range clauses {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		val, _ := clauses[key]
		b = b.Set(key, val)
	}
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b UpdateBuilder) Where(pred interface{}, args ...interface{}) UpdateBuilder {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(UpdateBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b UpdateBuilder) OrderBy(orderBys ...string) UpdateBuilder {
	return builder.Extend(b, "OrderBys", orderBys).(UpdateBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b UpdateBuilder) Limit(limit uint64) UpdateBuilder {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(UpdateBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b UpdateBuilder) Offset(offset uint64) UpdateBuilder {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(UpdateBuilder)
}

// Suffix adds an expression to the end of the query.
func (b UpdateBuilder) Suffix(sql string, args ...interface{}) UpdateBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// SuffixExpr adds an expression to the end of the query.
func (b UpdateBuilder) SuffixExpr(expr SQLizer) UpdateBuilder {
	return builder.Append(b, "Suffixes", expr).(UpdateBuilder)
}
