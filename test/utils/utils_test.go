package testutils

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

func TestPrepareQueryParams(t *testing.T) {
	tests := []struct {
		params   interface{}
		expected map[string]string
		err      error
	}{
		{
			params: &models.GetTradesParams{
				Limit: nil,
			},
			expected: map[string]string{},
			err:      nil,
		},
		{
			params: &models.GetTradesParams{
				Limit:     api.PtrInt(10),
				StartTime: api.PtrInt64(20),
				EndTime:   api.PtrInt64(30),
			},
			expected: map[string]string{
				"limit":      "10",
				"start_time": "20",
				"end_time":   "30",
			},
			err: nil,
		},
		{
			params: &models.GetTradesParams{
				Limit:     api.PtrInt(10),
				StartTime: api.PtrInt64(20),
				EndTime:   api.PtrInt64(0),
			},
			expected: map[string]string{
				"limit":      "10",
				"start_time": "20",
				"end_time":   "0",
			},
			err: nil,
		},
		{
			params: &models.GetHistoricalPricesParams{
				Limit: api.PtrInt(10),
			},
			expected: map[string]string{},
			err:      errors.New("required field: resolution"),
		},
		{
			params: &models.GetHistoricalPricesParams{
				Resolution: models.Minute,
				Limit:      api.PtrInt(10),
				StartTime:  api.PtrInt(20),
				EndTime:    api.PtrInt(0),
			},
			expected: map[string]string{
				"resolution": "60",
				"limit":      "10",
				"start_time": "20",
				"end_time":   "0",
			},
			err: nil,
		},
	}

	for i, test := range tests {
		msg := fmt.Sprintf("test #%d", i+1)
		result, err := api.PrepareQueryParams(test.params)
		if err != nil {
			if test.err == nil {
				t.Fatal("test err should be nil")
			}
			if test.err.Error() != err.Error() {
				t.Fatalf("Should be equal: %s, %s - %s", test.err.Error(), err.Error(), msg)
			}
		}
		if len(result) != len(test.expected) {
			t.Fatalf("Length inequality: %d, %d, %s", len(result), len(test.expected), msg)
		}
		for k, v := range test.expected {
			value, ok := result[k]
			if !ok {
				t.Fatalf("Could not find result: %s, %s", k, msg)
			}
			if v != value {
				t.Fatalf("Should be equal: %+v, %+v, %s", v, value, msg)
			}
		}
	}
}
