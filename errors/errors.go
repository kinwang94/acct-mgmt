package errors

import "errors"

var (
	ErrUsernameAlreadyExists  = errors.New("username already exists")
	ErrUsernameLengthTooShort = errors.New("invalid username: username should be longer than 3 characters")
	ErrUsernameLengthTooLong  = errors.New("invalid username: username should be shorter than 32 characters")

	ErrPasswordLengthTooShort = errors.New("invalid password: password should be longer than 8 characters")
	ErrPasswordLengthTooLong  = errors.New("invalid password: password should be shorter than 32 characters")
	ErrPasswordInvalidFormat  = errors.New("invalid password: password should contain at least 1 uppercase, 1 lowercase, and 1 number")

	ErrIncorrectCredential = errors.New("incorrect username or password")
	ErrTooManyAttempts     = errors.New("password verification attempts exceeded, please wait for one minute to retry")

	ErrNoAccountFound = errors.New("no account found")
)
