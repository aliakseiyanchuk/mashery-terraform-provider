package mashschemag

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"terraform-provider-mashery/funcsupport"
	"terraform-provider-mashery/tfmapper"
	"testing"
)

func TestMaps(t *testing.T) {
	obj := masherytypes.MasheryOAuth{}

	val := reflect.Indirect(reflect.ValueOf(&obj))
	fld := val.FieldByName("ForwardedHeaders")

	assert.Equal(t, reflect.Slice, fld.Type().Kind())
	assert.Equal(t, reflect.String, fld.Type().Elem().Kind())
}

func autoTestIdentity[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], ref Ident) {
	mapper := sb.Mapper()
	testState := sb.TestResourceData()

	err := mapper.AssignIdentity(ref, testState)
	assert.Nil(t, err)

	readBack, idErr := mapper.Identity(testState)
	assert.Nil(t, idErr)

	assert.True(t, reflect.DeepEqual(ref, readBack))
}

func autoTestParentIdentity[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], ref ParentIdent) {
	mapper := sb.Mapper()
	testState := sb.TestResourceData()

	err := mapper.TestSetPrentIdentity(ref, testState)
	assert.Nil(t, err)

	readBack, idErr := mapper.ParentIdentity(testState)
	assert.Nil(t, idErr)

	assert.True(t, reflect.DeepEqual(ref, readBack))
}

func autoTestMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], except ...string) {
	autoTestBoolMappings(t, sb, supplier, except...)
	autoTestStringMappings(t, sb, supplier, except...)
	autoTestIntMappings(t, sb, supplier, except...)
	autoTestInt64PtrMappings(t, sb, supplier, except...)
	autoTestEAVMappings(t, sb, supplier, except...)
	autoTestStringArrayMappings(t, sb, supplier, except...)
}

func autoTestBoolMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], except ...string) {
	ref := supplier()
	boolFields := matchingFieldsOf(ref, func(in reflect.Type) bool {
		return in.Kind() == reflect.Bool
	}, except...)

	for _, fldName := range boolFields {
		mapper := sb.Mapper()
		testState := sb.TestResourceData()

		in := supplier()
		reflectSetBool(&in, fldName, true)

		mapper.RemoteToSchema(&in, testState)

		readBack := supplier()
		mapper.SchemaToRemote(testState, &readBack)

		assert.True(t, reflectGetBool(readBack, fldName), "mismatching read/write on field %s", fldName)
	}
}

func reflectSetBool(i interface{}, fldName string, boolVal bool) {
	val := reflect.Indirect(reflect.ValueOf(i))
	val.FieldByName(fldName).SetBool(boolVal)
}

func reflectGetBool(i interface{}, fldName string) bool {
	val := reflect.Indirect(reflect.ValueOf(i))
	return val.FieldByName(fldName).Bool()
}

func autoTestStringMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], except ...string) {
	ref := supplier()
	stringFields := matchingFieldsOf(ref, func(in reflect.Type) bool {
		return in.Kind() == reflect.String
	}, except...)

	for _, fldName := range stringFields {
		fmt.Println("Testing field " + fldName)
		mapper := sb.Mapper()
		testState := sb.TestResourceData()

		in := supplier()

		fldValue := "string-under-test-" + fldName
		reflectSetString(&in, fldName, fldValue)

		mapper.RemoteToSchema(&in, testState)

		readBack := supplier()
		mapper.SchemaToRemote(testState, &readBack)

		assert.Equal(t, fldValue, reflectGetString(readBack, fldName), "mismatching read/write on string field %s", fldName)
	}
}

func reflectSetString(i interface{}, fldName string, stringVal string) {
	val := reflect.Indirect(reflect.ValueOf(i))
	val.FieldByName(fldName).SetString(stringVal)
}

func reflectGetString(i interface{}, fldName string) string {
	val := reflect.Indirect(reflect.ValueOf(i))
	return val.FieldByName(fldName).String()
}

func autoTestStringArrayMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], except ...string) {
	ref := supplier()
	stringFields := matchingFieldsOf(ref, func(in reflect.Type) bool {
		return in.Kind() == reflect.Slice && in.Elem().Kind() == reflect.String
	}, except...)

	for _, fldName := range stringFields {
		mapper := sb.Mapper()
		testState := sb.TestResourceData()

		in := supplier()

		fldValue := []string{"string-array-under-test-0-" + fldName, "string-array-under-test-1-" + fldName, "string-array-under-test-2-" + fldName}
		reflectSetStringArray(&in, fldName, fldValue)

		mapper.RemoteToSchema(&in, testState)

		readBack := supplier()
		mapper.SchemaToRemote(testState, &readBack)

		rbVal := reflectGetStringArray(readBack, fldName)

		assertArrayIn(t, fldName, &fldValue, &rbVal)
		assertArrayIn(t, fldName, &rbVal, &fldValue)
	}
}

func assertArrayIn(t *testing.T, field string, in *[]string, dest *[]string) {
outer:
	for _, k := range *in {
		for _, v := range *dest {
			if k == v {
				continue outer
			}
		}

		assert.Failf(t, "mismatching read/write on string array field %s: string %s is not found in the target array", field, k)
	}
}

func reflectSetStringArray(i interface{}, fldName string, stringVal []string) {
	val := reflect.Indirect(reflect.ValueOf(i))
	val.FieldByName(fldName).Set(reflect.ValueOf(stringVal))
}

