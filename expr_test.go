package sq

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatExpr(t *testing.T) {
	b := ConcatExpr("COALESCE(name,", Expr("CONCAT(?,' ',?)", "f", "l"), ")")
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "COALESCE(name,CONCAT(?,' ',?))"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"f", "l"}
	assert.Equal(t, expectedArgs, args)
}

func TestConcatExprBadType(t *testing.T) {
	b := ConcatExpr("prefix", 123, "suffix")
	_, _, err := b.ToSQL()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "123 is not")
}

func TestEqToSQL(t *testing.T) {
	b := Eq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id = ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestEqEmptyToSQL(t *testing.T) {
	sql, args, err := Eq{}.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "(1=1)"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, args)
}

func TestEqInToSQL(t *testing.T) {
	b := Eq{"id": []int{1, 2, 3}}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id IN (?,?,?)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1, 2, 3}
	assert.Equal(t, expectedArgs, args)
}

func TestNotEqToSQL(t *testing.T) {
	b := NotEq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id <> ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestEqNotInToSQL(t *testing.T) {
	b := NotEq{"id": []int{1, 2, 3}}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id NOT IN (?,?,?)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1, 2, 3}
	assert.Equal(t, expectedArgs, args)
}

func TestEqInEmptyToSQL(t *testing.T) {
	b := Eq{"id": []int{}}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "(1=0)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestNotEqInEmptyToSQL(t *testing.T) {
	b := NotEq{"id": []int{}}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "(1=1)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestEqBytesToSQL(t *testing.T) {
	b := Eq{"id": []byte("test")}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id = ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{[]byte("test")}
	assert.Equal(t, expectedArgs, args)
}

func TestLtToSQL(t *testing.T) {
	b := Lt{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id < ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestLtOrEqToSQL(t *testing.T) {
	b := LtOrEq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id <= ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestGtToSQL(t *testing.T) {
	b := Gt{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id > ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestGtOrEqToSQL(t *testing.T) {
	b := GtOrEq{"id": 1}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "id >= ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestExprNilToSQL(t *testing.T) {
	var b SQLizer
	b = NotEq{"name": nil}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)

	expectedSQL := "name IS NOT NULL"
	assert.Equal(t, expectedSQL, sql)

	b = Eq{"name": nil}
	sql, args, err = b.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)

	expectedSQL = "name IS NULL"
	assert.Equal(t, expectedSQL, sql)
}

func TestNullTypeString(t *testing.T) {
	var b SQLizer
	var name sql.NullString

	b = Eq{"name": name}
	sql, args, err := b.ToSQL()

	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "name IS NULL", sql)

	name.Scan("Name")
	b = Eq{"name": name}
	sql, args, err = b.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"Name"}, args)
	assert.Equal(t, "name = ?", sql)
}

func TestNullTypeInt64(t *testing.T) {
	var userID sql.NullInt64
	userID.Scan(nil)
	b := Eq{"user_id": userID}
	sql, args, err := b.ToSQL()

	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "user_id IS NULL", sql)

	userID.Scan(int64(10))
	b = Eq{"user_id": userID}
	sql, args, err = b.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, []interface{}{int64(10)}, args)
	assert.Equal(t, "user_id = ?", sql)
}

func TestNilPointer(t *testing.T) {
	var name *string = nil
	eq := Eq{"name": name}
	sql, args, err := eq.ToSQL()

	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "name IS NULL", sql)

	neq := NotEq{"name": name}
	sql, args, err = neq.ToSQL()

	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "name IS NOT NULL", sql)

	var ids *[]int = nil
	eq = Eq{"id": ids}
	sql, args, err = eq.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "id IS NULL", sql)

	neq = NotEq{"id": ids}
	sql, args, err = neq.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "id IS NOT NULL", sql)

	var ida *[3]int = nil
	eq = Eq{"id": ida}
	sql, args, err = eq.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "id IS NULL", sql)

	neq = NotEq{"id": ida}
	sql, args, err = neq.ToSQL()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "id IS NOT NULL", sql)

}

