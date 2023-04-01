package repository

import "context"

type IMPartnerRepository interface {
	IsExistById(ctx context.Context, id string) (bool, error)
}