func reflectGetStringArray(i interface{}, fldName string) []string {
	val := reflect.Indirect(reflect.ValueOf(i))
	return val.FieldByName(fldName).Interface().([]string)
}

func matchingFieldsOf(i interface{}, predicate funcsupport.Predicate[reflect.Type], except ...string) []string {
	rvFieldsRaw := enumerateFields(i, predicate)
	var rvFields []string

outer:
	for _, fld := range rvFieldsRaw {
		for _, exclField := range except {
			if fld == exclField {
				continue outer
			}
		}
		rvFields = append(rvFields, fld)
	}

	return rvFields
}

// ------------------------------------------------------------------------------
// Int field mappers

func autoTestIntMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], expectFields ...string) {
	ref := supplier()

	intFields := matchingFieldsOf(&ref, func(in reflect.Type) bool {
		return in.Kind() == reflect.Int
	}, expectFields...)

	for idx, fldName := range intFields {
		mapper := sb.Mapper()
		testState := sb.TestResourceData()

		in := supplier()

		fldValue := 100 + idx
		reflectSetInt(&in, fldName, fldValue)

		mapper.RemoteToSchema(&in, testState)

		readBack := supplier()
		mapper.SchemaToRemote(testState, &readBack)

		assert.Equal(t, fldValue, reflectGetInt(readBack, fldName), "mismatching read/write on int field %s", fldName)
	}
}

func reflectSetInt(i interface{}, fldName string, intVal int) {
	val := reflect.Indirect(reflect.ValueOf(i))
	val.FieldByName(fldName).SetInt(int64(intVal))
}

func reflectGetInt(i interface{}, fldName string) int {
	val := reflect.Indirect(reflect.ValueOf(i))
	return int(val.FieldByName(fldName).Int())
}

func autoTestInt64PtrMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], expectFields ...string) {
	ref := supplier()

	intFields := matchingFieldsOf(&ref, func(in reflect.Type) bool {
		return in.Kind() == reflect.Ptr && in.Elem().Kind() == reflect.Int64
	}, expectFields...)

	for idx, fldName := range intFields {
		mapper := sb.Mapper()
		testState := sb.TestResourceData()

		in := supplier()

		var fldValue = int64(100 + idx)
		reflectSetInt64Ptr(&in, fldName, &fldValue)

		mapper.RemoteToSchema(&in, testState)

		readBack := supplier()
		mapper.SchemaToRemote(testState, &readBack)

		readBackPtr := reflectGetInt64Ptr(&readBack, fldName)
		assert.NotNil(t, readBackPtr, "read back pointer should not be null at this point")
		assert.Equal(t, fldValue, *readBackPtr, "mismatching read/write on int field %s", fldName)
	}
}

func reflectSetInt64Ptr(i interface{}, fldName string, intVal *int64) {
	val := reflect.Indirect(reflect.ValueOf(i))
	val.FieldByName(fldName).Set(reflect.ValueOf(intVal))
}

func reflectGetInt64Ptr(i interface{}, fldName string) *int64 {
	val := reflect.Indirect(reflect.ValueOf(i))
	if rvPtr := val.FieldByName(fldName).Interface(); rvPtr != nil {
		return rvPtr.(*int64)
	} else {
		return nil
	}
}

func autoTestEAVMappings[ParentIdent, Ident, MType any](t *testing.T, sb *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], supplier funcsupport.Supplier[MType], expectFields ...string) {
	ref := supplier()

	mapFields := matchingFieldsOf(&ref, func(in reflect.Type) bool {
		return in.Kind() == reflect.Ptr && in.Elem().Name() == "EAV"
	}, expectFields...)

	for _, fldName := range mapFields {
		mapper := sb.Mapper()
		testState := sb.TestResourceData()

		in := supplier()

		var fldValue = masherytypes.EAV{
			"Field-" + fldName: "Field-" + fldName + "-Value",
		}
		reflectSetEAV(&in, fldName, &fldValue)

		mapper.RemoteToSchema(&in, testState)

		readBack := supplier()
		mapper.SchemaToRemote(testState, &readBack)

		readBackPtr := reflectGetEAV(&readBack, fldName)
		assert.NotNil(t, readBackPtr, "read back pointer should not be null at this point")
		assert.True(t, reflect.DeepEqual(&fldValue, readBackPtr), "mismatching read/write on map field %s", fldName)
	}
}

func reflectSetEAV(i interface{}, fldName string, intVal *masherytypes.EAV) {
	val := reflect.Indirect(reflect.ValueOf(i))
	val.FieldByName(fldName).Set(reflect.ValueOf(intVal))
}

func reflectGetEAV(i interface{}, fldName string) *masherytypes.EAV {
	val := reflect.Indirect(reflect.ValueOf(i))
	if rvPtr := val.FieldByName(fldName).Interface(); rvPtr != nil {
		return rvPtr.(*masherytypes.EAV)
	} else {
		return nil
	}
}

// --------------
// Enumerate the fields of struct based on the predicate of the field kind type

func enumerateFields(i interface{}, predicate funcsupport.Predicate[reflect.Type]) []string {
	var rv []string

	val := reflect.Indirect(reflect.ValueOf(i))

	for i := 0; i < val.NumField(); i++ {
		if predicate(val.Field(i).Type()) {
			rv = append(rv, val.Type().Field(i).Name)
		}
	}

	return rv
}
