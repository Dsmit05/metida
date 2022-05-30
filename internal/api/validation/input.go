package validation

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

const (
	minPassLength     = 8
	maxPassLength     = 32
	minUserNameLength = 4
	maxUserNameLength = 24
)

// IsEmailValid checks the correct user email.
func IsEmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// IsPasswordValid check password for next rules:
// at least 8 characters
// maximum of 32 characters
// at least 1 number
// at least 1 upper case
// at least 1 lower case
// at least 1 special character.
func IsPasswordValid(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}

	appendError := func(err string) {
		if strings.TrimSpace(errorString) != "" {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}

	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if errorString != "" {
		return fmt.Errorf(errorString)
	}
	return nil
}

// IsUserNameValid checks the correct user name.
func IsUserNameValid(userName string) error {
	userLen := len(userName)

	if !(minUserNameLength <= userLen && userLen <= maxUserNameLength) {
		return fmt.Errorf("name length must be between %d to %d characters long", minUserNameLength, maxUserNameLength)
	}

	return nil
}
