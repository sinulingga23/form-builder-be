package model

import (
	"database/sql"
	"time"
)

type (
	MPartner struct {
		Id          string
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   sql.NullTime
	}
)
