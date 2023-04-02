package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/sinulingga23/form-builder-be/api/repository"
	"github.com/sinulingga23/form-builder-be/model"
)

type mFieldTypeRepository struct {
	db *sql.DB
}

func NewMFieldTypeRepository(db *sql.DB) repository.IMFieldTypeRepository {
	return &mFieldTypeRepository{db: db}
}

func (repository *mFieldTypeRepository) FindListMFieldTypeByIds(ctx context.Context, ids []string) ([]*model.MFieldType, error) {
	param := `(`
	lenIds := len(ids)
	for i := 0; i < lenIds; i++ {
		if i != lenIds-1 {
			param += fmt.Sprintf(`'%s',`, ids[i])
		} else {
			param += fmt.Sprintf(`'%s'`, ids[i])
		}
	}
	param += `)`

	query := `
	select
		id, name, created_at, updated_at
	from
		partner.m_field_type
	where
		id in
	`
	query += fmt.Sprintf(" %s", param)
	rows, errQuery := repository.db.Query(query)
	if errQuery != nil {
		return []*model.MFieldType{}, errQuery
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			log.Printf("errClose: %v", errClose)
		}
	}()

	listMFielType := make([]*model.MFieldType, 0)
	for rows.Next() {
		mFieldType := model.MFieldType{}
		errScan := rows.Scan(
			&mFieldType.Id,
			&mFieldType.Name,
			&mFieldType.CreatedAt,
			&mFieldType.UpdatedAt,
		)
		if errScan != nil {
			return []*model.MFieldType{}, errScan
		}

		listMFielType = append(listMFielType, &mFieldType)
	}

	if lastErr := rows.Err(); lastErr != nil {
		return []*model.MFieldType{}, lastErr
	}

	return listMFielType, nil
}
