package sq

import (
	"fmt"
	"io"
)

type part struct {
	pred interface{}
	args []interface{}
}

func newPart(pred interface{}, args ...interface{}) SQLizer {
	return &part{pred, args}
}

func (p part) ToSQL() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case SQLizer:
		sql, args, err = nestedToSQL(pred)
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string or SQLizer, not %T", pred)
	}
	return
}

func nestedToSQL(s SQLizer) (string, []interface{}, error) {
	if raw, ok := s.(rawSQLizer); ok {
		return raw.toSQLRaw()
	} else {
		return s.ToSQL()
	}
}

func appendToSQL(parts []SQLizer, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		partSQL, partArgs, err := nestedToSQL(p)
		if err != nil {
			return nil, err
		} else if len(partSQL) == 0 {
			continue
		}

		if i > 0 {
			_, err := io.WriteString(w, sep)
			if err != nil {
				return nil, err
			}
		}

		_, err = io.WriteString(w, partSQL)
		if err != nil {
			return nil, err
		}
		args = append(args, partArgs...)
	}
	return args, nil
}
