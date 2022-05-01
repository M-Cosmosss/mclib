package dbutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

type NullTime time.Time

func (t *NullTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	tmp, ok := value.(time.Time)
	if !ok {
		return errors.Errorf("dbutil.NullTime: Scan wants time.Time but got %T", value)
	}

	*t = NullTime(tmp)
	return nil
}

func (t NullTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return t, nil
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(time.Time(t))
}

var _ json.Marshaler = (*NullTime)(nil)
var _ sql.Scanner = (*NullTime)(nil)
var _ driver.Valuer = (*NullTime)(nil)
