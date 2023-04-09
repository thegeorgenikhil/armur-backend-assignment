package database

const (
	createUserTableSQL = `CREATE TABLE IF NOT EXISTS user (
	"uid" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"first_name" TEXT,
	"last_name" TEXT,
	"email" TEXT NOT NULL,
	"password" TEXT NOT NULL,
	"phone_number" TEXT NOT NULL,
	"country" TEXT NOT NULL,
	"job_type" TEXT,
	"income_range" TEXT,
	"refresh_token" TEXT,
	"user_verification_token" TEXT,
	"is_verified" INTEGER
	);`

	insertNewUserSQL = `INSERT INTO user (
		first_name, 
		last_name, 
		email, 
		password, 
		phone_number, 
		country,
		job_type,
		income_range,
		refresh_token,
		user_verification_token,
		is_verified
	) values(?,?,?,?,?,?,?,?,?,?,?);`

	createResetPasswordTableSQL = `CREATE TABLE IF NOT EXISTS password (
	"uid" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"email" TEXT NOT NULL,
	"reset_password_token" TEXT NOT NULL,
	"reset_password_token_expiry" INTEGER NOT NULL
	);`

	insertResetPasswordTokenSQL = `INSERT INTO password (
		email,
		reset_password_token,
		reset_password_token_expiry
	) values(?,?,?);`
)
