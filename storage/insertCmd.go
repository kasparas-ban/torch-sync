package storage

import (
	"fmt"
	"strings"
)

func insertRecord(userID string, itemID string, diffs Diffs) error {
	query, args := buildInsertQuery(userID, itemID, diffs)

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
func buildInsertQuery(userID string, itemID string, diffs Diffs) (string, []interface{}) {
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

	if diffs.ItemType != nil {
		setColName = append(setColName, "item_type")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.ItemType.Val)
		argID++
	}
	if diffs.Title != nil {
		setColName = append(setColName, "title")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.Title.Val)
		argID++
	}
	if diffs.Status != nil {
		setColName = append(setColName, "status")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.Status.Val)
		argID++
	}
	if diffs.TargetDate != nil {
		setColName = append(setColName, "target_date")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.TargetDate.Val)
		argID++
	}
	if diffs.Priority != nil {
		setColName = append(setColName, "priority")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.Priority.Val)
		argID++
	}
	if diffs.Duration != nil {
		setColName = append(setColName, "duration")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.Duration.Val)
		argID++
	}
	if diffs.TimeSpent != nil {
		setColName = append(setColName, "time_spent")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.TimeSpent.Val)
		argID++
	}
	if diffs.RecTimes != nil {
		setColName = append(setColName, "rec_times")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.RecTimes.Val)
		argID++
	}
	if diffs.RecPeriod != nil {
		setColName = append(setColName, "rec_period")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.RecPeriod.Val)
		argID++
	}
	if diffs.RecProgress != nil {
		setColName = append(setColName, "rec_progress")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.RecProgress.Val)
		argID++
	}
	if diffs.ParentID != nil {
		setColName = append(setColName, "parent_id")
		cmdArgs = append(cmdArgs, fmt.Sprintf("$%d", argID))
		args = append(args, diffs.ParentID.Val)
		argID++
	}

	query := fmt.Sprintf("INSERT INTO items (%s) VALUES (%v)", strings.Join(setColName, ","), strings.Join(cmdArgs, ","))

	return query, args
}
