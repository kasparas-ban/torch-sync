package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"torch/torch-sync/types"
)

type Op struct {
	Cmd  string `json:"cmd"` // INSERT, UPDATE or DELETE
	Data Item   `json:"data"`
}

type Msg struct {
	Item
}

type Item struct {
	ItemID       string           `json:"itemID"`
	Title        types.NullString `json:"title"`
	ItemType     types.NullString `json:"itemType"`
	Status       types.NullString `json:"status"`
	TargetDate   types.NullString `json:"targetDate"`
	Priority     types.NullString `json:"priority"`
	Duration     types.NullInt64  `json:"duration"`
	TimeSpent    types.NullInt64  `json:"timeSpent"`
	RecTimes     types.NullInt64  `json:"recTimes"`
	RecPeriod    types.NullString `json:"recPeriod"`
	RecProgress  types.NullInt64  `json:"recProgress"`
	RecUpdatedAt types.NullString `json:"recUpdatedAt"`
	ParentID     types.NullString `json:"parentID"`
	UpdatedAt    types.NullString `json:"updatedAt"`
	CreatedAt    types.NullString `json:"createdAt"`
}

func ProcessCmd(msg []byte) error {
	var o Op
	err := json.Unmarshal(msg, &o)
	if err != nil {
		fmt.Printf("Failed to read msg: %v\n", err)
		return err
	}

	switch o.Cmd {
	case "UPDATE":
		err = updateRecord(o.Data)
		return err
	default:
		return errors.New("cmd not found")
	}
}

func updateRecord(item Item) error {
	query, args := buildUpdateQuery(item)

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

	// Execute your UPDATE command
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

func buildUpdateQuery(item Item) (string, []interface{}) {
	var setClauses []string
	var args []interface{}
	argID := 1

	if item.Title.Valid {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argID))
		args = append(args, item.Title.String)
		argID++
	}
	if item.Status.Valid {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argID))
		args = append(args, item.Status.String)
		argID++
	}
	if item.TargetDate.Valid {
		setClauses = append(setClauses, fmt.Sprintf("target_date = $%d", argID))
		args = append(args, item.TargetDate.String)
		argID++
	}
	if item.Priority.Valid {
		setClauses = append(setClauses, fmt.Sprintf("priority = $%d", argID))
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
	query := fmt.Sprintf("UPDATE items SET %s WHERE item_id = $%d", setClause, argID)
	args = append(args, item.ItemID)

	return query, args
}
