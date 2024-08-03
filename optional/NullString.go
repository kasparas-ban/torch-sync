package optional

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullString struct {
	Val     string
	IsValid bool
}

func NewNullString(val interface{}) NullString {
	ni := NullString{}
	ni.Set(val)
	return ni
}

func (ni *NullString) Scan(value interface{}) error {
	if val, ok := value.(time.Time); ok {
		val, err := val.MarshalText()
		if err == nil {
			ni.Val = string(val)
			ni.IsValid = true
			return nil
		}
	}

	if val, ok := value.([]uint8); ok {
		ni.Val = string(val)
		ni.IsValid = true
		return nil
	}

	if val, ok := value.(string); ok {
		ni.Val = val
		ni.IsValid = true
		return nil
	}

	ni.IsValid = false
	return nil
}

func (ni NullString) Value() (driver.Value, error) {
	if !ni.IsValid {
		return nil, nil
	}
	return ni.Val, nil
}

func (ni *NullString) Set(val interface{}) {
	ni.Val, ni.IsValid = val.(string)
}

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
		return []byte(`null`), nil
	}

	return json.Marshal(ni.Val)
}

func (ni *NullString) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.IsValid = false
		return nil
	}

	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		ni.IsValid = false
		return nil
	}

	ni.Val = val
	ni.IsValid = true

	return nil
}

func (ni NullString) String() string {
	if !ni.IsValid {
		return `<nil>`
	}

	return ni.Val
}
