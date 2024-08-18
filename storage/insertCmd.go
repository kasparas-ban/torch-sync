package storage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

type InsertData struct {
	Title       string
	Item_type   string
	Status      *string
	Target_date *string
	Priority    *string
	Duration    *int64
	Parent_id   *string
	Time_spent  *int64
	Created_at  *string
}

func (d *InsertData) UnmarshalJSON(b []byte) error {
	type Alias InsertData
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	return nil
}

func insertRecord(userID string, itemID string, data InsertData, wsId string) error {
	query, args := buildInsertQuery(userID, itemID, data)

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Set WebSocket ID
	_, err = tx.Exec(fmt.Sprintf(`SET custom.ws_id TO '%s'`, wsId))
	if err != nil {
		tx.Rollback()
		return err
	}

	slog.Info("Executing INSERT", "query", query, "args", args)

	// Execute INSERT command
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func buildInsertQuery(userID string, itemID string, data InsertData) (string, []interface{}) {
	var setColName []string
	var cmdArgs []string
	var args []interface{}
	argID := 1

	setColName = append(setColName, "item_id")
	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	args = append(args, itemID)
	argID++

	setColName = append(setColName, "user_id")
	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	args = append(args, userID)
	argID++

	setColName = append(setColName, "title")
	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	args = append(args, data.Title)
	argID++

	setColName = append(setColName, "item_type")
	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	args = append(args, data.Item_type)
	argID++

	if data.Status != nil {
		setColName = append(setColName, "status")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, *data.Status)
		argID++
	}
	if data.Target_date != nil {
		setColName = append(setColName, "target_date")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, *data.Target_date)
		argID++
	}
	if data.Priority != nil {
		setColName = append(setColName, "priority")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, *data.Priority)
		argID++
	}
	if data.Duration != nil {
		setColName = append(setColName, "duration")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, data.Duration)
		argID++
	}
	if data.Time_spent != nil {
		setColName = append(setColName, "time_spent")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, *data.Time_spent)
		argID++
	}
	// if data.RecTimes != nil {
	// 	setColName = append(setColName, "rec_times")
	// 	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	// 	args = append(args, data.RecTimes.Val)
	// 	argID++
	// }
	// if data.RecPeriod != nil {
	// 	setColName = append(setColName, "rec_period")
	// 	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	// 	args = append(args, data.RecPeriod.Val)
	// 	argID++
	// }
	// if data.RecProgress != nil {
	// 	setColName = append(setColName, "rec_progress")
	// 	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	// 	args = append(args, data.RecProgress.Val)
	// 	argID++
	// }
	if data.Parent_id != nil {
		setColName = append(setColName, "parent_id")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, *data.Parent_id)
		argID++
	}
	if data.Created_at != nil {
		setColName = append(setColName, "created_at")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, *data.Created_at)
		argID++
	}

	query := fmt.Sprintf("INSERT INTO items (%s) VALUES (%v)", strings.Join(setColName, ", "), strings.Join(cmdArgs, ", "))

	return query, args
}
