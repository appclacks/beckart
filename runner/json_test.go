package runner

// import (
// 	"testing"

// 	"github.com/appclacks/beckart/store"
// 	"github.com/stretchr/testify/assert"
// )

// func TestTemplateMapJSON(t *testing.T) {
// 	variables := make(map[string]any)
// 	variables["simple"] = "str"
// 	variables["simple2"] = "str2"
// 	variables["complex"] = []map[string]any{
// 		{
// 			"a": "bcd",
// 		},
// 	}
// 	store := store.New(variables)

// 	bodyJSON := make(map[string]interface{})

// 	bodyJSON["key"] = "val"
// 	bodyJSON["template"] = `{{index .Variables "simple"}}`
// 	bodyJSON["map"] = map[string]interface{}{
// 		"foo": `{{index .Variables "simple"}}`,
// 		"bar": []any{
// 			`{{index .Variables "simple"}}`,
// 		},
// 		"another": []any{
// 			map[string]any{
// 				"complex": `{{index .Variables "complex"}}`,
// 				"foo": []any{
// 					"bar",
// 					map[string]any{
// 						"a": `{{index .Variables "simple"}}`,
// 					},
// 					`{{index .Variables "simple2"}}`,
// 				},
// 			},
// 		},
// 	}
// 	err := templateMapJSON(store, bodyJSON)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "str", bodyJSON["template"])
// 	assert.Equal(t, "str", bodyJSON["map"].(map[string]any)["foo"])
// 	assert.Equal(t, "str", bodyJSON["map"].(map[string]any)["bar"].([]any)[0])
// 	assert.Equal(t, "[map[a:bcd]]", bodyJSON["map"].(map[string]any)["another"].([]any)[0].(map[string]any)["complex"])
// 	assert.Equal(t, "bar", bodyJSON["map"].(map[string]any)["another"].([]any)[0].(map[string]any)["foo"].([]any)[0])
// 	assert.Equal(t, "str", bodyJSON["map"].(map[string]any)["another"].([]any)[0].(map[string]any)["foo"].([]any)[1].(map[string]any)["a"])
// 	assert.Equal(t, "str2", bodyJSON["map"].(map[string]any)["another"].([]any)[0].(map[string]any)["foo"].([]any)[2])
// }
