package runner

import (
	"testing"

	"github.com/appclacks/beckart/store"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	cases := []struct {
		payload any
		key     any
		result  any
		err     string
	}{
		{
			payload: map[string]any{
				"foo": "bar",
			},
			key:    "foo",
			result: "bar",
		},
		{
			payload: map[string]any{
				"foo": map[string]string{
					"bar": "baz",
				},
			},
			key: "foo",
			result: map[string]string{
				"bar": "baz",
			},
		},
		{
			payload: []any{"a", "b", "c"},
			key:     1,
			result:  "b",
		},
		{
			payload: []any{},
			key:     1,
			err:     "no enough value",
		},
		{
			payload: []any{"1"},
			key:     "foo",
			err:     "expected int",
		},
		{
			payload: []string{"1"},
			key:     1,
			err:     "fail to extract",
		},
		{
			payload: map[string]any{"foo": "bar"},
			key:     1,
			err:     "invalid key type",
		},
		{
			payload: map[string]string{"foo": "bar"},
			key:     "foo",
			err:     "fail to extract",
		},
		{
			payload: map[string]any{"foo": "bar"},
			key:     "ko",
			err:     "no value found",
		},
	}
	for _, c := range cases {
		result, err := extractKey(c.payload, c.key)
		if c.err != "" {
			assert.ErrorContains(t, err, c.err)
		}

		assert.Equal(t, result, c.result)
	}
}

func TestExtractBody(t *testing.T) {
	cases := []struct {
		body       []byte
		extractor  map[string][]any
		storeKey   string
		storeValue any
		err        string
	}{
		{
			body:       []byte(`{"foo":"bar"}`),
			storeKey:   "var",
			storeValue: "bar",
			extractor: map[string][]any{
				"var": {"foo"},
			},
		},
		{
			body:       []byte(`[10, 11]`),
			storeKey:   "var",
			storeValue: float64(11),
			extractor: map[string][]any{
				"var": {1},
			},
		},
		{
			body:       []byte(`[10, 11]`),
			storeKey:   "var",
			storeValue: float64(10),
			extractor: map[string][]any{
				"var": {0},
			},
		},
	}
	variables := make(map[string]any)
	store := store.New(variables)

	for _, c := range cases {
		err := extractBody(store, c.body, c.extractor)
		if c.err != "" {
			assert.ErrorContains(t, err, c.err)
		}
		assert.Equal(t, store.Variables[c.storeKey], c.storeValue)
	}

}
