package optional

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

type NullUint struct {
	Val     uint
	IsValid bool
}

func NewNullUint(val interface{}) NullUint {
	ni := NullUint{}
	ni.Set(val)
	return ni
}

func (ni *NullUint) Scan(value interface{}) error {
	data, ok := value.(int64)
	if ok {
		ni.Val = uint(data)
		ni.IsValid = true
	}

	strData, ok := value.([]uint8)
	if !ok {
		return nil
	}
	val, err := strconv.ParseUint(string(strData), 10, 64)
	if err == nil {
		ni.Val = uint(val)
		ni.IsValid = true
	}

	return nil
}

func (ni NullUint) Value() (driver.Value, error) {
	if !ni.IsValid {
		return nil, nil
	}
	return int64(ni.Val), nil
}

func (ni *NullUint) Set(val interface{}) {
	ni.Val, ni.IsValid = val.(uint)
}

func (ni NullUint) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
		return []byte(`null`), nil
	}

	return json.Marshal(ni.Val)
}

func (ni *NullUint) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.IsValid = false
		return nil
	}

	val, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		ni.IsValid = false
		return err
	}

	ni.Val = uint(val)
	ni.IsValid = true

	return nil
}

func (ni NullUint) String() string {
	if !ni.IsValid {
		return `<nil>`
	}

	return strconv.FormatUint(uint64(ni.Val), 10)
}
