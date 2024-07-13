package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type updateCol struct {
	Column string
	Value  interface{}
	Clock  int64
}

type Diffs struct {
	Title        *FieldVal[string] `json:"title,omitempty"`
	Status       *FieldVal[string] `json:"status,omitempty"`
	TargetDate   *FieldVal[string] `json:"target_date,omitempty"`
	Priority     *FieldVal[string] `json:"priority,omitempty"`
	Duration     *FieldVal[int64]  `json:"duration,omitempty"`
	TimeSpent    *FieldVal[int64]  `json:"time_spent,omitempty"`
	RecTimes     *FieldVal[int64]  `json:"rec_times,omitempty"`
	RecPeriod    *FieldVal[string] `json:"rec_period,omitempty"`
	RecProgress  *FieldVal[int64]  `json:"rec_progress,omitempty"`
	RecUpdatedAt *FieldVal[int64]  `json:"rec_updated_at,omitempty"`
	ParentID     *FieldVal[string] `json:"parent_id,omitempty"`
}

type FieldVal[T any] struct {
	Val T     `json:"val"`
	Cl  int64 `json:"cl"`
}

func (d *Diffs) UnmarshalJSON(b []byte) error {
	type Alias Diffs
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
func updateRecord(userID string, itemID string, diffs Diffs, wsId string) error {
	updateCols := getUpdateData(diffs)

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

	// Get all relevant clocks from DB
	clockQuery, clockArgs, clockPointers, clockValues := getClockQuery(userID, itemID, updateCols)
	err = tx.QueryRow(clockQuery, clockArgs...).Scan(clockPointers...)
	if err != nil {
		return err
	}

	// Determine which columns need to be updated
	var newCols []updateCol
	for idx, clock := range clockValues {
		if clock <= updateCols[idx].Clock {
			newCols = append(newCols, updateCols[idx])
		}
	}

	if len(newCols) == 0 {
		tx.Rollback()
		return errors.New("no data to update")
	}

	// Build the update query
	query, args := getUpdateQuery(userID, itemID, newCols)

	// Execute UPDATE command
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

func getUpdateData(diffs Diffs) []updateCol {
	var updateCols []updateCol

	if diffs.Title != nil {
		updateCols = append(updateCols, updateCol{"title", diffs.Title.Val, diffs.Title.Cl})
	}
	if diffs.Status != nil {
		updateCols = append(updateCols, updateCol{"status", diffs.Status.Val, diffs.Status.Cl})
	}
	if diffs.TargetDate != nil {
		updateCols = append(updateCols, updateCol{"target_date", diffs.TargetDate.Val, diffs.TargetDate.Cl})
	}
	if diffs.Priority != nil {
		updateCols = append(updateCols, updateCol{"priority", diffs.Priority.Val, diffs.Priority.Cl})
	}
	if diffs.Duration != nil {
		updateCols = append(updateCols, updateCol{"duration", diffs.Duration.Val, diffs.Duration.Cl})
	}
	if diffs.TimeSpent != nil {
		updateCols = append(updateCols, updateCol{"time_spent", diffs.TimeSpent.Val, diffs.TimeSpent.Cl})
	}
	if diffs.RecTimes != nil {
		updateCols = append(updateCols, updateCol{"rec_times", diffs.RecTimes.Val, diffs.RecTimes.Cl})
	}
	if diffs.RecPeriod != nil {
		updateCols = append(updateCols, updateCol{"rec_period", diffs.RecPeriod.Val, diffs.RecPeriod.Cl})
	}
	if diffs.RecProgress != nil {
		updateCols = append(updateCols, updateCol{"rec_progress", diffs.RecProgress.Val, diffs.RecProgress.Cl})
	}
	if diffs.RecUpdatedAt != nil {
		updateCols = append(updateCols, updateCol{"rec_updated_at", diffs.RecUpdatedAt.Val, diffs.RecUpdatedAt.Cl})
	}
	if diffs.ParentID != nil {
		updateCols = append(updateCols, updateCol{"parent_id", diffs.ParentID.Val, diffs.ParentID.Cl})
	}

	return updateCols
}

func getClockQuery(userID string, itemID string, updateCols []updateCol) (string, []interface{}, []interface{}, []int64) {
	// Column value pointer for scanning
	columns := make([]int64, len(updateCols))
	columnPointers := make([]interface{}, len(updateCols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	// Get all clock columns
	var clockCols []string
	for _, col := range updateCols {
		clockCols = append(clockCols, col.Column+"__c")
	}

	query := fmt.Sprintf("SELECT %s FROM items WHERE user_id = $1 AND item_id = $2", strings.Join(clockCols, ", "))
	args := []interface{}{userID, itemID}

	return query, args, columnPointers, columns
}

func getUpdateQuery(userID string, itemID string, cols []updateCol) (string, []interface{}) {
	var setClause []string
	var args []interface{}

	clauseCounter := 1
	for i := 0; i < len(cols); i += 1 {
		c := cols[i]
		setClause = append(setClause, fmt.Sprintf(c.Column+" = $%d", clauseCounter))
		args = append(args, c.Value)
		clauseCounter++

		setClause = append(setClause, fmt.Sprintf(c.Column+"__c = $%d", clauseCounter))
		args = append(args, c.Clock+1)
		clauseCounter++
	}

	clause := strings.Join(setClause, ", ")
	query := fmt.Sprintf("UPDATE items SET %s WHERE user_id = $%d AND item_id = $%d", clause, 2*len(cols)+1, 2*len(cols)+2)
	args = append(args, userID)
	args = append(args, itemID)

	return query, args
}
