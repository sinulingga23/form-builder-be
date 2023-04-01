package payload

type (
	CreateMFormRequest struct {
		MPartnerId  string       `json:"mPartnerId"`
		MFormName   string       `json:"mFormName"`
		MFormFields []MFormField `json:"mFormFields"`
	}

	MFormField struct {
		MFieldTypeId          string `json:"mFieldTypeId"`
		MFormFieldName        string `json:"mFormFieldName"`
		MFormFieldIsMandatory bool   `json:"mFormFieldIsMandatory"`
		MFormFieldOrdering    int    `json:"mFormFieldOrdering"`
		MFormFieldPlaceholder string `json:"mFormFieldPlaceholder"`
		MFormFieldChilds      []MFormFieldChild
	}

	MFormFieldChild struct {
		MFormFieldChildName string `json:"mFormFieldChildName"`
	}
)
