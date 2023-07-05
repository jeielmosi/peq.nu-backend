package repository_helpers

import (
	"errors"
	"time"
)

const timestamp1e8Pattern = "2006-01-02T15:04:05.99999999"

func NewTimeFromTimestamp1e8(timestamp *string) (*time.Time, error) {
	if timestamp == nil {
		return nil, errors.New("Timestamp is nil")
	}
	res, err := time.Parse(timestamp1e8Pattern, *timestamp)

	return &res, err
}

func TimeToTimestamp1e8(t time.Time) string {
	return t.UTC().Format(timestamp1e8Pattern)
}

func NowTimestamp1e8() string {
	return time.Now().UTC().Format(timestamp1e8Pattern)
}