func TestNotNilPointer(t *testing.T) {
	c := "Name"
	name := &c
	eq := Eq{"name": name}
	sql, args, err := eq.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"Name"}, args)
	assert.Equal(t, "name = ?", sql)

	neq := NotEq{"name": name}
	sql, args, err = neq.ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"Name"}, args)
	assert.Equal(t, "name <> ?", sql)

	s := []int{1, 2, 3}
	ids := &s
	eq = Eq{"id": ids}
	sql, args, err = eq.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{1, 2, 3}, args)
	assert.Equal(t, "id IN (?,?,?)", sql)

	neq = NotEq{"id": ids}
	sql, args, err = neq.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{1, 2, 3}, args)
	assert.Equal(t, "id NOT IN (?,?,?)", sql)

	a := [3]int{1, 2, 3}
	ida := &a
	eq = Eq{"id": ida}
	sql, args, err = eq.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{1, 2, 3}, args)
	assert.Equal(t, "id IN (?,?,?)", sql)

	neq = NotEq{"id": ida}
	sql, args, err = neq.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{1, 2, 3}, args)
	assert.Equal(t, "id NOT IN (?,?,?)", sql)
}

func TestEmptyAndToSQL(t *testing.T) {
	sql, args, err := And{}.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "(1=1)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestEmptyOrToSQL(t *testing.T) {
	sql, args, err := Or{}.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "(1=0)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestLikeToSQL(t *testing.T) {
	b := Like{"name": "%irrel"}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "name LIKE ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"%irrel"}
	assert.Equal(t, expectedArgs, args)
}

func TestNotLikeToSQL(t *testing.T) {
	b := NotLike{"name": "%irrel"}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "name NOT LIKE ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"%irrel"}
	assert.Equal(t, expectedArgs, args)
}

func TestILikeToSQL(t *testing.T) {
	b := ILike{"name": "sq%"}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "name ILIKE ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"sq%"}
	assert.Equal(t, expectedArgs, args)
}

func TestNotILikeToSQL(t *testing.T) {
	b := NotILike{"name": "sq%"}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "name NOT ILIKE ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{"sq%"}
	assert.Equal(t, expectedArgs, args)
}

func TestSQLEqOrder(t *testing.T) {
	b := Eq{"a": 1, "b": 2, "c": 3}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "a = ? AND b = ? AND c = ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1, 2, 3}
	assert.Equal(t, expectedArgs, args)
}

func TestSQLLtOrder(t *testing.T) {
	b := Lt{"a": 1, "b": 2, "c": 3}
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "a < ? AND b < ? AND c < ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1, 2, 3}
	assert.Equal(t, expectedArgs, args)
}

func TestExprEscaped(t *testing.T) {
	b := Expr("count(??)", Expr("x"))
	sql, args, err := b.ToSQL()
	assert.NoError(t, err)

	expectedSQL := "count(??)"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{Expr("x")}
	assert.Equal(t, expectedArgs, args)
}

func TestExprRecursion(t *testing.T) {
	{
		b := Expr("count(?)", Expr("nullif(a,?)", "b"))
		sql, args, err := b.ToSQL()
		assert.NoError(t, err)

		expectedSQL := "count(nullif(a,?))"
		assert.Equal(t, expectedSQL, sql)

		expectedArgs := []interface{}{"b"}
		assert.Equal(t, expectedArgs, args)
	}
	{
		b := Expr("extract(? from ?)", Expr("epoch"), "2001-02-03")
		sql, args, err := b.ToSQL()
		assert.NoError(t, err)

		expectedSQL := "extract(epoch from ?)"
		assert.Equal(t, expectedSQL, sql)

		expectedArgs := []interface{}{"2001-02-03"}
		assert.Equal(t, expectedArgs, args)
	}
	{
		b := Expr("JOIN t1 ON ?", And{Eq{"id": 1}, Expr("NOT c1"), Expr("? @@ ?", "x", "y")})
		sql, args, err := b.ToSQL()
		assert.NoError(t, err)

		expectedSQL := "JOIN t1 ON (id = ? AND NOT c1 AND ? @@ ?)"
		assert.Equal(t, expectedSQL, sql)

		expectedArgs := []interface{}{1, "x", "y"}
		assert.Equal(t, expectedArgs, args)
	}
}

func ExampleEq() {
	Select("id", "created", "first_name").From("users").Where(Eq{
		"company": 20,
	})
}
