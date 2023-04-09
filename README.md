# Armur Backend Assignment

## Introduction

This is a backend assignment for ArmurAI. It is a simple Golang based backend server with the following features:

1. Login
2. Logout
3. Register
4. User Activation Email
5. Resend Activation Email
6. Forgot Password
7. Get User Details

The project also uses a frontend for the user to interact with the APIs. The frontend is composed of static HTML files served by the backend file server.

It uses chi for routing and sqlite3 for storing the data.

Packages used:

- `github.com/dgrijalva/jwt-go` v3.2.0+incompatible
- `github.com/go-chi/chi` v1.5.4
- `github.com/go-playground/validator/v10` v10.12.0
- `github.com/joho/godotenv` v1.5.1
- `github.com/mattn/go-sqlite3` v1.14.16
- `golang.org/x/crypto` v0.8.0

## Installation

> For this to work you need to have Go 1.16 or above installed on your system. You can download it from [here](https://golang.org/dl/) or you can use a docker container with Go image.

### Prerequisites

- Go 1.16 or above
- Git

1. Clone the repository

```bash
git clone github.com/thegeorgenikhil/armur-backend-assignment
```

2. Install the dependencies

```bash
go mod download
```

3. Rename the `.env.example` file to `.env` and fill in the required values

```bash
mv .env.example .env
```

_.env file_

```bash
PORT= # Server port
JWT_SECRET= # JWT Secret used for signing the JWT token
FROM_EMAIL= # Email address used for sending emails
FROM_EMAIL_PASSWORD= # App password for send emails using the FROM_EMAIL address
FRONTEND_URL= # URL of the frontend eg: http://localhost:8080 if hosted locally
```

4. Run the server

```bash
go run cmd/main.go
```

5. After starting the server you can access the frontend at `http://localhost:8080`. It will redirect you to the login page.

## Folder Structure

```bash
├── cmd
│   └── main.go            # entry point for the application
├── data.db                # SQLite database file
├── go.mod
├── go.sum
├── internal
│   ├── database           # database access layer
│   ├── handlers           # HTTP handlers for request processing
│   ├── middleware         # HTTP middleware for request processing
│   ├── models             # data models
│   └── routes             # HTTP routes definition
|   └── utils              # utility functions
├── pkg
│   ├── bcrypt             # bcrypt password hashing utility
│   ├── jwt                # JSON Web Token (JWT) utility
│   └── mail               # email utility
└── static                 # HTML templates and static assets
```

## API Endpoints

---

### Register

Creates a new user account and sends an activation email. Returns the JWT access and refresh tokens.

`POST /api/register`

#### Request

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "johndoe@email.com",
  "password": "password",
  "phone_number": "9999999999",
  "country": "Country",
  "job_type": "Job Type",
  "income_range": "200000"
}
```

#### Response

_Status: 201 Created_

```json
{
  "msg": "Signed up successfully. Please verify you email using the confirmation mail",
  "access_token": "jwt-access-token-here",
  "refresh_token": "jwt-refresh-token-here",
  "is_verified": false
}
```

---

### Login

Logs in the user and returns the JWT access and refresh tokens.

`POST /api/login`

#### Request

```json
{
  "email": "johndoe@gmail.com",
  "password": "password"
}
```

#### Response

_Status: 200 OK_

```json
{
  "msg": "Logged In Successfully",
  "access_token": "jwt-access-token-here",
  "refresh_token": "jwt-refresh-token-here",
  "is_verified": false
}
```

---

### Logout

Logs out the user by invalidating the refresh token.

`GET /api/logout`

`Authorization: Bearer <access_token>`

#### Request

No request body

#### Response

_Status: 200 OK_

```json
{
  "msg": "Successfully logged out"
}
```

---

### Refresh Token

Returns a new access token using the refresh token.

`POST /api/refresh`

#### Request

```json
{
  "refresh_token": "jwt-refresh-token-here"
}
```

#### Response

_Status: 200 OK_

```json
{
  "access_token": "jwt-refresh-token-here"
}
```

---

### Verify Activation Email

Verifies the user's email address by sending a verification link

Verification link: `GET /api/user-activation/verify?token=<user-verification-token>`

#### Response

_Status: 200 OK_

```
Verification successful. You can close this tab and continue to the home page
```

---

### Check User Verification Status

Checks if the user's email address is verified and returns a boolean value.

`GET /api/user-activation/check`

`Authorization: Bearer <access_token>`

#### Response

_Status: 200 OK_

```json
{
  "is_user_verified": false
}
```

---

### Resend Verification Email

Resends the verification email to the user's email address. **It has a user specific rate limit of 1 email per 5 minutes.**

`GET /api/user-activation/resend`

`Authorization: Bearer <access_token>`

#### Response

_Status: 200 OK_

```json
{
  "msg": "Resent verification email. Please check your registered email"
}
```

---

### Forgot Password

Sends a password reset link to the user's email address.

`POST /api/forgot`

#### Request

```json
{
  "email": "johndoe@gmail.com"
}
```

#### Response

_Status: 200 OK_

```json
{
  "msg": "Reset password email sent. Please check your email"
}
```

---

### Reset Password

After clicking the password reset link, the user is redirected to the reset password page. The user can enter the new password and submit the form. The user's password is updated in the database.

`POST /api/reset-password`

#### Request

```json
{
  "email": "johndoe@email.com",
  "password": "password",
  "token": "reset-password-token-here"
}
```

#### Response

_Status: 200 OK_

```json
{
  "msg": "Password reset successful. You can now login with your new password"
}
```

---

### Get User Details

Returns the user's details.

`GET /api/user`

`Authorization: Bearer <access_token>`

#### Response

_Status: 200 OK_

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "johndoe@email.com",
  "phone_number": "9999999999",
  "country": "Country",
  "job_type": "Job Type",
  "income_range": "200000",
  "is_verified": true
}
```

---

### Frontend

The frontend is built with plain old HTML, CSS and JavaScript. The go server serves the frontend files from the `static` folder. Appropriate redirects are made in case of invalid credentials.
