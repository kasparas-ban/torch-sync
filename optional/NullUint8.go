package optional

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

type NullUint8 struct {
	Val     uint8
	IsValid bool
}

func NewNullUint8(val interface{}) NullUint {
	ni := NullUint{}
	ni.Set(val)
	return ni
}

func (ni *NullUint8) Scan(value interface{}) error {
	int64Data, ok := value.(int64)
	if ok {
		ni.Val = uint8(int64Data)
		ni.IsValid = true
		return nil
	}

	uint8Data, ok := value.(uint8)
	if ok {
		ni.Val = uint8Data
		ni.IsValid = true
		return nil
	}

	strData, ok := value.([]uint8)
	if ok {
		val, err := strconv.ParseUint(string(strData), 10, 64)
		if err == nil {
			ni.Val = uint8(val)
			ni.IsValid = true
		}

		return nil
	}

	return nil
}

func (ni NullUint8) Value() (driver.Value, error) {
	if !ni.IsValid {
		return nil, nil
	}
	return int64(ni.Val), nil
}

func (ni *NullUint8) Set(val interface{}) {
	ni.Val, ni.IsValid = val.(uint8)
}

func (ni NullUint8) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
		return []byte(`null`), nil
	}

	return json.Marshal(ni.Val)
}

func (ni *NullUint8) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.IsValid = false
		return nil
	}

	val, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		ni.IsValid = false
		return err
	}

	ni.Val = uint8(val)
	ni.IsValid = true

	return nil
}

func (ni NullUint8) String() string {
	if !ni.IsValid {
		return `<nil>`
	}

	return strconv.FormatUint(uint64(ni.Val), 10)
}
