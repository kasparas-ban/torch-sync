package storage

import (
	"errors"
	"fmt"
	"strings"
)

func updateUserRecord(userID string, diffs Diffs, wsID string) error {
	updateCols := getUpdateUserData(diffs)

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Set WebSocket ID
	_, err = tx.Exec(fmt.Sprintf(`SET custom.ws_id TO '%s'`, wsID))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get all relevant clocks from DB
	clockQuery, clockArgs, clockPointers, clockValues := getUserClockQuery(userID, updateCols)
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
	query, args := getUpdateUserQuery(userID, newCols)

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

func getUpdateUserData(diffs Diffs) []updateCol {
	var updateCols []updateCol

	if diffs.FocusTime != nil {
		updateCols = append(updateCols, updateCol{"focus_time", diffs.FocusTime.Val, diffs.FocusTime.Cl})
	}

	return updateCols
}

func getUserClockQuery(userID string, updateCols []updateCol) (string, []interface{}, []interface{}, []int64) {
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

	query := fmt.Sprintf("SELECT %s FROM users WHERE user_id = $1", strings.Join(clockCols, ", "))
	args := []interface{}{userID}

	return query, args, columnPointers, columns
}

func getUpdateUserQuery(userID string, cols []updateCol) (string, []interface{}) {
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
	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = $%d", clause, 2*len(cols)+1)
	args = append(args, userID)

	return query, args
}
