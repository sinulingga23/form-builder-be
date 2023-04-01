package define

import "errors"

var (
	ErrMPartnerIdEmpty     = errors.New("Partner Id can't be empty.")
	ErrMPartnerNotFound    = errors.New("Partner not found.")
	ErrMFormNameEmpty      = errors.New("Form Name can't be empty.")
	ErrMFormFieldNameEmpty = errors.New("Form Field Name can't be empty.")
	ErrMFieldTypeIdEmpty   = errors.New("FieldType Id can't be empty.")
	ErrMFieldTypeNotFound  = errors.New("FieldType not found.")
	ErrInternalServerError = errors.New("Internal Server Error.")
	ErrQueryData           = errors.New("Error query dtata.")
)
