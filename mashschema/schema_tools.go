package mashschema

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"regexp"
	"strings"
)

func cloneAsComputedSchemaElem(inp *schema.Schema) interface{} {
	if inpSchema, ok := inp.Elem.(*schema.Schema); ok {
		return cloneAsComputedSchema(inpSchema, false)
	}

	if inpResource, ok := inp.Elem.(*schema.Resource); ok {
		return &schema.Resource{
			Schema: CloneAsComputed(inpResource.Schema),
		}
	}

	return nil
}

func cloneAsComputedSchema(inp *schema.Schema, isKey bool) *schema.Schema {
	rv := schema.Schema{
		Type: inp.Type,
	}
	if isKey {
		rv.Computed = true
	}

	switch inp.Type {
	case schema.TypeSet:
		fallthrough
	case schema.TypeMap:
		fallthrough
	case schema.TypeList:
		rv.Elem = cloneAsComputedSchemaElem(inp)
		break
	default:
		break
	}

	return &rv
}

// CloneAsComputed Clone mashschema as a computer mashschema
func CloneAsComputed(inp map[string]*schema.Schema) map[string]*schema.Schema {
	rv := make(map[string]*schema.Schema, len(inp))

	for k, v := range inp {
		rv[k] = cloneAsComputedSchema(v, true)
	}

	return rv
}

// Append the source mashschema as computed mashschema in the target map.
// source: source mashschema
// dest: destination mashschema
func appendAsComputedInto(source *map[string]*schema.Schema, dest *map[string]*schema.Schema) {
	rv := *dest

	for k, v := range *source {
		rv[k] = cloneAsComputedSchema(v, true)
	}
}

func FindInArray(v string, lst *[]string) int {
	for idx, lv := range *lst {
		if v == lv {
			return idx
		}
	}

	return -1
}

func toStringArray(inp interface{}) []string {
	if elem, ok := inp.([]interface{}); ok {
		rv := make([]string, len(elem))
		for idx, v := range elem {
			if str, ok := v.(string); ok {
				rv[idx] = str
			} else {
				rv[idx] = fmt.Sprintf("%s", v)
			}
		}

		return rv
	} else {
		return []string{}
	}
}

func ValidateStringValueInSet(inp interface{}, pth cty.Path, lst *[]string) diag.Diagnostics {
	rv := diag.Diagnostics{}

	if str, ok := inp.(string); ok {
		if FindInArray(str, lst) < 0 {
			rv = append(rv, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "unacceptable value",
				Detail:        fmt.Sprintf("string value '%s' is not one of the allowed options", str),
				AttributePath: pth,
			})
		}
	} else {
		rv = append(rv, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "incompatible type",
			Detail:        fmt.Sprintf("input should be string, but is %s", reflect.TypeOf(inp)),
			AttributePath: pth,
		})
	}

	return rv
}

func ValidateIntValueInSet(inp interface{}, pth cty.Path, lst *[]int) diag.Diagnostics {
	v := inp.(int)
	for _, c := range *lst {
		if v == c {
			return diag.Diagnostics{}
		}
	}

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Specified value not supported",
			Detail:        fmt.Sprintf("value %d is not a valid option", v),
			AttributePath: pth,
		},
	}
}

func validateListValues(inp interface{}, pth cty.Path, lst *[]string) diag.Diagnostics {
	if elem, ok := inp.([]interface{}); ok {
		rv := diag.Diagnostics{}

		for idx, inpElem := range elem {
			if str, ok := inpElem.(string); ok {
				if FindInArray(str, lst) < 0 {
					rv = append(rv, diag.Diagnostic{
						Severity:      diag.Error,
						AttributePath: pth,
						Summary:       "value outside of supported options",
						Detail:        fmt.Sprintf("The value should be string from these options: %s", strings.Join(*lst, ",")),
					})
				}
			} else {
				rv = append(rv, diag.Diagnostic{
					Severity:      diag.Error,
					AttributePath: pth,
					Summary:       fmt.Sprintf("Value at index %d not a string but %s", idx, reflect.TypeOf(inpElem)),
				})
			}

			return rv
		}

		return diag.Diagnostics{}
	} else {
		return diag.Errorf("Expected string array type at path %s", pth)
	}
}

func inheritAll(source *map[string]*schema.Schema, dest *map[string]*schema.Schema) {
	for k, v := range *source {
		(*dest)[k] = v
	}
}

func getSetLength(inp interface{}) int {
	if inpAsSet, ok := inp.(*schema.Set); ok {
		return inpAsSet.Len()
	} else if inpAsArr, ok := inp.([]interface{}); ok {
		return len(inpAsArr)
	} else {
		return -1
	}
}

