package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"torch/torch-sync/types"
)

type Op struct {
	Op     string          `json:"op"` // INSERT, UPDATE or DELETE
	ItemID string          `json:"item_id"`
	Diffs  *Diffs          `json:"diffs"`
	Data   *InsertData     `json:"data"`
	Cl     types.NullInt64 `json:"cl,omitempty"`
}

func ProcessCmd(msg []byte, userID string, wsId string) error {
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
		err = updateRecord(userID, o.ItemID, *o.Diffs, wsId)
		slog.Info("UPDATE", "op", o)
		return err
	case "INSERT":
		if o.Data == nil {
			return errors.New("incorrect msg body format")
		}
		err = insertRecord(userID, o.ItemID, *o.Data, wsId)
		slog.Info("INSERT", "op", o)
		return err
	case "DELETE":
		err = deleteRecord(userID, o.ItemID, o.Cl.Int64, wsId)
		slog.Info("DELETE", "op", o)
		return err
	default:
		return errors.New("cmd not found")
	}
}
