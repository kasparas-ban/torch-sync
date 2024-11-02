package storage

import (
	"context"
	"errors"
	"fmt"
	o "torch/torch-sync/optional"
	"torch/torch-sync/pkg"
)

type User struct {
	UserID      string       `json:"user_id"`
	ClerkID     string       `json:"clerk_id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Birthday    o.NullString `json:"birthday"`
	Gender      o.NullString `json:"gender"`
	CountryCode o.NullString `json:"country_code"`
	City        o.NullString `json:"city"`
	Description o.NullString `json:"description"`
	FocusTime   uint64       `json:"focus_time"`
	FocusTimeCl uint64       `json:"focus_time__c"`
	UpdatedAt   string       `json:"updated_at"`
	CreatedAt   string       `json:"created_at"`
}

type ExistingUser struct {
	User
	CountryCode o.NullString `json:"country_code"`
}

type NewUser struct {
	ClerkID  string
	Email    string
	Username string
}

type UpdateUserReq struct {
	Username    string       `json:"username" validate:"required,gt=5,lt=21"`
	Birthday    o.NullString `json:"birthday"`
	Gender      o.NullString `json:"gender"`
	CountryCode o.NullString `json:"countryCode" validate:"lt=3"`
	City        o.NullString `json:"city"`
	Description o.NullString `json:"description"`
}

type RegisterUserReq struct {
	Username    string       `json:"username" validate:"required,gt=5,lt=21"`
	Email       string       `json:"email" validate:"required,email"`
	Password    string       `json:"password" validate:"required,password"`
	Birthday    o.NullString `json:"birthday"`
	Gender      o.NullString `json:"gender"`
	CountryCode o.NullString `json:"countryCode" validate:"lt=3"`
	City        o.NullString `json:"city"`
	Description o.NullString `json:"description"`
}

type UpdateUserEmailReq struct {
	Email string `json:"email" validate:"email"`
}

func GetUserByClerkID(clerkID string) (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description, u.focus_time, u.focus_time__c, u.updated_at, u.created_at
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.clerk_id = $1
		LIMIT 1
	`, clerkID).Scan(&user.UserID, &user.ClerkID, &user.Username,
		&user.Email, &user.Birthday, &user.Gender, &user.CountryCode,
		&user.City, &user.Description, &user.FocusTime, &user.FocusTimeCl,
		&user.UpdatedAt, &user.CreatedAt)

	return user, err
}

func GetUser(userID string) (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description, u.focus_time, u.focus_time__c, u.updated_at, u.created_at
		FROM users u
		LEFT JOIN countries c ON u.country_id = c.country_id
		WHERE u.user_id = $1
		LIMIT 1
	`, userID).Scan(&user.UserID, &user.ClerkID, &user.Username,
		&user.Email, &user.Birthday, &user.Gender, &user.CountryCode,
		&user.City, &user.Description, &user.FocusTime, &user.FocusTimeCl,
		&user.UpdatedAt, &user.CreatedAt)

	return user, err
}

func AddUser(u NewUser) (ExistingUser, error) {
	ctx := context.Background()

	var newUser ExistingUser
	userID, err := pkg.NewRandomID()
	if err != nil {
		return newUser, err
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return newUser, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SET custom.ws_id TO '0'`)
	if err != nil {
		tx.Rollback()
		return newUser, err
	}

	// Add user
	_, err = tx.ExecContext(ctx, "INSERT INTO users (user_id, clerk_id, username, email) VALUES ($1, $2, $3, $4)",
		userID, u.ClerkID, u.Username, u.Email)
	if err != nil {
		return newUser, err
	}

	// Select newly added user
	err = tx.QueryRowContext(ctx, `
	        SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description, u.focus_time, u.focus_time__c, u.updated_at, u.created_at
			FROM users u
			LEFT JOIN countries c ON u.country_id = c.country_id
			WHERE u.user_id = $1 LIMIT 1`, userID).Scan(
		&newUser.UserID, &newUser.ClerkID, &newUser.Username,
		&newUser.Email, &newUser.Birthday, &newUser.Gender, &newUser.CountryCode,
		&newUser.City, &newUser.Description, &newUser.FocusTime, &newUser.FocusTimeCl,
		&newUser.UpdatedAt, &newUser.CreatedAt)
	if err != nil {
		return newUser, err
	}

	if err = tx.Commit(); err != nil {
		return newUser, err
	}

	return newUser, err
}