// UnwrapStructFromTerraformSet Unwraps a struct -- or, effectively, map[string]interface{} from terraform-encoded
// set input.
func UnwrapStructFromTerraformSet(inp interface{}) map[string]interface{} {
	if inpAsSet, ok := inp.(*schema.Set); ok {
		if inpAsSet.Len() >= 1 {
			uStruct := inpAsSet.List()[0]
			if rv, ok := uStruct.(map[string]interface{}); ok {
				return rv
			}
		}
	} else if inpAsArr, ok := inp.([]interface{}); ok {
		uStruct := inpAsArr[0]
		if rv, ok := uStruct.(map[string]interface{}); ok {
			return rv
		}
	}

	return make(map[string]interface{}, 0)
}

func UnwrapStructArrayFromTerraformSet(inp interface{}) []map[string]interface{} {
	if inpAsSet, ok := inp.(*schema.Set); ok {
		rv := make([]map[string]interface{}, inpAsSet.Len())
		for idx, v := range inpAsSet.List() {
			if setVal, ok := v.(map[string]interface{}); ok {
				rv[idx] = setVal
			}
		}

		return rv
	} else if inpAsArr, ok := inp.([]interface{}); ok {
		rv := make([]map[string]interface{}, len(inpAsArr))
		for idx, v := range inpAsArr {
			if setVal, ok := v.(map[string]interface{}); ok {
				rv[idx] = setVal
			}
		}

		return rv
	}

	return make([]map[string]interface{}, 0)
}

func ExtractBool(d *schema.ResourceData, key string, impliedValue bool) bool {
	if v := d.Get(key); v != nil {
		return v.(bool)
	} else {
		return impliedValue
	}
}

func ExtractString(d *schema.ResourceData, key string, impliedValue string) string {
	if v, exists := d.GetOk(key); exists {
		return v.(string)
	} else {
		return impliedValue
	}
}

func ToSet(inp []interface{}) *schema.Set {
	rv := schema.Set{}
	for k := range inp {
		rv.Add(k)
	}
	return &rv
}

func convertSetToStringArray(inp interface{}) []string {
	if set, ok := inp.(*schema.Set); ok {
		rv := make([]string, set.Len())
		for idx, v := range set.List() {
			if str, ok := v.(string); ok {
				rv[idx] = str
			} else {
				rv[idx] = fmt.Sprintf("%s", v)
			}
		}

		return rv
	} else {
		return []string{}
	}
}

func extractSetOrPrefixedString(d *schema.ResourceData, key, prefix string) string {
	if v, exists := d.GetOk(key); exists {
		return v.(string)
	} else {
		var rv string

		if prefix, exists := d.GetOk(prefix); exists {
			rv = resource.PrefixedUniqueId(prefix.(string))
		} else {
			rv = resource.UniqueId()
		}

		_ = d.Set(key, rv)

		return rv
	}
}

func ExtractStringPointer(d *schema.ResourceData, key string) *string {
	if v, exists := d.GetOk(key); exists {
		rv := v.(string)
		return &rv
	} else {
		return nil
	}
}

