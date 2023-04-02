package model

import (
	"database/sql"
	"time"
)

type (
	MFormField struct {
		Id          string
		Name        string
		MFormId     string
		MFormTypeId string
		IsMandatory bool
		Ordering    string
		Placeholder string
		CreatedAt   time.Time
		UpdatedAt   sql.NullTime
	}
)
