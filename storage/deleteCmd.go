package storage

import "errors"

func deleteRecord(userID string, itemID string, cl int64) error {
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

	// Compare item clocks
	var itemClock int64
	err = tx.QueryRow("SELECT item__c FROM items WHERE user_id = $1 AND item_id = $2", userID, itemID).Scan(&itemClock)
	if err != nil {
		tx.Rollback()
		return err
	}
	if itemClock >= cl {
		tx.Rollback()
		return errors.New("delete aborted: outdated command")
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
