package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/sinulingga23/form-builder-be/api/repository"
	"github.com/sinulingga23/form-builder-be/model"
)

type mFormFieldRepository struct {
	db *sql.DB
}

func NewMFormFieldRepository(db *sql.DB) repository.IMFormFieldRepository {
	return &mFormFieldRepository{db: db}
}

func (repository *mFormFieldRepository) FindListFormFieldByMFormId(ctx context.Context, mFormId string) ([]*model.MFormField, error) {
	// TODO: Change column name m_from_type_id into m_field_type_id
	query := `
	select
		id, name, m_form_id, m_field_type_id, is_mandatory, ordering, placeholder, created_at, updated_at
	from
		partner.m_form_field
	where
		m_form_id = $1
	`
	rows, errQuery := repository.db.Query(query, mFormId)
	if errQuery != nil {
		return []*model.MFormField{}, errQuery
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			log.Printf("errClose: %v", errClose)
		}
	}()

	listMFormField := make([]*model.MFormField, 0)
	for rows.Next() {
		mFormField := model.MFormField{}
		errScan := rows.Scan(
			&mFormField.Id,
			&mFormField.Name,
			&mFormField.MFormId,
			&mFormField.MFieldTypeId,
			&mFormField.IsMandatory,
			&mFormField.Ordering,
			&mFormField.Placeholder,
			&mFormField.CreatedAt,
			&mFormField.UpdatedAt,
		)
		if errScan != nil {
			return []*model.MFormField{}, errScan
		}

		listMFormField = append(listMFormField, &mFormField)
	}

	if lastErr := rows.Err(); lastErr != nil {
		return []*model.MFormField{}, lastErr
	}

	return listMFormField, nil
}
