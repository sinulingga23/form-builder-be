package model

import (
	"database/sql"
	"time"
)

type (
	MFieldType struct {
		Id        string
		Name      string
		CreatedAt time.Time
		UpdatedAt sql.NullTime
	}
)
