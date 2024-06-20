package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"torch/torch-sync/types"
)

type Op struct {
	Op     string          `json:"op"` // INSERT, UPDATE or DELETE
	ItemID string          `json:"item_id"`
	Diffs  *Diffs          `json:"diffs"`
	Data   *InsertData     `json:"data"`
	Cl     types.NullInt64 `json:"cl,omitempty"`
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
		if o.Diffs == nil {
			return errors.New("incorrect msg body format")
		}
		err = updateRecord(userID, o.ItemID, *o.Diffs)
		return err
	case "INSERT":
		if o.Data == nil {
			return errors.New("incorrect msg body format")
		}
		err = insertRecord(userID, o.ItemID, *o.Data)
		return err
	case "DELETE":
		err = deleteRecord(userID, o.ItemID, o.Cl.Int64)
		return err
	default:
		return errors.New("cmd not found")
	}
}
