package ssmodels

import (
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02T15:04:05"
)

type MyTime struct {
	time.Time
}

func (m *MyTime) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(dateFormat, strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	))

	if err != nil {
		return err
	}

	m.Time = t

	return nil
}
