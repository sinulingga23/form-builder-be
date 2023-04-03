package model

import (
	"database/sql"
	"time"
)

type (
	MFormFieldChilds struct {
		Id           string
		Name         string
		MFormFieldId string
		CreatedAt    time.Time
		UpdatedAt    sql.NullTime
	}
)
