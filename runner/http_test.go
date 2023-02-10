package runner

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBuildURL(t *testing.T) {
	cases := []struct {
		action config.HTTPAction
		result string
	}{
		{
			action: config.HTTPAction{
				Target:   "localhost",
				Port:     10000,
				Protocol: "http",
			},
			result: "http://localhost:10000",
		},
		{
			action: config.HTTPAction{
				Target:   "127.0.0.1",
				Port:     10001,
				Protocol: "https",
			},
			result: "https://127.0.0.1:10001",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.result, buildURL(c.action))
	}
}

func TestIsSuccessful(t *testing.T) {
	cases := []struct {
		action   config.HTTPAction
		response *http.Response
		result   bool
	}{
		{
			action: config.HTTPAction{
				Target:   "localhost",
				Port:     10000,
				Protocol: "http",
				ValidStatus: []uint{
					200,
				},
			},
			response: &http.Response{
				StatusCode: 200,
			},
			result: true,
		},
		{
			action: config.HTTPAction{
				Target:   "localhost",
				Port:     10000,
				Protocol: "http",
				ValidStatus: []uint{
					201,
					301,
				},
			},
			response: &http.Response{
				StatusCode: 301,
			},
			result: true,
		},
		{
			action: config.HTTPAction{
				Target:   "localhost",
				Port:     10000,
				Protocol: "http",
				ValidStatus: []uint{
					201,
					301,
				},
			},
			response: &http.Response{
				StatusCode: 300,
			},
			result: false,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.result, isSuccessful(c.action, c.response))
	}
}

type httpTestCase struct {
	Action            config.Action
	ExpectedBody      string
	ExpectedHeaders   map[string]string
	ExpectedVariables map[string]any
}

var testBody = `[{"test1": "v1"}, {"test2": {"k3": "v2"}}, {"test3": {"a": [1, 2]}}]`

func TestExecuteHTTP(t *testing.T) {
	variables := make(map[string]any)
	variables["simple"] = "str"
	variables["complex"] = []map[string]any{
		{
			"a": "bcd",
		},
	}
	store := store.New(variables)

	var testCase *httpTestCase

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("my-header", "header-value")
		_, err := w.Write([]byte(testBody))
		assert.NoError(t, err)
		assert.Equal(t, r.Method, testCase.Action.HTTP.Method)
		bodyBytes, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		body := string(bodyBytes)
		assert.Equal(t, testCase.ExpectedBody, body, "invalid request body")
	}))
	port, err := strconv.ParseUint(strings.Split(ts.URL, ":")[2], 10, 16)
	assert.NoError(t, err)
	cases := []httpTestCase{
		{
			Action: config.Action{
				Name: "foo",
				HTTP: config.HTTPAction{ValidStatus: []uint{200},
					Target:   "localhost",
					Method:   "POST",
					Port:     uint(port),
					Timeout:  "5s",
					Protocol: "http",
					Path:     "/",
					Body:     `body: {{ index .Variables "simple"}} : {{ index .Variables "complex" 0 "a"}}`,
					Extractors: config.HTTPExtractors{
						BodyJSON: map[string][]any{
							"foo": {0, "test1"},
							"bar": {1, "test2"},
							"baz": {1, "test2", "k3"},
						},
						Headers: map[string]string{
							"var-header": "my-header",
						},
					},
				},
			},
			ExpectedBody: "body: str : bcd",
			ExpectedVariables: map[string]any{
				"var-header": "header-value",
				"foo":        "v1",
				"bar": map[string]any{
					"k3": "v2",
				},
				"baz": "v2",
			},
		},
		{
			Action: config.Action{
				Name: "bar",
				HTTP: config.HTTPAction{ValidStatus: []uint{200},
					Target:   "localhost",
					Method:   "POST",
					Port:     uint(port),
					Timeout:  "5s",
					Protocol: "http",
					Path:     "/",
					Body:     `{"foo": {{ json (index .Variables "complex")}}}`,
					Extractors: config.HTTPExtractors{
						BodyJSON: map[string][]any{
							"slice":  {2, "test3", "a"},
							"intvar": {2, "test3", "a", 0},
						},
					},
				},
			},
			ExpectedBody: `{"foo": [{"a":"bcd"}]}`,
			ExpectedVariables: map[string]any{
				"slice":  []any{float64(1), float64(2)},
				"intvar": float64(1),
			},
		},
		{
			Action: config.Action{
				Name: "baz",
				HTTP: config.HTTPAction{ValidStatus: []uint{200},
					Target:   "localhost",
					Method:   "POST",
					Port:     uint(port),
					Timeout:  "5s",
					Protocol: "http",
					Path:     "/",
					Body:     `{"foo": {{ json (index .Variables "complex")}}}`,
					Extractors: config.HTTPExtractors{
						Body: "fullbody",
					},
				},
			},
			ExpectedBody: `{"foo": [{"a":"bcd"}]}`,
			ExpectedVariables: map[string]any{
				"slice":    []any{float64(1), float64(2)},
				"intvar":   float64(1),
				"fullbody": testBody,
			},
		},
	}
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	for _, c := range cases {
		testCase = &c
		err := ExecuteHTTP(logger, store, nil, c.Action)
		assert.NoError(t, err)
		for k, v := range c.ExpectedVariables {
			result, ok := store.Get(k)
			assert.Truef(t, ok, "variable %s not found", k)
			assert.Equalf(t, v, result, "invalid value for variable %s", k)
		}
	}
}
