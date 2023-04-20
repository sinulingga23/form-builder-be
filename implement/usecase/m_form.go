package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sinulingga23/form-builder-be/api/repository"
	"github.com/sinulingga23/form-builder-be/api/usecase"
	"github.com/sinulingga23/form-builder-be/cache"
	"github.com/sinulingga23/form-builder-be/define"
	"github.com/sinulingga23/form-builder-be/model"
	"github.com/sinulingga23/form-builder-be/monitoring"
	"github.com/sinulingga23/form-builder-be/payload"
)

type mFormUsecase struct {
	db                         *sql.DB
	mPartnerRepository         repository.IMPartnerRepository
	mFieldTypeRepository       repository.IMFieldTypeRepository
	mFormRepository            repository.IMFormRepository
	mFormFieldRepository       repository.IMFormFieldRepository
	mFormFieldChildsRepository repository.IMFormFieldChildsRepository
	metric                     *monitoring.Metric
}

func NewMFormUsecase(
	db *sql.DB,
	mPartnerRepository repository.IMPartnerRepository,
	mFieldRepository repository.IMFieldTypeRepository,
	mFormRepository repository.IMFormRepository,
	mFormFieldRepository repository.IMFormFieldRepository,
	mFormFieldChildsRepository repository.IMFormFieldChildsRepository,
	metric *monitoring.Metric,
) usecase.IMFormUsecase {
	return &mFormUsecase{
		db:                         db,
		mPartnerRepository:         mPartnerRepository,
		mFieldTypeRepository:       mFieldRepository,
		mFormRepository:            mFormRepository,
		mFormFieldRepository:       mFormFieldRepository,
		mFormFieldChildsRepository: mFormFieldChildsRepository,
		metric:                     metric,
	}
}

