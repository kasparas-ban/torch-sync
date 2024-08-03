package optional

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

type NullUint64 struct {
	Val   uint64
	Valid bool
}

func NewNullUint64(val interface{}) NullUint64 {
	ni := NullUint64{}
	ni.Set(val)
	return ni
}

func (ni *NullUint64) Scan(value interface{}) error {
	data, ok := value.(int64)
	if ok {
		ni.Val = uint64(data)
		ni.Valid = true
		return nil
	}

	strData, ok := value.([]uint8)
	if !ok {
		return nil
	}
	val, err := strconv.ParseUint(string(strData), 10, 64)
	if err == nil {
		ni.Val = val
		ni.Valid = true
	}

	return nil
}

func (ni NullUint64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Val, nil
}

func (ni *NullUint64) Set(val interface{}) {
	ni.Val, ni.Valid = val.(uint64)
}

func (ni NullUint64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte(`null`), nil
	}

	return json.Marshal(ni.Val)
}

func (ni *NullUint64) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.Valid = false
		return nil
	}

	val, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		ni.Valid = false
		return err
	}

	ni.Val = val
	ni.Valid = true

	return nil
}

func (ni NullUint64) String() string {
	if !ni.Valid {
		return `<nil>`
	}

	return strconv.FormatUint(ni.Val, 10)
}
