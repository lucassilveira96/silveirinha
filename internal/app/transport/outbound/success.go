package outbound

type (
	Success struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Params  interface{} `json:"params"`
	}

	SuccessWithCount struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Params  interface{} `json:"params"`
		Total   int         `json:"total"`
	}
)
