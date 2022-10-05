package postgresql

import (
	"errors"
	ricardoerr "gitlab.com/ricardo-public/errors/pkg/errors"
	"gorm.io/gorm"
)

func notFoundOrElseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ricardoerr.New(ricardoerr.ErrNotFound, "record not found")
	}

	return ricardoerr.New(ricardoerr.ErrDatabaseError, err.Error())
}
