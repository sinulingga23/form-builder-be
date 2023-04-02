package repository

import (
	"context"
	"database/sql"

	"github.com/sinulingga23/form-builder-be/api/repository"
	"github.com/sinulingga23/form-builder-be/model"
)

type mFormRepository struct {
	db *sql.DB
}

func NewMFormRepository(db *sql.DB) repository.IMFormRepository {
	return &mFormRepository{db: db}
}

func (repository *mFormRepository) FindOne(ctx context.Context, id string) (model.MForm, error) {
	query := `
	select
		id, code, name, m_partner_id, created_at, updated_at
	from
		partner.m_form
	where
		id = $1
	`
	row := repository.db.QueryRow(query, id)

	mForm := model.MForm{}
	errScan := row.Scan(
		&mForm.Id,
		&mForm.Code,
		&mForm.Name,
		&mForm.MPartnerId,
		&mForm.CreatedAt,
		&mForm.UpdatedAt,
	)
	if errScan != nil {
		return model.MForm{}, errScan
	}

	if lastErr := row.Err(); lastErr != nil {
		return model.MForm{}, lastErr
	}

	return mForm, nil
}
