package runner

// import (
// 	"fmt"
// 	"reflect"

// 	"github.com/appclacks/beckart/store"
// 	"github.com/google/uuid"
// )

// func templateSliceJSON(store *store.Store, slice []any) error {
// 	for i := range slice {
// 		value := slice[i]
// 		if reflect.ValueOf(value).Kind() == reflect.Map {
// 			m, ok := value.(map[string]interface{})
// 			if !ok {
// 				return fmt.Errorf("fail to cast to map type %s", reflect.TypeOf(value))
// 			}
// 			err := templateMapJSON(store, m)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if reflect.ValueOf(value).Kind() == reflect.Slice {
// 			slice, ok := value.([]any)
// 			if !ok {
// 				return fmt.Errorf("fail to cast to slice type %s", reflect.TypeOf(value))
// 			}
// 			err := templateSliceJSON(store, slice)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if reflect.ValueOf(value).Kind() == reflect.String {
// 			string, ok := value.(string)
// 			if ok {
// 				id := uuid.New().String()
// 				result, err := genTemplate(store, id, string)
// 				if err != nil {
// 					return err
// 				}
// 				slice[i] = result
// 			}
// 		}
// 	}
// 	return nil
// }

// func templateMapJSON(store *store.Store, bodyJSON map[string]interface{}) error {
// 	if len(bodyJSON) != 0 {
// 		for k, v := range bodyJSON {
// 			if reflect.ValueOf(v).Kind() == reflect.Map {
// 				m, ok := v.(map[string]interface{})
// 				if !ok {
// 					return fmt.Errorf("fail to cast to map type %s", reflect.TypeOf(v))
// 				}
// 				err := templateMapJSON(store, m)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			if reflect.ValueOf(v).Kind() == reflect.Slice {
// 				slice, ok := v.([]any)
// 				if !ok {
// 					return fmt.Errorf("fail to cast to slice type %s", reflect.TypeOf(v))
// 				}
// 				err := templateSliceJSON(store, slice)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			if reflect.ValueOf(v).Kind() == reflect.String {
// 				string, ok := v.(string)
// 				if !ok {
// 					return fmt.Errorf("fail to cast to string type %s", reflect.TypeOf(v))
// 				}

// 				id := uuid.New().String()
// 				result, err := genTemplate(store, id, string)
// 				if err != nil {
// 					return err
// 				}
// 				fmt.Println(result)
// 				bodyJSON[k] = result
// 			}
// 		}
// 	}
// 	return nil
// }
