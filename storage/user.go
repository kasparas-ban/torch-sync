package storage

import "database/sql"

type User struct {
	UserID      string         `json:"userID"`
	ClerkID     string         `json:"clerkID"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	Birthday    sql.NullString `json:"birthday"`
	Gender      sql.NullString `json:"gender"`
	CountryCode sql.NullString `json:"countryCode"`
	City        sql.NullString `json:"city"`
	Description sql.NullString `json:"description"`
	FocusTime   uint64         `json:"focusTime"`
	UpdatedAt   string         `json:"updatedAt"`
	CreatedAt   string         `json:"createdAt"`
}

func GetUserByClerkID(clerkID string) (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description_, u.focus_time, u.updated_at, u.created_at
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.clerk_id = $1
		LIMIT 1
	`, clerkID).Scan(&user.UserID, &user.ClerkID, &user.Username,
		&user.Email, &user.Birthday, &user.Gender, &user.CountryCode,
		&user.City, &user.Description, &user.FocusTime, &user.UpdatedAt,
		&user.CreatedAt)

	return user, err
}

func GetUser(userID string) (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description_, u.focus_time, u.updated_at, u.created_at
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.user_id = $1
		LIMIT 1
	`, userID).Scan(&user.UserID, &user.ClerkID, &user.Username,
		&user.Email, &user.Birthday, &user.Gender, &user.CountryCode,
		&user.City, &user.Description, &user.FocusTime, &user.UpdatedAt,
		&user.CreatedAt)

	return user, err

}