func ExtractIfModifiedStringPointer(d *schema.ResourceData, key string) *string {
	if d.HasChange(key) {
		if v, exists := d.GetOk(key); exists {
			rv := v.(string)
			return &rv
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func ExtractInt(d *schema.ResourceData, key string, impliedValue int) int {
	if v, exists := d.GetOk(key); exists {
		return v.(int)
	} else {
		return impliedValue
	}
}

func ExtractInt64(d *schema.ResourceData, key string, impliedValue int64) int64 {
	if v, exists := d.GetOk(key); exists {
		return v.(int64)
	} else {
		return impliedValue
	}
}

func ExtractIntPointer(d *schema.ResourceData, key string) *int {
	if v, exists := d.GetOk(key); exists {
		if rv, ok := v.(int); ok {
			return &rv
		}
	}

	// If the key does not exist, or if the conversion to int was not possible,
	// then nil will be returned.
	return nil
}

func intPointerToInt(inp *int64) int {
	if inp != nil {
		return int(*inp)
	} else {
		return -1
	}
}

func ExtractInt64Pointer(d *schema.ResourceData, key string, threshold int64) *int64 {
	if v, exists := d.GetOk(key); exists {
		if rv, ok := v.(int); ok {
			rv64 := int64(rv)
			if rv64 > threshold {
				return &rv64
			}
		}
	}

	// If the key does not exist, or if the conversion to int was not possible,
	// then nil will be returned.
	return nil
}

func SchemaMapToStringMap(v interface{}) map[string]string {
	return SchemaMapToStringMapWithEmpty(v, false)
}

func SchemaMapToStringMapWithEmpty(v interface{}, emptyAsNil bool) map[string]string {
	if mp, ok := v.(map[string]interface{}); ok {
		if emptyAsNil && len(mp) == 0 {
			return nil
		}

		rv := make(map[string]string, len(mp))

		for k, v := range mp {
			if str, ok := v.(string); ok {
				rv[k] = str
			} else {
				rv[k] = fmt.Sprintf("%s", v)
			}
		}

		return rv
	} else if mp, ok := v.(map[string]string); ok {
		if emptyAsNil && len(mp) == 0 {
			return nil
		}

		rv := make(map[string]string, len(mp))
		for k, v := range mp {
			rv[k] = v
		}

		return rv
	} else {
		if emptyAsNil && len(mp) == 0 {
			return nil
		} else {
			return map[string]string{}
		}
	}
}

// ExtractStringMap extracts and clone the map from the resource data. The can be modified by the calling client as needed
// without affecting the state of data in the mashschema.ResourceData structure.
func ExtractStringMap(d *schema.ResourceData, key string) map[string]string {
	if v, exists := d.GetOk(key); exists {
		return SchemaMapToStringMap(v)
	}

	return map[string]string{}
}

func extractEAVPointer(d *schema.ResourceData, key string) *masherytypes.EAV {
	if v, exists := d.GetOk(key); exists {
		if mp, ok := v.(map[string]interface{}); ok {
			rv := masherytypes.EAV{}

			for k, v := range mp {
				if str, ok := v.(string); ok {
					rv[k] = str
				} else {
					rv[k] = fmt.Sprintf("%s", v)
				}
			}

			return &rv
		}
	}

	return nil
}

func extractStringFieldOfStruct(d *schema.ResourceData, key string, structKey string) string {
	if v, exists := d.GetOk(key); exists {
		if mp, ok := v.(*schema.Set); ok {
			if len(mp.List()) == 1 {
				if strMap, ok := mp.List()[0].(map[string]interface{}); ok {
					rv := strMap[structKey]
					if str, ok := rv.(string); ok {
						return str
					} else {
						return fmt.Sprintf("%s", str)
					}
				}
			}
		}
	}

	return ""
}

// Utility method allowing using Set or List interchangeably when it comes to the extraction
// of data from the Terraform resource state.
func SchemaSetToStringArray(v interface{}) []string {
	if schSet, ok := v.(*schema.Set); ok {
		rv := make([]string, schSet.Len())
		for i, iStr := range schSet.List() {
			if str, ok := iStr.(string); ok {
				rv[i] = str
			} else {
				rv[i] = fmt.Sprintf("%s", iStr)
			}
		}

		return rv
	} else if arr, ok := v.([]interface{}); ok {
		rv := make([]string, len(arr))
		for idx, v := range arr {
			if str, ok := v.(string); ok {
				rv[idx] = str
			} else {
				rv[idx] = fmt.Sprintf("%s", v)
			}
		}

		return rv
	} else if arr, ok := v.([]string); ok {
		return arr
	} else {
		return []string{}
	}
}

func ConvertInterfaceArrayToStringArray(inp []interface{}) []string {
	rv := make([]string, len(inp))
	for idx, vraw := range inp {
		if str, ok := vraw.(string); ok {
			rv[idx] = str
		} else {
			rv[idx] = fmt.Sprintf("%s", vraw)
		}
	}

	return rv
}

// ExtractStringArray Extract a string array from the Terraform's internal data structure, otherwise supplying
// implied value
func ExtractStringArray(d *schema.ResourceData, key string, implied *[]string) []string {
	if v, exists := d.GetOk(key); exists {
		return SchemaSetToStringArray(v)
	}

	// Return the implied data elements, since the input is not understood.
	rv := make([]string, len(*implied))
	for i, v := range *implied {
		rv[i] = v
	}

	return rv
}

func SafeLookupStringPointer(src *map[string]interface{}, key string) *string {
	v := (*src)[key]
	if v != nil {
		if str, ok := v.(string); ok {
			return &str
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func afterApplyKnown(source string) string {
	if len(source) > 0 {
		return source
	} else {
		return ("known after apply")
	}
}

func NullForEmptyString(inp string) interface{} {
	if len(inp) > 0 {
		return inp
	} else {
		return nil
	}
}

func nilArrayForEmptyString(inp, sep string) interface{} {
	if len(inp) == 0 {
		return nil
	} else {
		return strings.Split(inp, sep)
	}
}

func toRegexpArray(exp []string) []*regexp.Regexp {
	rgxArr := make([]*regexp.Regexp, len(exp))
	for idx, v := range exp {
		rgxArr[idx] = regexp.MustCompile(v)
	}

	return rgxArr
}
