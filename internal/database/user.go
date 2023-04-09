package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
)

// InsertNewUser inserts a new user into the database
func InsertNewUser(conn *sql.DB, user models.User) int64 {
	stmt, err := conn.Prepare(insertNewUserSQL)

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber, user.Country, user.JobType, user.IncomeRange, user.RefreshToken, user.UserVerificationToken, user.IsVerified)

	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	return id
}

// GetUserWithEmail gets a user from the database with the given email
func GetUserWithEmail(conn *sql.DB, email string) (*models.User, error) {
	var user models.User
	getUserWithEmailSQL := `SELECT * FROM user WHERE email = ?;`

	stmt, err := conn.Prepare(getUserWithEmailSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := stmt.Query(email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.PhoneNumber, &user.Country, &user.JobType, &user.IncomeRange, &user.RefreshToken, &user.UserVerificationToken, &user.IsVerified)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// CheckIfUserExists checks if a user exists in the database with the given email and returns a boolean
func CheckIfUserExists(conn *sql.DB, email string) bool {
	checkUserSQL := `SELECT * FROM user WHERE email = ?;`
	stmt, err := conn.Prepare(checkUserSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	rows, err := stmt.Query(email)

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer rows.Close()

	return rows.Next()
}

// UpdateUserPasswordInDB updates the password of the user in the database with the given email
func UpdateUserPasswordInDB(conn *sql.DB, email string, password string) error {
	updateUserPasswordSQL := `UPDATE user SET password = ? WHERE email = ?`

	stmt, err := conn.Prepare(updateUserPasswordSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(password, email)

	if err != nil {
		return err
	}

	return nil
}
