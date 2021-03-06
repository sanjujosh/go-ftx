package api

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	CommaOmitEmpty string = "," + OmitEmpty
	OmitEmpty      string = "omit_empty"
	ZeroString     string = ""
)

func PrepareQueryParams(params interface{}) (map[string]string, error) {

	result, val := make(map[string]string), reflect.ValueOf(params).Elem()

	if val.Kind() != reflect.Struct {
		return result, nil
	}

	for i := 0; i < val.NumField(); i++ {

		vf := val.Field(i)
		vt := val.Type().Field(i)
		tag := vt.Tag.Get("json")

		omitempty := strings.Contains(tag, OmitEmpty)
		tag = strings.ReplaceAll(
			strings.ReplaceAll(tag, " ", ZeroString),
			CommaOmitEmpty,
			ZeroString,
		)

		switch vf.Kind() {
		case reflect.Ptr:
			if vf.IsNil() {
				continue
			}
			result[tag] = fmt.Sprintf("%v", vf.Elem().Interface())
		default:
			if vf.IsZero() {
				if omitempty {
					continue
				} else {
					return result, errors.Errorf("Required field: %v", tag)
				}
			}
			result[tag] = fmt.Sprintf("%v", vf.Interface())
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
