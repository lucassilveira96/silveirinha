package presenter

import (
	"reflect"
	"teste/internal/app/transport/outbound"
)

func Success(message string, params interface{}) *outbound.Success {
	if params == nil {
		params = make(map[string]interface{}, 0)
	}

	return &outbound.Success{
		Status:  "success",
		Message: message,
		Params:  params,
	}
}

func SuccessWithCount(message string, params interface{}) *outbound.SuccessWithCount {
	if params == nil {
		params = make(map[string]interface{}, 0)
	}

	return &outbound.SuccessWithCount{
		Status:  "success",
		Message: message,
		Params:  params,
		Total:   reflect.ValueOf(params).Len(),
	}
}
