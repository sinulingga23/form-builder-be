package repository

import (
	"context"
	"database/sql"

	"github.com/sinulingga23/form-builder-be/api/repository"
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
