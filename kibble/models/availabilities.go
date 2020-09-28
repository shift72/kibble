package models

import "time"

// Period stores a range
type Period struct {
	From *time.Time
	To   *time.Time
}