func (usecase *mFormUsecase) AddFrom(ctx context.Context, createMFormRequest payload.CreateMFormRequest) payload.Response {
	response := payload.Response{
		StatusCode: http.StatusOK,
		Message:    "Success add a new form.",
	}

	serviceName := "m_form_service:add_form"
	now := time.Now()

	if strings.Trim(createMFormRequest.MFormName, " ") == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = define.ErrMFormNameEmpty.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
		return response
	}

	if strings.Trim(createMFormRequest.MPartnerId, " ") == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = define.ErrMPartnerIdEmpty.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
		return response
	}

	_, errParsePartnerId := uuid.Parse(createMFormRequest.MPartnerId)
	if errParsePartnerId != nil {
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMPartnerNotFound.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusNotFound), &now)
		return response
	}

	// validate all the MFieldTypeId
	mFieldTypeIds := make([]string, 0)
	mapOrdering := make(map[int]int)
	for _, mFormField := range createMFormRequest.MFormFields {
		if strings.Trim(mFormField.MFormFieldName, " ") == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = define.ErrMFormFieldNameEmpty.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
			return response
		}

		if strings.Trim(mFormField.MFieldTypeId, " ") == "" {
			response.StatusCode = http.StatusBadRequest
			response.Message = define.ErrMFieldTypeIdEmpty.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
			return response
		}

		_, errParseMFieldTypeId := uuid.Parse(mFormField.MFieldTypeId)
		if errParseMFieldTypeId != nil {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusNotFound), &now)
			return response
		}

		mFieldTypeIds = append(mFieldTypeIds, mFormField.MFieldTypeId)

		// egde case, ensure ordering not duplicate from client
		mapOrdering[mFormField.MFormFieldOrdering] += 1
		if mapOrdering[mFormField.MFormFieldOrdering] > 1 {
			response.StatusCode = http.StatusBadRequest
			response.Message = define.ErrOrderingCantDuplicet.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
			return response
		}
	}

	// check m_parther by mPartnerId, ensure it's exists
	isMPartnerExists, errIsExistsById := usecase.mPartnerRepository.IsExistById(ctx, createMFormRequest.MPartnerId)
	if errIsExistsById != nil {
		log.Printf("errIsExistsById: %v", errIsExistsById)
		if errors.Is(errIsExistsById, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMPartnerNotFound.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusNotFound), &now)
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrInternalServerError.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
		return response
	}

	if !isMPartnerExists {
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMPartnerNotFound.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusNotFound), &now)
		return response
	}

	tx, errBegin := usecase.db.Begin()
	if errBegin != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
		}
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrInternalServerError.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
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

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusNotFound), &now)
			return response
		}

		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
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

				go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
				return response
			}

			log.Printf("errScan: %v", errScan)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrQueryData.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
			return response
		}

		listMFieldType = append(listMFieldType, mFieldType)
	}
	if lastErrQueryMFieldTypeByIds := rowsQueryMFieldTypeByIds.Err(); lastErrQueryMFieldTypeByIds != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
			return response
		}
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
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
		define.GenerateRandomString(define.LENGTH_CODE_FOR_M_FORM),
		createMFormRequest.MFormName,
		createMFormRequest.MPartnerId,
		time.Now())
	if errInsertMForm != nil {
		log.Printf("errInsertMForm: %v", errInsertMForm)
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("errRollback: %v", errRollback)
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
		return response
	}

	paramInsertListMFormField := ``
	mapMFormFieldIdWithMFieldType := make(map[string]string)
	mapMFromFieldIdWithChilds := make(map[string][]payload.MFormFieldChild)
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
				time.Now().Format(time.RFC3339))
		} else {
			paramInsertListMFormField += fmt.Sprintf(`('%s','%s','%s', '%s', '%v', '%v', '%s', '%v')`,
				mFormFieldId,
				mFormField.MFormFieldName,
				mFormId,
				mFormField.MFieldTypeId,
				mFormField.MFormFieldIsMandatory,
				mFormField.MFormFieldOrdering,
				mFormField.MFormFieldPlaceholder,
				time.Now().Format(time.RFC3339))
		}

		mapMFormFieldIdWithMFieldType[mFormFieldId] = mFormField.MFieldTypeId
		mapMFromFieldIdWithChilds[mFormFieldId] = mFormField.MFormFieldChilds
	}
	queryInsertListMFormField := `
	insert into partner.m_form_field
		(id, name, m_form_id, m_field_type_id, is_mandatory, ordering, placeholder, created_at)
	values
	`
	queryInsertListMFormField += fmt.Sprintf(" %s", paramInsertListMFormField)
	_, errInsertListMFormField := tx.Exec(queryInsertListMFormField)
	if errInsertListMFormField != nil {
		log.Printf("errInsertListMFormField: %v", errInsertListMFormField)
		if errRollback := tx.Rollback(); errRollback != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Message = define.ErrInternalServerError.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
			return response
		}
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
		return response
	}

	// for the m_field_type which have child, insert the childs
	// into partner.m_form_field_childs table.
	for mFormFieldId, mFieldTypeId := range mapMFormFieldIdWithMFieldType {
		mFieldTypeName, oKMFieldTypeName := mapMFieldType[mFieldTypeId]
		if !oKMFieldTypeName {
			log.Printf("oKMFieldTypeName: %v", oKMFieldTypeName)
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Printf("errRollback: %v", errRollback)
				response.StatusCode = http.StatusInternalServerError
				response.Message = define.ErrInternalServerError.Error()

				go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
				return response
			}
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()

			go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusNotFound), &now)
			return response
		}

		if mFieldTypeName == define.M_FIELD_TYPE_DROPDOWN ||
			mFieldTypeName == define.M_FIELD_TYPE_RADIO_BUTTON ||
			mFieldTypeName == define.M_FIELD_TYPE_CHECKBOX {
			// ennsure the child not empty
			mFormFieldChilds, oKMFormFieldChilds := mapMFromFieldIdWithChilds[mFormFieldId]
			log.Printf("mFormFieldChilds: %v, oKMFormFieldChilds: %v", mFormFieldChilds, oKMFormFieldChilds)
			if !oKMFormFieldChilds {
				if errRollback := tx.Rollback(); errRollback != nil {
					log.Printf("errRollback: %v", errRollback)
					response.StatusCode = http.StatusInternalServerError
					response.Message = define.ErrInternalServerError.Error()

					go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
					return response
				}
				response.StatusCode = http.StatusBadRequest
				response.Message = define.ErrMFormFieldsEmpty.Error()

				go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
				return response
			}

			if len(mFormFieldChilds) == 0 {
				if errRollback := tx.Rollback(); errRollback != nil {
					log.Printf("errRollback: %v", errRollback)
					response.StatusCode = http.StatusInternalServerError
					response.Message = define.ErrInternalServerError.Error()

					go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
					return response
				}

				response.StatusCode = http.StatusBadRequest
				response.Message = define.ErrMFormFieldsEmpty.Error()

				go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusBadRequest), &now)
				return response
			}

			paramInsertMFormFieldChilds := ``
			lenMFormFieldChilds := len(mFormFieldChilds)
			for i := 0; i < lenMFormFieldChilds; i++ {
				mFormFieldChildsId := uuid.NewString()
				if i != lenMFormFieldChilds-1 {
					paramInsertMFormFieldChilds += fmt.Sprintf(`('%s','%s','%s', '%s'),`,
						mFormFieldChildsId,
						mFormFieldChilds[i].MFormFieldChildName,
						mFormFieldId,
						time.Now().Format(time.RFC3339),
					)
				} else {
					paramInsertMFormFieldChilds += fmt.Sprintf(`('%s','%s','%s', '%s')`,
						mFormFieldChildsId,
						mFormFieldChilds[i].MFormFieldChildName,
						mFormFieldId,
						time.Now().Format(time.RFC3339),
					)
				}
			}

			queryInsertMFormFieldChilds := `
			insert into	partner.m_form_field_childs
				(id, name, m_form_field_id, created_at)
			values
			`
			queryInsertMFormFieldChilds += fmt.Sprintf(` %s`, paramInsertMFormFieldChilds)
			_, errInsertMFormFieldChilds := tx.Exec(queryInsertMFormFieldChilds)
			if errInsertMFormFieldChilds != nil {
				log.Printf("errInsertMFormFieldChilds: %v", errInsertMFormFieldChilds)
				if errRollback := tx.Rollback(); errRollback != nil {
					log.Printf("errRollback: %v", errRollback)
					response.StatusCode = http.StatusInternalServerError
					response.Message = define.ErrInternalServerError.Error()

					go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
					return response
				}

				response.StatusCode = http.StatusInternalServerError
				response.Message = define.ErrQueryData.Error()

				go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
				return response
			}
		}
	}

	// commit transaction
	if errCommit := tx.Commit(); errCommit != nil {
		log.Printf("errCommit: %v", errCommit)
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusInternalServerError), &now)
		return response
	}

	go usecase.metric.CaptureMetrics(serviceName, http.MethodPost, strconv.Itoa(http.StatusOK), &now)
	return response
}

