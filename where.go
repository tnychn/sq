package sq

import (
	"fmt"
)

type wherePart part

func newWherePart(pred interface{}, args ...interface{}) SQLizer {
	return &wherePart{pred: pred, args: args}
}

func (p wherePart) ToSQL() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case rawSQLizer:
		return pred.toSQLRaw()
	case SQLizer:
		return pred.ToSQL()
	case map[string]interface{}:
		return Eq(pred).ToSQL()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string-keyed map or string, not %T", pred)
	}
	return
}
