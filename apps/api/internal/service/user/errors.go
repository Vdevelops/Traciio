package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRoleNotFound      = errors.New("role not found")
	ErrGroupNotFound     = errors.New("group not found")
	ErrBrickNotFound     = errors.New("brick not found")
	ErrIncorrectPassword = errors.New("current password is incorrect")
)
