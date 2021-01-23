package api

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func PrepareQueryParams(params interface{}) (map[string]string, error) {
	result := make(map[string]string)

	val := reflect.ValueOf(params).Elem()
	if val.Kind() != reflect.Struct {
		return result, nil
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("json")

		switch valueField.Kind() {
		case reflect.Ptr:
			if valueField.IsNil() {
				continue
			}
			result[tag] = fmt.Sprintf("%v", valueField.Elem().Interface())
		default:
			if valueField.IsZero() {
				return result, errors.Errorf("required field: %v", tag)
			}
			result[tag] = fmt.Sprintf("%v", valueField.Interface())
		}
	}

	return result, nil
}

func FormURL(s string) string {
	return fmt.Sprintf("%s%s", apiUrl, s)
}

func PtrInt(i int) *int {
	return &i
}

func PtrInt64(i int64) *int64 {
	return &i
}

func PtrString(s string) *string {
	return &s
}

func PtrFloat64(f float64) *float64 {
	return &f
}

func PtrDecimal(d decimal.Decimal) *decimal.Decimal {
	return &d
}

func PtrBool(b bool) *bool {
	return &b
}

func PtrDuration(d time.Duration) *time.Duration {
	return &d
}
