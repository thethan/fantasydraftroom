package mysql

import "time"

type NullTime struct {
	time.Time
	Valid bool
}