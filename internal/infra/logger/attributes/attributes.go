package attributes

import (
	"reflect"

	"github.com/pkg/errors"
)

type Attributes map[string]interface{}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func New() Attributes {
	return Attributes{}
}

func (attr Attributes) With(key string, value interface{}) Attributes {
	attr[key] = value
	return attr
}

func (attr Attributes) WithError(err error) Attributes {
	if err == nil {
		return attr
	}

	attr["exception.type"] = typeOfError(err)
	attr["exception.message"] = err.Error()

	if cause := errors.Cause(err); cause != nil {
		attr["exception.cause"] = cause.Error()
	}

	if st, ok := err.(stackTracer); ok {
		attr["exception.stacktrace"] = st.StackTrace()
	}

	return attr
}

func typeOfError(err error) string {
	if err == nil {
		return "nil"
	}

	t := reflect.TypeOf(err)
	if t.Kind() == reflect.Ptr {
		return t.Elem().String()
	}
	return t.String()
}