func (usecase *mFormUsecase) GetFormById(ctx context.Context, id string) payload.Response {
	response := payload.Response{
		StatusCode: http.StatusOK,
		Message:    "Success to get the form.",
	}

	serviceName := "m_form_service:get_form_by_id"
	now := time.Now()
	httpMethod := http.MethodGet

	cacheResponse, errGetValue := cache.GetValue(ctx, id)
	if errGetValue == nil {
		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusOK), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return cacheResponse.(payload.Response)
	}
	log.Printf("errGetValue: %v", errGetValue)

	if strings.Trim(id, " ") == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = define.ErrIdEmpty.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusBadRequest), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()

		return response
	}

	_, errPase := uuid.Parse(id)
	if errPase != nil {
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMFormNotFound.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()

		return response
	}

	mForm, errFindOneMForm := usecase.mFormRepository.FindOne(ctx, id)
	if errFindOneMForm != nil {
		log.Printf("errFindOneMForm: %v", errFindOneMForm)
		if errors.Is(errFindOneMForm, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFormNotFound.Error()

			go func() {
				errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
				if errCaptureMetrics != nil {
					log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
				}
			}()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusInternalServerError), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return response
	}

	mPartner, errFindOneMPartner := usecase.mPartnerRepository.FindOne(ctx, mForm.MPartnerId)
	if errFindOneMPartner != nil {
		log.Printf("errFindOneMPartner: %v", errFindOneMPartner)
		if errors.Is(errFindOneMPartner, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMPartnerNotFound.Error()

			go func() {
				errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
				if errCaptureMetrics != nil {
					log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
				}
			}()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrInternalServerError.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusInternalServerError), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()

		return response
	}

	listMFormField, errFindListFormFieldByMFormId := usecase.mFormFieldRepository.FindListFormFieldByMFormId(ctx, mForm.Id)
	if errFindListFormFieldByMFormId != nil {
		log.Printf("FindListFormFieldByMFormId: %v", errFindListFormFieldByMFormId)
		if errors.Is(errFindListFormFieldByMFormId, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFormFieldNotFound.Error()

			go func() {
				errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
				if errCaptureMetrics != nil {
					log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
				}
			}()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusInternalServerError), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return response
	}

	if len(listMFormField) == 0 {
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMFormFieldNotFound.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return response
	}

	mFieldTypeIds := make([]string, 0)
	mFormFieldIds := make([]string, 0)
	for _, mFormField := range listMFormField {
		mFieldTypeIds = append(mFieldTypeIds, mFormField.MFieldTypeId)
		mFormFieldIds = append(mFormFieldIds, mFormField.Id)
	}

	listMFieldType, errFindListMFieldTypeByIds := usecase.mFieldTypeRepository.FindListMFieldTypeByIds(ctx, mFieldTypeIds)
	if errFindListMFieldTypeByIds != nil {
		log.Printf("errFindListMFieldTypeByIds: %v", errFindListMFieldTypeByIds)
		if errors.Is(errFindListMFieldTypeByIds, sql.ErrNoRows) {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()

			go func() {
				errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
				if errCaptureMetrics != nil {
					log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
				}
			}()
			return response
		}

		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusInternalServerError), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return response
	}

	if len(listMFieldType) == 0 {
		log.Printf("listMFieldType: %v", listMFieldType)
		response.StatusCode = http.StatusNotFound
		response.Message = define.ErrMFieldTypeNotFound.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return response
	}

	listMFormFieldResponse := make([]payload.MFormFieldResponse, 0)
	mapFieldType := make(map[string]string)
	for _, mFiedlType := range listMFieldType {
		mapFieldType[mFiedlType.Id] = mFiedlType.Name
	}

	listMFormFieldChilds, errFindListMFormChildsByMFormFieldById := usecase.mFormFieldChildsRepository.FindListMFormFieldChildsByMFormFieldByIds(ctx, mFormFieldIds)
	if errFindListMFormChildsByMFormFieldById != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = define.ErrQueryData.Error()

		go func() {
			errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusInternalServerError), &now)
			if errCaptureMetrics != nil {
				log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
			}
		}()
		return response
	}

	// ('d35ffd80-c932-4768-8fb5-77c7212058f2', '6ecd65d6-50f0-4831-9943-b4a290d09340', '357ae4c8-818b-4939-8c19-23e155907cb8', '883b222e-a1f8-4437-9460-fccb8221d390','76157af8-f6de-4b1c-a545-da3a92f61b63')
	mapFormFieldIdWithListMFormFieldChilds := make(map[string][]payload.MFormFieldChildResponse)
	lenMFormFieldChilds := len(listMFormFieldChilds)
	for i := 0; i < lenMFormFieldChilds; i++ {
		listMFormFieldChildResponse := make([]payload.MFormFieldChildResponse, 0)
		for j := 0; j < lenMFormFieldChilds; j++ {
			if listMFormFieldChilds[i].MFormFieldId == listMFormFieldChilds[j].MFormFieldId {
				listMFormFieldChildResponse = append(listMFormFieldChildResponse, payload.MFormFieldChildResponse{
					MFormFieldChildId:   listMFormFieldChilds[j].Id,
					MFormFieldChildName: listMFormFieldChilds[j].Name,
				})
			}
		}

		mapFormFieldIdWithListMFormFieldChilds[listMFormFieldChilds[i].MFormFieldId] = listMFormFieldChildResponse
	}

	for _, mFormField := range listMFormField {
		mFormFieldResponse := payload.MFormFieldResponse{}

		mFormFieldResponse.MFormFieldId = mFormField.Id
		mFormFieldResponse.MFormFieldName = mFormField.Name
		mFormFieldResponse.MFormFieldIsMandatory = mFormField.IsMandatory
		mFormFieldResponse.MFormFieldOrdering = mFormField.Ordering
		mFormFieldResponse.MFormFieldPlaceholder = mFormField.Placeholder

		mFormFieldResponse.MFieldTypeId = mFormField.MFieldTypeId
		mFieldTypeName, okMFieldTypeName := mapFieldType[mFormField.MFieldTypeId]
		if !okMFieldTypeName {
			response.StatusCode = http.StatusNotFound
			response.Message = define.ErrMFieldTypeNotFound.Error()

			go func() {
				errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusNotFound), &now)
				if errCaptureMetrics != nil {
					log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
				}
			}()
			return response
		}
		mFormFieldResponse.MFieldTypeName = mFieldTypeName

		listMFormFieldChilds, oKListMFormFieldChilds := mapFormFieldIdWithListMFormFieldChilds[mFormField.Id]
		if oKListMFormFieldChilds {
			mFormFieldResponse.MFormFieldChildsResponse = listMFormFieldChilds
		} else {
			mFormFieldResponse.MFormFieldChildsResponse = []payload.MFormFieldChildResponse{}
		}

		listMFormFieldResponse = append(listMFormFieldResponse, mFormFieldResponse)
	}

	updateAt := time.Time{}
	if mForm.UpdatedAt.Valid {
		updateAt = mForm.UpdatedAt.Time
	}

	mFormDetailResponse := payload.MFormDetailResponse{
		MFormId:             mForm.Id,
		MFormName:           mForm.Name,
		MFormCode:           mForm.Code,
		MPartnerId:          mPartner.Id,
		MPartnerName:        mPartner.Name,
		MFormFieldsResponse: listMFormFieldResponse,
		MFormCreatedAt:      mForm.CreatedAt,
		MFormUpdatedAt:      updateAt,
	}

	response.Data = mFormDetailResponse
	go func() {
		if errSetValue := cache.SetValue(ctx, id, response); errSetValue != nil {
			log.Printf("errSetValue: %v", errSetValue)
		}

		errCaptureMetrics := usecase.metric.CaptureMetrics(serviceName, httpMethod, strconv.Itoa(http.StatusOK), &now)
		if errCaptureMetrics != nil {
			log.Printf("errCaptureMetrics: %v", errCaptureMetrics)
		}
	}()
	return response
}
