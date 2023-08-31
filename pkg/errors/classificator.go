package errors

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	MaxYear    = 3000
	MinYear    = 1900
	MaxMonth   = 12
	MinMonth   = 1
	MaxPercent = 100
	MinPercent = 1
)

var (
	ErrInternal          = errors.New("internal server error")
	ErrUserNotFound      = errors.New("user not found")
	ErrSegmentNotFound   = errors.New("segment not found")
	ErrUserExists        = errors.New("user with this nickname already exists")
	ErrSegmentExists     = errors.New("segment with this slug already exists")
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidForm       = errors.New("invalid form")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrYearIsRequired    = errors.New("year is required")
	ErrYearIsInvalid     = errors.New("year is invalid")
	ErrMonthIsRequired   = errors.New("month is required")
	ErrMonthIsInvalid    = errors.New("month is invalid")
	ErrPercentIsInvalid  = errors.New("percent is invalid")
	ErrUntilIsInvalid    = errors.New("field until is invalid. format: YYYY-MM-DD HH:MM")
)

var HttpCodes = map[string]int{
	ErrInternal.Error():          http.StatusInternalServerError,
	ErrUserNotFound.Error():      http.StatusNotFound,
	ErrSegmentNotFound.Error():   http.StatusNotFound,
	ErrUserExists.Error():        http.StatusConflict,
	ErrSegmentExists.Error():     http.StatusConflict,
	ErrInvalidURL.Error():        http.StatusBadRequest,
	ErrInvalidForm.Error():       http.StatusBadRequest,
	ErrInvalidParameters.Error(): http.StatusBadRequest,
	ErrYearIsRequired.Error():    http.StatusBadRequest,
	ErrYearIsInvalid.Error():     http.StatusBadRequest,
	ErrMonthIsRequired.Error():   http.StatusBadRequest,
	ErrMonthIsInvalid.Error():    http.StatusBadRequest,
	ErrPercentIsInvalid.Error():  http.StatusBadRequest,
	ErrUntilIsInvalid.Error():    http.StatusBadRequest,
}

var LogLevels = map[string]logrus.Level{
	ErrInternal.Error():          logrus.ErrorLevel,
	ErrUserNotFound.Error():      logrus.WarnLevel,
	ErrSegmentNotFound.Error():   logrus.WarnLevel,
	ErrUserExists.Error():        logrus.WarnLevel,
	ErrSegmentExists.Error():     logrus.WarnLevel,
	ErrInvalidURL.Error():        logrus.WarnLevel,
	ErrInvalidForm.Error():       logrus.WarnLevel,
	ErrInvalidParameters.Error(): logrus.WarnLevel,
	ErrYearIsRequired.Error():    logrus.WarnLevel,
	ErrYearIsInvalid.Error():     logrus.WarnLevel,
	ErrMonthIsRequired.Error():   logrus.WarnLevel,
	ErrMonthIsInvalid.Error():    logrus.WarnLevel,
	ErrPercentIsInvalid.Error():  logrus.WarnLevel,
	ErrUntilIsInvalid.Error():    logrus.WarnLevel,
}

func HttpCode(err error) int {
	code, ok := HttpCodes[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}

func LogLevel(err error) logrus.Level {
	level, ok := LogLevels[err.Error()]
	if !ok {
		return logrus.ErrorLevel
	}

	return level
}
