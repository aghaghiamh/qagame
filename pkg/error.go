package userservice

import (
	"fmt"
)

const (
	// Server Errors (1000-1999)
	GeneralServerErr = 1000

	// Database Errors (2000-2999)
	GeneralDatabaseErr = 2000

	// Entity Errors (3000-3999)
	EntityNotFound = 3001

	// Auth Errors (4000-4999)
	InvalidCredentialsErr = 4002

	// Business Logic Errors (5000-5999)
	GeneralBusinessLogicError = 5000
)

type RichErr struct {
	Message string
	Code    int
}

func (err RichErr) Error() string {
	switch err.Code {
	case GeneralServerErr:
		return fmt.Sprintf("Server Error: %s", err.Message)
	case GeneralDatabaseErr:
		return fmt.Sprintf("Database Error: %s", err.Message)
	case EntityNotFound:
		return fmt.Sprintf("Not Found: %s", err.Message)
	case InvalidCredentialsErr:
		return fmt.Sprintf("Invalid phone number or password")
	case GeneralBusinessLogicError:
		return fmt.Sprintf("General Bussiness Logic Error: %s", err.Message)
	default:
		return fmt.Sprintf("Not impelemented: %s", err.Message)
	}
}
