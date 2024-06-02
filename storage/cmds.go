package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"torch/torch-sync/types"
)

type Op struct {
	Op     string          `json:"op"` // INSERT, UPDATE or DELETE
	ItemID string          `json:"itemID"`
	Diffs  Diffs           `json:"diffs"`
	Cl     types.NullInt64 `json:"cl,omitempty"`
}

type Diffs struct {
	Title       *FieldVal[string] `json:"title,omitempty"`
	ItemType    *FieldVal[string] `json:"itemType,omitempty"`
	Status      *FieldVal[string] `json:"status,omitempty"`
	TargetDate  *FieldVal[string] `json:"targetDate,omitempty"`
	Priority    *FieldVal[string] `json:"priority,omitempty"`
	Duration    *FieldVal[int64]  `json:"duration,omitempty"`
	TimeSpent   *FieldVal[int64]  `json:"timeSpent,omitempty"`
	RecTimes    *FieldVal[int64]  `json:"recTimes,omitempty"`
	RecPeriod   *FieldVal[string] `json:"recPeriod,omitempty"`
	RecProgress *FieldVal[int64]  `json:"recProgress,omitempty"`
	ParentID    *FieldVal[string] `json:"parentID,omitempty"`
}

type FieldVal[T any] struct {
	Val T     `json:"val"`
	Cl  int64 `json:"cl"`
}

func ProcessCmd(msg []byte, userID string) error {
	var o Op
	err := json.Unmarshal(msg, &o)
	if err != nil {
		fmt.Printf("Failed to read msg: %v\n", err)
		return err
	}

	switch o.Op {
	case "UPDATE":
		err = updateRecord(userID, o.ItemID, o.Diffs)
		return err
	case "INSERT":
		err = insertRecord(userID, o.ItemID, o.Diffs)
		return err
	case "DELETE":
		err = deleteRecord(userID, o.ItemID, o.Cl.Int64)
		return err
	default:
		return errors.New("cmd not found")
	}
}
