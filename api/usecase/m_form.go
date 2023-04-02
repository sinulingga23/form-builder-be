package usecase

import (
	"context"

	"github.com/sinulingga23/form-builder-be/payload"
)

type IMFormUsecase interface {
	AddFrom(ctx context.Context, createMFormRequest payload.CreateMFormRequest) payload.Response
	GetFormById(ctx context.Context, id string) payload.Response
}
