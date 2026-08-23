package httpx

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

type TrueType bool
type FalseType bool

var True TrueType = true
var False FalseType = false

var trueConst = any(true)
var falseConst = any(false)

var TrueTypeSchemaTypes = map[reflect.Type]*jsonschema.Schema{
	reflect.TypeFor[TrueType](): {Type: "bool", Const: &trueConst},
}

var FalseTypeSchemaTypes = map[reflect.Type]*jsonschema.Schema{
	reflect.TypeFor[FalseType](): {Type: "bool", Const: &falseConst},
}

func (t TrueType) MarshalJSON() ([]byte, error) {
	return []byte("true"), nil
}

func (t *TrueType) UnmarshalJSON(data []byte) error {
	*t = True
	return nil
}

func (f FalseType) MarshalJSON() ([]byte, error) {
	return []byte("false"), nil
}

func (f *FalseType) UnmarshalJSON(data []byte) error {
	*f = False
	return nil
}
