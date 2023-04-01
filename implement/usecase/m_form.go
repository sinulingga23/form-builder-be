package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sinulingga23/form-builder-be/api/usecase"
	"github.com/sinulingga23/form-builder-be/define"
	"github.com/sinulingga23/form-builder-be/model"
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
	mFieldTypeIds := make([]string, 0)
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

		mFieldTypeIds = append(mFieldTypeIds, mFormField.MFieldTypeId)
	}

	tx, errBegin := usecase.db.Begin()
	if errBegin != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			response.StatusCode = http.StatusInternalServerError
			return response
		}
		response.StatusCode = http.StatusInternalServerError
		return response
	}

	// check m_parther by mPartnerId, ensure it's exists
	queryCheckMPartner := `
	select
		count(id)
	from
		partner.m_partner
	where
		id = $1
	`
	rowCheckMPartner := tx.QueryRow(queryCheckMPartner, createMFormRequest.MPartnerId)
	countCheckMPartner := 0
	errRowCheckMpartner := rowCheckMPartner.Scan(&countCheckMPartner)
	if errRowCheckMpartner != nil {
		log.Printf("errRowCheckMpartner: %v", errRowCheckMpartner)
		if errRollback := tx.Rollback(); errRollback != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}

		if errors.Is(errRowCheckMpartner, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMPartnerNotFound.Error()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()
		return response
	}

	if lastErrRowCheckMPartner := rowCheckMPartner.Err(); lastErrRowCheckMPartner != nil {
		log.Printf("lastErrRowCheckMPartner: %v", lastErrRowCheckMPartner)
		if errRollback := tx.Rollback(); errRollback != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()
		return response
	}

	if countCheckMPartner != 1 {
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMPartnerNotFound.Error()
		return response
	}

	paramMFieldTypeIds := `(`
	lenMFieldTypeIds := len(mFieldTypeIds)
	for i := 0; i < lenMFieldTypeIds; i++ {
		if i != lenMFieldTypeIds-1 {
			paramMFieldTypeIds += fmt.Sprintf(`'%s',`, mFieldTypeIds[i])
		} else {
			paramMFieldTypeIds += fmt.Sprintf(`'%s'`, mFieldTypeIds[i])
		}
	}
	paramMFieldTypeIds += `)`

	// get list m_field_type based on paramMFieldTypeIds
	queryMFieldTypeByIds := `
	select
		id, name
	from
		partner.m_field_type
	where
		id in
	`
	queryMFieldTypeByIds += fmt.Sprintf(` %s`, paramMFieldTypeIds)
	rowsQueryMFieldTypeByIds, errQueryMFieldTypeByIds := tx.Query(queryMFieldTypeByIds)
	if errQueryMFieldTypeByIds != nil {
		log.Printf("errQueryMFieldTypeByIds: %v", errQueryMFieldTypeByIds)
		if errors.Is(errQueryMFieldTypeByIds, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()
			return response
		}

		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()
		return response
	}
	defer func() {
		if errClose := rowsQueryMFieldTypeByIds.Close(); errClose != nil {
			log.Printf("errClose: %v", errClose)
		}
	}()

	listMFieldType := make([]model.MFieldType, 0)
	for rowsQueryMFieldTypeByIds.Next() {
		mFieldType := model.MFieldType{}
		errScan := rowsQueryMFieldTypeByIds.Scan(
			&mFieldType.Id,
			&mFieldType.Name,
		)
		if errScan != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Printf("errRollback: %v", errRollback)
				response.StatusCode = http.StatusInternalServerError
				response.Message = define.ErrInternalServerError.Error()
				return response
			}

			log.Printf("errScan: %v", errScan)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrQueryData.Error()
			return response
		}

		listMFieldType = append(listMFieldType, mFieldType)
	}
	if lastErrQueryMFieldTypeByIds := rowsQueryMFieldTypeByIds.Err(); lastErrQueryMFieldTypeByIds != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()
		return response
	}

	// compare length of listMFieldType with mFormFields, it should same.
	// if not same, maybe is there some m_field_type not exists
	if len(listMFieldType) != len(createMFormRequest.MFormFields) {
		log.Printf("len(listMFieldType): %v, len(createMFormRequest.MFormFields): %v", len(listMFieldType), len(createMFormRequest.MFormFields))
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}

		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMFieldTypeNotFound.Error()
		return response
	}

	// convert listMFieldType into an map, the
	mapMFieldType := make(map[string]string)
	for _, mFieldType := range listMFieldType {
		mapMFieldType[mFieldType.Id] = mFieldType.Name
	}

	// insert row into partner.m_form table
	queryInsertMForm := `
	insert into partner.m_form
		(id, code, name, m_partner_id, created_at)
	values
		($1, $2, $3, $4, $5)
	`
	mFormId := uuid.NewString()
	_, errInsertMForm := tx.Exec(queryInsertMForm,
		mFormId,
		"code",
		createMFormRequest.MFormName,
		createMFormRequest.MPartnerId,
		time.Now())
	if errInsertMForm != nil {
		log.Printf("errInsertMForm: %v", errInsertMForm)
		if errRollback := tx.Rollback(); errRollback != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()
		return response
	}

	paramInsertListMFormField := ``
	mapMFormFieldIdWithMFieldType := make(map[string]string)
	lenMFormFields := len(createMFormRequest.MFormFields)
	for i := 0; i < lenMFormFields; i++ {
		mFormFieldId := uuid.NewString()
		mFormField := createMFormRequest.MFormFields[i]
		if i != lenMFieldTypeIds-1 {
			paramInsertListMFormField += fmt.Sprintf(`('%s','%s','%s', '%s', '%v', '%v', '%s', '%v'),`,
				mFormFieldId,
				mFormField.MFormFieldName,
				mFormId,
				mFormField.MFieldTypeId,
				mFormField.MFormFieldIsMandatory,
				mFormField.MFormFieldOrdering,
				mFormField.MFormFieldPlaceholder,
				time.Now())
		} else {
			paramInsertListMFormField += fmt.Sprintf(`('%s','%s','%s', '%s', '%v', '%v', '%s', '%v')`,
				mFormFieldId,
				mFormField.MFormFieldName,
				mFormId,
				mFormField.MFieldTypeId,
				mFormField.MFormFieldIsMandatory,
				mFormField.MFormFieldOrdering,
				mFormField.MFormFieldPlaceholder,
				time.Now())
		}

		mapMFormFieldIdWithMFieldType[mFormFieldId] = mFormField.MFieldTypeId
	}
	queryInsertListMFormField := `
	insert into partner.m_form_field
		(id, name, m_form_id, m_form_type_id, is_mandatory, ordering, created_at)
	values
	`
	queryInsertListMFormField += fmt.Sprintf(" %s", paramInsertListMFormField)
	_, errInsertListMFormField := tx.Exec(queryInsertListMFormField)
	if errInsertListMFormField != nil {
		log.Printf("errInsertListMFormField: %v", errInsertListMFormField)
		if errRollback := tx.Rollback(); errRollback != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()
			return response
		}
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()
		return response
	}

	// for the m_field_type which have child, insert the childs
	// into partner.m_form_field_childs table.
	for mFormFieldId, mFieldTypeId := range mapMFormFieldIdWithMFieldType {
		mFieldTypeName, oKMFieldTypeName := mapMFieldType[mFieldTypeId]
		if !oKMFieldTypeName {
			if errRollback := tx.Rollback(); errRollback != nil {
				response.StatusCode = http.StatusInternalServerError
				response.Message = define.ErrInternalServerError.Error()
				return response
			}
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()
			return response
		}

		if mFieldTypeName == define.HAS_CHILD_M_FIELD_TYPE_DROPDOWN ||
			mFieldTypeName == define.HAS_CHILD_M_FIELD_TYPE_RADIO_BUTTON {
			// ennsure the child not empty
		}

		_ = mFormFieldId
	}

	return response
}
