package sq

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DBStub struct {
	err error

	LastPrepareSQL string
	PrepareCount   int

	LastExecSQL  string
	LastExecArgs []interface{}

	LastQuerySQL  string
	LastQueryArgs []interface{}

	LastQueryRowSQL  string
	LastQueryRowArgs []interface{}
}

func (s *DBStub) Prepare(query string) (*sql.Stmt, error) {
	s.LastPrepareSQL = query
	s.PrepareCount++
	return nil, nil
}

var sqlizer = Select("test")

var testDebugUpdateSQL = Update("table").SetMap(Eq{"x": 1, "y": "val"})
var expectedDebugUpateSQL = "UPDATE table SET x = '1', y = 'val'"

func TestDebugSQLizerUpdateColon(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Colon)
	assert.Equal(t, expectedDebugUpateSQL, DebugSQLizer(testDebugUpdateSQL))
}

func TestDebugSQLizerUpdateAtp(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(AtP)
	assert.Equal(t, expectedDebugUpateSQL, DebugSQLizer(testDebugUpdateSQL))
}

func TestDebugSQLizerUpdateDollar(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Dollar)
	assert.Equal(t, expectedDebugUpateSQL, DebugSQLizer(testDebugUpdateSQL))
}

func TestDebugSQLizerUpdateQuestion(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Question)
	assert.Equal(t, expectedDebugUpateSQL, DebugSQLizer(testDebugUpdateSQL))
}

var testDebugDeleteSQL = Delete("table").Where(And{
	Eq{"column": "val"},
	Eq{"other": 1},
})
var expectedDebugDeleteSQL = "DELETE FROM table WHERE (column = 'val' AND other = '1')"

func TestDebugSQLizerDeleteColon(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Colon)
	assert.Equal(t, expectedDebugDeleteSQL, DebugSQLizer(testDebugDeleteSQL))
}

func TestDebugSQLizerDeleteAtp(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(AtP)
	assert.Equal(t, expectedDebugDeleteSQL, DebugSQLizer(testDebugDeleteSQL))
}

func TestDebugSQLizerDeleteDollar(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Dollar)
	assert.Equal(t, expectedDebugDeleteSQL, DebugSQLizer(testDebugDeleteSQL))
}

func TestDebugSQLizerDeleteQuestion(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Question)
	assert.Equal(t, expectedDebugDeleteSQL, DebugSQLizer(testDebugDeleteSQL))
}

var testDebugInsertSQL = Insert("table").Values(1, "test")
var expectedDebugInsertSQL = "INSERT INTO table VALUES ('1','test')"

func TestDebugSQLizerInsertColon(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Colon)
	assert.Equal(t, expectedDebugInsertSQL, DebugSQLizer(testDebugInsertSQL))
}

func TestDebugSQLizerInsertAtp(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(AtP)
	assert.Equal(t, expectedDebugInsertSQL, DebugSQLizer(testDebugInsertSQL))
}

func TestDebugSQLizerInsertDollar(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Dollar)
	assert.Equal(t, expectedDebugInsertSQL, DebugSQLizer(testDebugInsertSQL))
}

func TestDebugSQLizerInsertQuestion(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Question)
	assert.Equal(t, expectedDebugInsertSQL, DebugSQLizer(testDebugInsertSQL))
}

var testDebugSelectSQL = Select("*").From("table").Where(And{
	Eq{"column": "val"},
	Eq{"other": 1},
})
var expectedDebugSelectSQL = "SELECT * FROM table WHERE (column = 'val' AND other = '1')"

func TestDebugSQLizerSelectColon(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Colon)
	assert.Equal(t, expectedDebugSelectSQL, DebugSQLizer(testDebugSelectSQL))
}

func TestDebugSQLizerSelectAtp(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(AtP)
	assert.Equal(t, expectedDebugSelectSQL, DebugSQLizer(testDebugSelectSQL))
}

func TestDebugSQLizerSelectDollar(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Dollar)
	assert.Equal(t, expectedDebugSelectSQL, DebugSQLizer(testDebugSelectSQL))
}

func TestDebugSQLizerSelectQuestion(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Question)
	assert.Equal(t, expectedDebugSelectSQL, DebugSQLizer(testDebugSelectSQL))
}

func TestDebugSQLizer(t *testing.T) {
	sqlizer := Expr("x = ? AND y = ? AND z = '??'", 1, "text")
	expectedDebug := "x = '1' AND y = 'text' AND z = '?'"
	assert.Equal(t, expectedDebug, DebugSQLizer(sqlizer))
}

func TestDebugSQLizerErrors(t *testing.T) {
	errorMsg := DebugSQLizer(Expr("x = ?", 1, 2)) // Not enough placeholders
	assert.True(t, strings.HasPrefix(errorMsg, "[DebugSQLizer error: "))

	errorMsg = DebugSQLizer(Expr("x = ? AND y = ?", 1)) // Too many placeholders
	assert.True(t, strings.HasPrefix(errorMsg, "[DebugSQLizer error: "))

	errorMsg = DebugSQLizer(Lt{"x": nil}) // Cannot use nil values with Lt
	assert.True(t, strings.HasPrefix(errorMsg, "[ToSQL error: "))
}
