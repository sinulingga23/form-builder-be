package repository

import (
	"context"

	"github.com/sinulingga23/form-builder-be/model"
)

type IMFormFieldChildsRepository interface {
	FindListMFormFieldChildsByMFormFieldByIds(ctx context.Context, mFormFieldIds []string) ([]*model.MFormFieldChilds, error)
}
