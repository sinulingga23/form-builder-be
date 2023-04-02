package repository

import (
	"context"

	"github.com/sinulingga23/form-builder-be/model"
)

type IMFormFieldRepository interface {
	FindListFormFieldByMFormId(ctx context.Context, mFormId string) ([]*model.MFormField, error)
}
