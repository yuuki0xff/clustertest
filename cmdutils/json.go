package cmdutils

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSONInterface unmarshal abstract data with concrete type.
// Input json MUST contain "type" field. The concrete object creates in according to the specified type.
//
//   Input: {"type": "foo",  "extra-keys": "..."}
//   Output: &FooType{...} if type is "foo".
//           &BarType{...} if type is "bar".
//           ...
//   Usage: var data AbstractType = UnmarshalJSONInterface(js, initializer)
//
// This method needs when you want to unmarshal a struct contains interfaces. For example, following code will raise an
// error because encoding/json package can not unmarshal interface type. UnmarshalJSONInterface() provides a way to
// solve this problem. For details, see example_json_test.go.
//
//   type AbstractType interface {
//   	SomeMethods()
//   }
//   type Config struct {
//   	Type string
//   	Data AbstractType
//   }
//   func main() {
//   	js := []byte(`{"type":"foo", "data":{"a":1,"b":"2"}}`)
//   	conf := &Config{}
//      // json: cannot unmarshal object into Go struct field Config.Data of type main.AbstractType
//   	err := json.Unmarshal(js, conf)
//   	if err != nil {
//   		panic(err)
//   	}
//   }
func UnmarshalJSONInterface(js []byte, fn func(typeName string) (concrete interface{}, err error)) error {
	var m map[string]interface{}
	err := json.Unmarshal(js, &m)
	if err != nil {
		return err
	}

	typeName := m["type"]
	if typeName == "" {
		return errors.New(`missing "type" field`)
	}
	if t, ok := typeName.(string); !ok {
		return errors.New(`"type" field is not string`)
	} else {
		concrete, err := fn(t)
		if err != nil {
			return err
		}
		return json.Unmarshal(js, concrete)
	}
}

// See UnmarshalJSONInterface() docs.
func UnmarshalYAMLInterface(unmarshal func(interface{}) error, fn func(typeName string) (concrete interface{}, err error)) error {
	var m map[string]interface{}
	err := unmarshal(&m)
	if err != nil {
		return err
	}

	typeName := m["type"]
	if typeName == "" {
		return errors.New(`missing "type" field`)
	}
	if t, ok := typeName.(string); !ok {
		return errors.New(`"type" field is not string`)
	} else {
		concrete, err := fn(t)
		if err != nil {
			return err
		}
		return unmarshal(concrete)
	}
}
