package payload

import "time"

type (
	CreateMFormRequest struct {
		MPartnerId  string       `json:"mPartnerId"`
		MFormName   string       `json:"mFormName"`
		MFormFields []MFormField `json:"mFormFields"`
	}

	MFormField struct {
		MFieldTypeId          string            `json:"mFieldTypeId"`
		MFormFieldName        string            `json:"mFormFieldName"`
		MFormFieldIsMandatory bool              `json:"mFormFieldIsMandatory"`
		MFormFieldOrdering    int               `json:"mFormFieldOrdering"`
		MFormFieldPlaceholder string            `json:"mFormFieldPlaceholder"`
		MFormFieldChilds      []MFormFieldChild `json:"mFormFieldChilds"`
	}

	MFormFieldChild struct {
		MFormFieldChildName string `json:"mFormFieldChildName"`
	}

	MFormDetailResponse struct {
		MFormId             string               `json:"mFormId"`
		MFormName           string               `json:"mFormName"`
		MFormCode           string               `json:"mFormCode"`
		MPartnerId          string               `json:"mPartnerId"`
		MPartnerName        string               `json:"mPartnerName"`
		MFormFieldsResponse []MFormFieldResponse `json:"mFormFields"`
		MFormCreatedAt      time.Time            `json:"mFormCreatedAt"`
		MFormUpdatedAt      time.Time            `json:"mFormUpdatedAT"`
	}

	MFormFieldResponse struct {
		MFieldTypeId             string            `json:"mFieldTypeId"`
		MFieldTypeName           string            `json:"mFieldTypeID"`
		MFormFieldId             string            `json:"mFormFieldId"`
		MFormFieldName           string            `json:"mFormFieldName"`
		MFormFieldIsMandatory    bool              `json:"mFormFieldIsMandatory"`
		MFormFieldOrdering       int               `json:"mFormFieldOrdering"`
		MFormFieldPlaceholder    string            `json:"mFormFieldPlaceholder"`
		MFormFieldChildsResponse []MFormFieldChild `json:"mFormFieldChilds"`
	}

	MFormFieldChildResponse struct {
		MFormFieldChildId   string `json:"mFormFieldChildId"`
		MFormFieldChildName string `json:"mFormFieldChildName"`
	}
)
