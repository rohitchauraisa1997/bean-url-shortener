package global

import "errors"

var ErrUserNotExists = errors.New("user with above emailId does not exist")
var ErrUserAlreadyExists = errors.New("user with above username or emailId already exists")
var ErrInvalidUserCredentials = errors.New("invalid credentials passed for signing in")
