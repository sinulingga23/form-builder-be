package repository

import (
	"context"

	"github.com/sinulingga23/form-builder-be/model"
)

type IMFormRepository interface {
	FindOne(ctx context.Context, id string) (model.MForm, error)
}
