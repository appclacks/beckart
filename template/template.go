package template

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"text/template"

	"github.com/appclacks/beckart/store"
)

var funcMap = template.FuncMap{
	"base64Encode": func(v string) string {
		return base64.StdEncoding.EncodeToString([]byte(v))
	},
	"json": func(v any) (string, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	},
}

func GenTemplate(store *store.Store, name string, content string) (*bytes.Buffer, error) {
	tmpl := template.New(name)
	tmpl.Funcs(funcMap)
	t, err := tmpl.Parse(content)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, store); err != nil {
		return nil, err
	}
	return buf, nil
}
