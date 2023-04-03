package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/sinulingga23/form-builder-be/api/repository"
	"github.com/sinulingga23/form-builder-be/model"
)

type mFormFieldChildsRepository struct {
	db *sql.DB
}

func NewMFormFieldChildsRepository(db *sql.DB) repository.IMFormFieldChildsRepository {
	return &mFormFieldChildsRepository{db: db}
}

func (repository *mFormFieldChildsRepository) FindListMFormFieldChildsByMFormFieldByIds(ctx context.Context, mFormFieldIds []string) ([]*model.MFormFieldChilds, error) {
	query := `
	select
		id, name, m_form_field_id, created_at, updated_at
	from
		partner.m_form_field_childs
	where
		m_form_field_id in
	`

	param := `(`
	lenMFormFields := len(mFormFieldIds)
	for i := 0; i < lenMFormFields; i++ {
		if i != lenMFormFields-1 {
			param += fmt.Sprintf(`'%s',`, mFormFieldIds[i])
		} else {
			param += fmt.Sprintf(`'%s'`, mFormFieldIds[i])
		}
	}
	param += `)`

	query += fmt.Sprintf(" %s", param)

	rows, errQuery := repository.db.Query(query)
	if errQuery != nil {
		return []*model.MFormFieldChilds{}, errQuery
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			log.Printf("errClose: %v", errClose)
		}
	}()

	listMFormFieldChilds := make([]*model.MFormFieldChilds, 0)
	for rows.Next() {
		mFormFieldChilds := model.MFormFieldChilds{}
		errScan := rows.Scan(
			&mFormFieldChilds.Id,
			&mFormFieldChilds.Name,
			&mFormFieldChilds.MFormFieldId,
			&mFormFieldChilds.CreatedAt,
			&mFormFieldChilds.UpdatedAt,
		)
		if errScan != nil {
			return []*model.MFormFieldChilds{}, errScan
		}

		listMFormFieldChilds = append(listMFormFieldChilds, &mFormFieldChilds)
	}

	if lastErr := rows.Err(); lastErr != nil {
		return []*model.MFormFieldChilds{}, lastErr
	}

	return listMFormFieldChilds, nil
}
