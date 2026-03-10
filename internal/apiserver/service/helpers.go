package service

import (
	"errors"

	"gpu-scheduler-platform/internal/util"
)

func IsNotFound(err error) bool {
	return errors.Is(err, util.ErrNotFound)
}

func IsInvalidArgument(err error) bool {
	return errors.Is(err, util.ErrInvalidArgument)
}

func IsUnavailable(err error) bool {
	return errors.Is(err, util.ErrUnavailable)
}
