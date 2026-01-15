package auth

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("record already exists")
	ErrUserNotFound       = errors.New("user not found")
)

func wrapGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotFound
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrUserAlreadyExists
	}
	return err
}
