package database

import (
	"database/sql"
	"log"
)

// UpdateRefreshTokenInDB updates the refresh token in the database
func UpdateRefreshTokenInDB(conn *sql.DB, email string, token string) error {
	updateRefreshTokenSQL := `UPDATE user SET refresh_token = ? WHERE email = ?`

	stmt, err := conn.Prepare(updateRefreshTokenSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(token, email)

	if err != nil {
		return err
	}

	return nil
}

// RemoveRefreshTokenFromDB removes the refresh token from the database
func RemoveRefreshTokenFromDB(conn *sql.DB, email string) error {
	removeRefreshTokenSQL := `UPDATE user SET refresh_token = NULL WHERE email = ?`
	stmt, err := conn.Prepare(removeRefreshTokenSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(email)

	if err != nil {
		return err
	}

	return nil
}

// VerifyUserInDB sets the is_verified column to 1 and removes the user_verification_token from the database after the user has verified their email
func VerifyUserInDB(conn *sql.DB, email string) error {
	verifyUserSQL := `UPDATE user SET is_verified = 1, user_verification_token = NULL WHERE uid = ?`
	stmt, err := conn.Prepare(verifyUserSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(email)

	if err != nil {
		return err
	}

	return nil
}

// UpdateVerificationTokenInDB updates the user_verification_token in the database after the user has requested a new verification email
func UpdateVerificationTokenInDB(conn *sql.DB, email string, token string) error {
	updateVerificationTokenSQL := `UPDATE user SET user_verification_token = ? WHERE email = ?`

	stmt, err := conn.Prepare(updateVerificationTokenSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(token, email)

	if err != nil {
		return err
	}

	return nil
}

// InsertResetPasswordTokenInDB inserts the reset_password_token and reset_password_token_expiry into the database when the user requests a password reset
func InsertResetPasswordTokenInDB(conn *sql.DB, email string, token string, tokenExpiry int64) error {
	stmt, err := conn.Prepare(insertResetPasswordTokenSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(email, token, tokenExpiry)

	if err != nil {
		return err
	}

	return nil
}

// GetResetPasswordTokenFromDB gets the reset_password_token and reset_password_token_expiry from the database
func GetResetPasswordTokenFromDB(conn *sql.DB, email string) (string, int64, error) {
	getResetPasswordTokenSQL := `SELECT reset_password_token,reset_password_token_expiry FROM password WHERE email = ?`

	stmt, err := conn.Prepare(getResetPasswordTokenSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	var token string
	var tokenExpiry int64

	err = stmt.QueryRow(email).Scan(&token, &tokenExpiry)

	if err != nil {
		return "", 0, err
	}

	return token, tokenExpiry, nil
}

// RemoveResetPasswordTokenFromDB removes the reset_password_token and reset_password_token_expiry from the database after the user has reset their password
func RemoveResetPasswordTokenFromDB(conn *sql.DB, email string) error {
	removeResetPasswordTokenSQL := `DELETE FROM password WHERE email = ?`
	stmt, err := conn.Prepare(removeResetPasswordTokenSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(email)

	if err != nil {
		return err
	}

	return nil
}
