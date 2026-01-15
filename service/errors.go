package service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRowsAffected     = errors.New("no rows were affected")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("record already exists")
	ErrNotFound           = errors.New("record not found")
)

func wrapGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		err = ErrAlreadyExists
	}
	return err
}
