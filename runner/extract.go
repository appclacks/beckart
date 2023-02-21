package runner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/appclacks/beckart/store"
)

func extractKey(payload any, key any) (any, error) {
	if reflect.ValueOf(payload).Kind() == reflect.Map {
		k, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("invalid key type during extraction: expected string, got %v", reflect.TypeOf(key))
		}
		p, ok := payload.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("fail to extract key %s from slice: invalid type %v", k, reflect.TypeOf(payload))
		}
		result, ok := p[k]
		if !ok {
			return nil, fmt.Errorf("no value found for key %s", k)
		}
		return result, nil

	}
	if reflect.ValueOf(payload).Kind() == reflect.Slice {
		k, ok := key.(int)
		if !ok {
			return nil, fmt.Errorf("invalid key type during extraction: expected int, got %v", reflect.TypeOf(key))
		}
		p, ok := payload.([]any)
		if !ok {
			return nil, fmt.Errorf("fail to extract index %d from slice: invalid type %v", k, reflect.TypeOf(payload))
		}
		if k >= len(p) {
			return nil, fmt.Errorf("no enough value in the list to extract index %d", k)
		}
		return p[k], nil
	}

	return nil, fmt.Errorf("extract not supported for type %s", reflect.TypeOf(payload))
}

func extractBody(store *store.Store, body []byte, extractor map[string][]any) error {
	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	for variable, path := range extractor {
		var result any
		var err error
		for i, key := range path {
			if i == 0 {
				result, err = extractKey(payload, key)
				if err != nil {
					return err
				}
			} else {
				result, err = extractKey(result, key)
				if err != nil {
					return err
				}
			}
		}
		store.Set(variable, result)
	}
	return nil
}

func extractHeaders(store *store.Store, headers http.Header, extractor map[string]string) error {
	for variable, headerKey := range extractor {
		result := headers.Get(headerKey)
		if result == "" {
			return fmt.Errorf("fail to extract header %s: the header is empty", headerKey)
		}
		store.Set(variable, result)
	}
	return nil
}

func extractCookies(store *store.Store, cookies []*http.Cookie, extractor map[string]string) error {
	// double for not needed but let's optimize later
	for variable, cookieName := range extractor {
		for i := range cookies {
			cookie := cookies[i]
			if cookie.Name == cookieName {
				store.Set(variable, *cookie)
				break
			}
		}
		return fmt.Errorf("cookie %s not found in http response", cookieName)
	}
	return nil
}
