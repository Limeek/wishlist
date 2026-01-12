package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserInvalidPassword = errors.New("invalid password")
var ErrUserNotAuthorized = errors.New("user not authorized")
var ErrRefreshTokenIsNotValid = errors.New("refresh token is not valid")
var ErrEmailInvalidFormat = errors.New("invalid email format")
var ErrPasswordInvalidFormat = errors.New("invalid password format")
var ErrPhoneInvalidFormat = errors.New("invalid phone format")
var ErrBirthDateInvalidFormat = errors.New("invalid birth date format")
