package main

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrDatabaseConnection is returned when there is a database connection error.
	ErrDatabaseConnection = errors.New("database connection error")
)
