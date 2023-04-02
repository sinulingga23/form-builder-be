package model

import (
	"database/sql"
	"time"
)

type (
	MForm struct {
		Id         string
		Code       string
		Name       string
		MPartnerId string
		CreatedAt  time.Time
		UpdatedAt  sql.NullTime
	}
)
