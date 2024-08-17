package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"torch/torch-sync/types"
)

type Op struct {
	Op    string          `json:"op"` // INSERT, UPDATE or DELETE
	Table string          `json:"table"`
	RowID string          `json:"row_id"`
	Diffs *Diffs          `json:"diffs"`
	Data  *InsertData     `json:"data"`
	Cl    types.NullInt64 `json:"cl,omitempty"`
}

func ProcessCmd(msg []byte, userID string, wsID string) error {
	var o Op
	err := json.Unmarshal(msg, &o)
	if err != nil {
		fmt.Printf("Failed to read msg: %v\n", err)
		return err
	}

	switch o.Op {
	case "UPDATE":
		slog.Info("UPDATE", "op", o)
		if o.Diffs == nil {
			return errors.New("incorrect msg body format")
		}
		if o.Table == "users" {
			return updateUserRecord(userID, *o.Diffs, wsID)
		}
		err = updateRecord(userID, o.RowID, *o.Diffs, wsID)
		return err
	case "INSERT":
		slog.Info("INSERT", "op", o)
		if o.Data == nil {
			return errors.New("incorrect msg body format")
		}
		err = insertRecord(userID, o.RowID, *o.Data, wsID)
		return err
	case "DELETE":
		slog.Info("DELETE", "op", o)
		err = deleteRecord(userID, o.RowID, o.Cl.Int64, wsID)
		return err
	default:
		return errors.New("cmd not found")
	}
}