func RegisterUser(u RegisterUserReq, clerkID string) (ExistingUser, error) {
	ctx := context.Background()

	var newUser ExistingUser
	userID, err := pkg.NewRandomID()
	if err != nil {
		return newUser, err
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return newUser, err
	}
	defer tx.Rollback()

	// Get country ID
	var countryId o.NullUint
	if u.CountryCode.IsValid && u.CountryCode.Val != "" {
		err := tx.QueryRowContext(ctx, `
				SELECT country_id FROM countries WHERE country_code = $1
			`, u.CountryCode.Val).Scan(&countryId)
		if err != nil {
			return newUser, err
		}
	}

	// Add user
	_, err = tx.ExecContext(ctx, `
			INSERT INTO users (user_id, clerk_id, username, email, birthday, gender, country_id, city, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, userID, clerkID, u.Username, u.Email, u.Birthday, u.Gender, countryId, u.City, u.Description)
	if err != nil {
		return newUser, err
	}

	// Select new user
	err = tx.QueryRowContext(ctx, `
			SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description, u.focus_time, u.focus_time__c, u.updated_at, u.created_at
			FROM users u
			LEFT JOIN countries c ON u.country_id = c.country_id
			WHERE u.user_id = $1 LIMIT 1
		`, userID).Scan(&newUser.UserID, &newUser.ClerkID, &newUser.Username,
		&newUser.Email, &newUser.Birthday, &newUser.Gender, &newUser.CountryCode,
		&newUser.City, &newUser.Description, &newUser.FocusTime, &newUser.FocusTimeCl,
		&newUser.UpdatedAt, &newUser.CreatedAt)
	if err != nil {
		return newUser, err
	}

	if err = tx.Commit(); err != nil {
		return newUser, err
	}

	return newUser, err
}

func UpdateUser(userID string, u UpdateUserReq) (ExistingUser, error) {
	ctx := context.Background()

	var updatedUser ExistingUser
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return updatedUser, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SET custom.ws_id TO '0'`)
	if err != nil {
		tx.Rollback()
		return updatedUser, err
	}

	// Get country ID
	var countryId o.NullUint
	if u.CountryCode.IsValid && u.CountryCode.Val != "" {
		err := tx.QueryRowContext(ctx, `
				SELECT country_id FROM countries WHERE country_code = $1
			`, u.CountryCode.Val).Scan(&countryId)
		if err != nil {
			fmt.Println("Error country", err)
			return updatedUser, err
		}
	}

	// Update user
	_, err = tx.ExecContext(ctx, `
			UPDATE users
			SET username = $1, birthday = $2, gender = $3, country_id = $4, city = $5, description = $6
			WHERE user_id = $7
		`, u.Username, u.Birthday, u.Gender, countryId, u.City, u.Description, userID)
	if err != nil {
		fmt.Println("Error adding", err, u.Username, userID)
		return updatedUser, err
	}

	// Select updated user
	err = tx.QueryRowContext(ctx, `
			SELECT u.user_id, u.clerk_id, u.username, u.email, u.birthday, u.gender, c.country_code, u.city, u.description, u.focus_time, u.focus_time__c, u.updated_at, u.created_at
	 		FROM users u
	 		LEFT JOIN countries c ON u.country_id = c.country_id
	 		WHERE u.user_id = $1 LIMIT 1
		`, userID).Scan(&updatedUser.UserID, &updatedUser.ClerkID, &updatedUser.Username,
		&updatedUser.Email, &updatedUser.Birthday, &updatedUser.Gender, &updatedUser.CountryCode,
		&updatedUser.City, &updatedUser.Description, &updatedUser.FocusTime, &updatedUser.FocusTimeCl,
		&updatedUser.UpdatedAt, &updatedUser.CreatedAt)
	if err != nil {
		fmt.Println("Error selecting", err)
		return updatedUser, err
	}

	if err = tx.Commit(); err != nil {
		return updatedUser, err
	}

	return updatedUser, err
}

func DeleteUser(userID string) error {
	ctx := context.Background()

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SET custom.ws_id TO '0'`)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all user items
	_, err = tx.ExecContext(ctx, "DELETE FROM items WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Delete user
	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (u *UpdateUserReq) Validate() error {
	if err := pkg.Validate.Struct(u); err != nil {
		return err
	}

	if u.Gender.IsValid && (u.Gender.Val != "MALE" && u.Gender.Val != "FEMALE" && u.Gender.Val != "OTHER") {
		return errors.New("incorrect gender value")
	}

	if u.CountryCode.IsValid && (len(u.CountryCode.Val) > 2) {
		return errors.New("incorrect country code value")
	}

	return nil
}

func (u *RegisterUserReq) Validate() error {
	if err := pkg.Validate.Struct(u); err != nil {
		return err
	}

	if u.Gender.IsValid && (u.Gender.Val != "MALE" && u.Gender.Val != "FEMALE" && u.Gender.Val != "OTHER") {
		return errors.New("incorrect gender value")
	}

	if u.CountryCode.IsValid && (len(u.CountryCode.Val) > 2) {
		return errors.New("incorrect country code value")
	}

	return nil
}
