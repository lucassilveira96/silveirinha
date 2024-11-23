package mapper

import (
	"reflect"
	"silveirinha/internal/app/domain/model"
	"silveirinha/internal/app/transport/inbound"
)

func TestepppppMapToModel(inbound inbound.Testeppppp) model.Testeppppp {
	var modelObj model.Testeppppp
	inboundValue := reflect.ValueOf(inbound)
	modelValue := reflect.ValueOf(&modelObj).Elem()

	// Loop through each field in the inbound struct
	for i := 0; i < inboundValue.NumField(); i++ {
		fieldName := inboundValue.Type().Field(i).Name
		modelField := modelValue.FieldByName(fieldName)

		// If the field exists in the model, copy the value
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(inboundValue.Field(i))
		}
	}

	return modelObj
}
