package storage

import (
	"torch/torch-sync/types"
)

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

func GetAllItemsByUser(userID string) ([]Item, error) {
	rows, err := DB.Query(`
		SELECT item_id, title, item_type, status, target_date, 
		priority, duration, time_spent, rec_times, rec_period, 
		rec_progress, rec_updated_at, parent_id, updated_at, 
		created_at 
		FROM items WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ItemID, &item.Title, &item.ItemType, &item.Status,
			&item.TargetDate, &item.Priority, &item.Duration, &item.TimeSpent, &item.RecTimes,
			&item.RecPeriod, &item.RecProgress, &item.RecUpdatedAt, &item.ParentID, &item.UpdatedAt,
			&item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
