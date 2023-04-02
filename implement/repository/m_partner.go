package repository

import (
	"context"
	"database/sql"

	"github.com/sinulingga23/form-builder-be/api/repository"
	"github.com/sinulingga23/form-builder-be/model"
)

type mPartnerRepository struct {
	db *sql.DB
}

func NewMPartnerRepository(db *sql.DB) repository.IMPartnerRepository {
	return &mPartnerRepository{db: db}
}

func (repository *mPartnerRepository) IsExistById(ctx context.Context, id string) (bool, error) {
	queryCheckMPartner := `
	select
		count(id)
	from
		partner.m_partner
	where
		id = $1
	`
	rowCheckMPartner := repository.db.QueryRow(queryCheckMPartner, id)
	countCheckMPartner := 0
	errRowCheckMpartner := rowCheckMPartner.Scan(
		&countCheckMPartner)
	if errRowCheckMpartner != nil {
		return false, errRowCheckMpartner
	}

	if lastErrRowCheckMPartner := rowCheckMPartner.Err(); lastErrRowCheckMPartner != nil {
		return false, lastErrRowCheckMPartner
	}

	if countCheckMPartner != 1 {
		return false, sql.ErrNoRows
	}

	return true, nil
}

func (repository *mPartnerRepository) FindOne(ctx context.Context, id string) (model.MPartner, error) {
	query := `
	select
		id, name, description, created_at, updated_at
	from
		partner.m_partner
	where
		id = $1
	`
	row := repository.db.QueryRow(query, id)
	mPartner := model.MPartner{}
	errScan := row.Scan(
		&mPartner.Id,
		&mPartner.Name,
		&mPartner.Description,
		&mPartner.CreatedAt,
		&mPartner.UpdatedAt,
	)
	if errScan != nil {
		return model.MPartner{}, errScan
	}

	return mPartner, nil
}
