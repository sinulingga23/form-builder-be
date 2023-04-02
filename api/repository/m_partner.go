package repository

import (
	"context"

	"github.com/sinulingga23/form-builder-be/model"
)

type IMPartnerRepository interface {
	IsExistById(ctx context.Context, id string) (bool, error)
	FindOne(ctx context.Context, id string) (model.MPartner, error)
}
