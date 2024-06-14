package storage

type Item struct {
	User_id string `json:"user_id"`
	ItemResponse
}

type ItemResponse struct {
	Item_id        string  `json:"item_id"`
	Title          string  `json:"title"`
	Item_type      string  `json:"item_type"`
	Status         string  `json:"status"`
	Target_date    *string `json:"target_date"`
	Priority       *string `json:"priority"`
	Duration       *string `json:"duration"`
	Time_spent     int64   `json:"time_spent"`
	Rec_times      *int64  `json:"rec_times"`
	Rec_period     *string `json:"rec_period"`
	Rec_progress   *int64  `json:"rec_progress"`
	Rec_updated_at *string `json:"rec_updated_at"`
	Parent_id      *string `json:"parent_id"`
	Updated_at     string  `json:"updated_at"`
	Created_at     string  `json:"created_at"`
	// Clocks
	Title__c        int64 `json:"title__c"`
	Status__c       int64 `json:"status__c"`
	Target_date__c  int64 `json:"target_date__c"`
	Priority__c     int64 `json:"priority__c"`
	Duration__c     int64 `json:"duration__c"`
	Time_spent__c   int64 `json:"time_spent__c"`
	Rec_times__c    int64 `json:"rec_times__c"`
	Rec_period__c   int64 `json:"rec_period__c"`
	Rec_progress__c int64 `json:"rec_progress__c"`
	Parent_id__c    int64 `json:"parent_id__c"`
	Item__c         int64 `json:"item__c"`
}

func GetAllItemsByUser(userID string) ([]ItemResponse, error) {
	rows, err := DB.Query(`
		SELECT * 
		FROM items 
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ItemResponse
	for rows.Next() {
		var item ItemResponse
		var user_id string

		err := rows.Scan(&item.Item_id, &user_id, &item.Title, &item.Item_type, &item.Status,
			&item.Target_date, &item.Priority, &item.Duration, &item.Time_spent, &item.Rec_times,
			&item.Rec_period, &item.Rec_progress, &item.Rec_updated_at, &item.Parent_id, &item.Updated_at,
			&item.Created_at, &item.Title__c, &item.Status__c, &item.Target_date__c, &item.Priority__c,
			&item.Duration__c, &item.Time_spent__c, &item.Rec_times__c, &item.Rec_period__c, &item.Rec_progress__c,
			&item.Parent_id__c, &item.Item__c)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
