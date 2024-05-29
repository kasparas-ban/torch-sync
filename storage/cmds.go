package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Op struct {
	Cmd  string `json:"cmd"` // INSERT, UPDATE or DELETE
	Data Item   `json:"data"`
}

type Msg struct {
	Item
}

func ProcessCmd(msg []byte, userID string) error {
	var o Op
	err := json.Unmarshal(msg, &o)
	if err != nil {
		fmt.Printf("Failed to read msg: %v\n", err)
		return err
	}

	switch o.Cmd {
	case "UPDATE":
		err = updateRecord(o.Data, userID)
		return err
	case "INSERT":
		err = insertRecord(o.Data, userID)
		return err
	case "DELETE":
		err = deleteRecord(o.Data.ItemID, userID)
		return err
	default:
		return errors.New("cmd not found")
	}
}

func insertRecord(item Item, userID string) error {
	query, args := buildInsertQuery(item, userID)

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Disable the trigger for this transaction
	_, err = tx.Exec("SET LOCAL custom.disable_trigger = 'true'")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Execute INSERT command
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction, which will automatically reset the custom setting
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func updateRecord(item Item, userID string) error {
	query, args := buildUpdateQuery(item, userID)

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Disable the trigger for this transaction
	_, err = tx.Exec("SET LOCAL custom.disable_trigger = 'true'")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Execute UPDATE command
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction, which will automatically reset the custom setting
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func deleteRecord(itemID string, userID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Disable the trigger for this transaction
	_, err = tx.Exec("SET LOCAL custom.disable_trigger = 'true'")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Execute DELETE command
	_, err = tx.Exec("DELETE FROM items WHERE user_id = $1 AND item_id = $2", userID, itemID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction, which will automatically reset the custom setting
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateQuery(item Item, userID string) (string, []interface{}) {
	var setClauses []string
	var args []interface{}
	argID := 1

	if item.Title.Valid {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argID))
		args = append(args, item.Title.String)
		argID++
	}
	if item.Status.Valid {
		setClauses = append(setClauses, fmt.Sprintf("status_ = $%d", argID))
		args = append(args, item.Status.String)
		argID++
	}
	if item.TargetDate.Valid {
		setClauses = append(setClauses, fmt.Sprintf("target_date = $%d", argID))
		args = append(args, item.TargetDate.String)
		argID++
	}
	if item.Priority.Valid {
		setClauses = append(setClauses, fmt.Sprintf("priority_ = $%d", argID))
		args = append(args, item.Priority.String)
		argID++
	}
	if item.Duration.Valid {
		setClauses = append(setClauses, fmt.Sprintf("duration = $%d", argID))
		args = append(args, item.Duration.Int64)
		argID++
	}
	if item.TimeSpent.Valid {
		setClauses = append(setClauses, fmt.Sprintf("time_spent = $%d", argID))
		args = append(args, item.TimeSpent.Int64)
		argID++
	}
	if item.RecTimes.Valid {
		setClauses = append(setClauses, fmt.Sprintf("rec_times = $%d", argID))
		args = append(args, item.RecTimes.Int64)
		argID++
	}
	if item.RecPeriod.Valid {
		setClauses = append(setClauses, fmt.Sprintf("rec_period = $%d", argID))
		args = append(args, item.RecPeriod.String)
		argID++
	}
	if item.RecProgress.Valid {
		setClauses = append(setClauses, fmt.Sprintf("rec_progress = $%d", argID))
		args = append(args, item.RecProgress.Int64)
		argID++
	}
	if item.RecUpdatedAt.Valid {
		setClauses = append(setClauses, fmt.Sprintf("rec_updated_at = $%d", argID))
		args = append(args, item.RecUpdatedAt.String)
		argID++
	}
	if item.ParentID.Valid {
		setClauses = append(setClauses, fmt.Sprintf("parent_id = $%d", argID))
		args = append(args, item.ParentID.String)
		argID++
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE items SET %s WHERE user_id = $%d AND item_id = $%d", setClause, argID, argID+1)
	args = append(args, userID)
	args = append(args, item.ItemID)

	return query, args
}

func buildInsertQuery(item Item, userID string) (string, []interface{}) {
	var setColName []string
	var cmdArgs []string
	var args []interface{}
	argID := 1

	setColName = append(setColName, "item_id")
	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	args = append(args, item.ItemID)
	argID++

	setColName = append(setColName, "user_id")
	cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
	args = append(args, userID)
	argID++

	if item.ItemType.Valid {
		setColName = append(setColName, "type_")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.ItemType.String)
		argID++
	}
	if item.Title.Valid {
		setColName = append(setColName, "title")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.Title.String)
		argID++
	}
	if item.Status.Valid {
		setColName = append(setColName, "status_")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.Status.String)
		argID++
	}
	if item.TargetDate.Valid {
		setColName = append(setColName, "target_date")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.TargetDate.String)
		argID++
	}
	if item.Priority.Valid {
		setColName = append(setColName, "priority_")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.Priority.String)
		argID++
	}
	if item.Duration.Valid {
		setColName = append(setColName, "duration")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.Duration.Int64)
		argID++
	}
	if item.TimeSpent.Valid {
		setColName = append(setColName, "time_spent")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.TimeSpent.Int64)
		argID++
	}
	if item.RecTimes.Valid {
		setColName = append(setColName, "rec_times")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.RecTimes.Int64)
		argID++
	}
	if item.RecPeriod.Valid {
		setColName = append(setColName, "rec_period")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.RecPeriod.String)
		argID++
	}
	if item.RecProgress.Valid {
		setColName = append(setColName, "rec_progress")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.RecProgress.Int64)
		argID++
	}
	if item.RecUpdatedAt.Valid {
		setColName = append(setColName, "rec_updated_at")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.RecUpdatedAt.String)
		argID++
	}
	if item.ParentID.Valid {
		setColName = append(setColName, "parent_id")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, item.ParentID.String)
		argID++
	}

	query := fmt.Sprintf("INSERT INTO items (%s) VALUES (%v)", strings.Join(setColName, ","), strings.Join(cmdArgs, ","))

	return query, args
}
