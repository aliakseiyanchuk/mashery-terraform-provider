package mashschema

import (
	"encoding/json"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"
)

const (
	CompoundIdSeparator = "::"
)

func enumerateStringFields(vif interface{}, op func(field reflect.StructField, val string)) {
	vof := reflect.ValueOf(vif)
	if vof.Kind() == reflect.Ptr {
		vof = vof.Elem()
	}

	t := reflect.TypeOf(vof.Interface())

	for i := 0; i < vof.NumField(); i++ {
		fld := vof.Field(i)

		if fld.CanInterface() {
			if fld.Kind() == reflect.String {
				valueField := vof.Field(i)
				str := valueField.Interface().(string)

				op(t.Field(i), str)
			} else if fld.Kind() == reflect.Struct {
				valueField := vof.Field(i)
				enumerateStringFields(&valueField, op)
			}
		}
	}
}

func CompoundId(identStruct interface{}) string {
	js, _ := json.Marshal(identStruct)
	return string(js)

	//rv := strings.Builder{}
	//enumerateStringFields(identStruct, func(field reflect.StructField, val string) *string {
	//	if rv.Len() > 0 {
	//		rv.WriteString(CompoundIdSeparator)
	//	}
	//	rv.WriteString(val)
	//	return nil
	//})
	//
	//return rv.String()
}

func CompoundIdFrom(identStruct interface{}, id string) bool {
	if err := json.Unmarshal([]byte(id), identStruct); err != nil {
		return false
	} else {
		return IsIdentified(identStruct)
	}
}

func IsIdentified(identStruct interface{}) bool {
	emptyFields := 0
	enumerateStringFields(identStruct, func(field reflect.StructField, val string) {
		if len(val) == 0 {
			emptyFields++
		}
	})

	return emptyFields == 0
}

func CompoundIdMalformedDiagnostic(path cty.Path) diag.Diagnostics {
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       "Incomplete id",
		Detail:        "Identifier supplies incomplete data or is malformed",
		AttributePath: path,
	}}
}
