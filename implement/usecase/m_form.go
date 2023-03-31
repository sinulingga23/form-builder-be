package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sinulingga23/form-builder-be/api/usecase"
	"github.com/sinulingga23/form-builder-be/define"
	"github.com/sinulingga23/form-builder-be/payload"
)

type mFormUsecase struct {
	db *sql.DB
}

func NewMFormUsecase(db *sql.DB) usecase.IMFormUsecase {
	return &mFormUsecase{db: db}
}

func (usecase *mFormUsecase) AddFrom(ctx context.Context, createMFormRequest payload.CreateMFormRequest) payload.Response {
	response := payload.Response{
		StatusCode: http.StatusOK,
		Message:    "Success add a new form.",
	}

	if strings.Trim(createMFormRequest.MFormName, " ") == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = define.ErrMFormNameEmpty.Error()
		return response
	}

	if strings.Trim(createMFormRequest.MPartnerId, " ") == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = define.ErrMPartnerIdEmpty.Error()
		return response
	}

	_, errParsePartnerId := uuid.Parse(createMFormRequest.MPartnerId)
	if errParsePartnerId != nil {
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMPartnerNotFound.Error()
		return response
	}

	// validate all the MFieldTypeId
	mapFieldTypes := make([]string, 0)
	for _, mFormField := range createMFormRequest.MFormFields {
		if strings.Trim(mFormField.MFormFieldName, " ") == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = define.ErrMFormFieldNameEmpty.Error()
			return response
		}

		if strings.Trim(mFormField.MFieldTypeId, " ") == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = define.ErrMFieldTypeIdEmpty.Error()
			return response
		}

		_, errParseMFieldTypeId := uuid.Parse(mFormField.MFieldTypeId)
		if errParseMFieldTypeId != nil {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()
			return response
		}

		mapFieldTypes = append(mapFieldTypes, mFormField.MFieldTypeId)
	}

	tx, errBegin := usecase.db.Begin()
	if errBegin != nil {
		response.StatusCode = http.StatusInternalServerError
		return response
	}

	_ = tx

	return response
}
