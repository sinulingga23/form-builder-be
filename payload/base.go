package payload

type (
	Response struct {
		StatusCode int         `json:"statusCode"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
	}
)
