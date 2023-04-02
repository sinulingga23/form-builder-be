package repository

import (
	"context"

	"github.com/sinulingga23/form-builder-be/model"
)

type IMFieldTypeRepository interface {
	FindListMFieldTypeByIds(ctx context.Context, ids []string) ([]*model.MFieldType, error)
}
