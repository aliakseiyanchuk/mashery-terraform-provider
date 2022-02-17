package mashschema

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"
	"time"
)

func validateDuration(i interface{}, path cty.Path) diag.Diagnostics {
	if _, err := time.ParseDuration(i.(string)); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "invalid duration",
			Detail:        fmt.Sprintf("expression %s is not a valid duration expression", i),
			AttributePath: path,
		}}
	} else {
		return diag.Diagnostics{}
	}
}

func validateZeroOrGreater(i interface{}, path cty.Path) diag.Diagnostics {
	if v, ok := i.(int); ok {
		if v < 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Field must be zero or positive",
				Detail:        fmt.Sprintf("Value %d is negative", v),
				AttributePath: path,
			}}
		} else {
			return nil
		}
	} else if v, ok := i.(int64); ok {
		if v < 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Field must be zero or positive",
				Detail:        fmt.Sprintf("Value %d is negative", v),
				AttributePath: path,
			}}
		} else {
			return nil
		}
	}

	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       "int or in64 required at this path",
		Detail:        fmt.Sprintf("unsupported type is %s", reflect.TypeOf(i).Name()),
		AttributePath: path,
	}}

}

func ValidateCompoundIdent(i interface{}, path cty.Path, supplier Supplier) diag.Diagnostics {
	if str, ok := i.(string); ok {
		ci := supplier()
		if !CompoundIdFrom(ci, str) {
			return CompoundIdMalformedDiagnostic(path)
		} else {
			return nil
		}
	} else {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Not a valid input",
			Detail:        fmt.Sprintf("Input to this field should be string, but got %s", reflect.TypeOf(i)),
			AttributePath: path,
		}}
	}
}

func validateDiagInputIsEndpointMethodIdentifier(i interface{}, path cty.Path) diag.Diagnostics {
	if str, ok := i.(string); ok {
		mid := ServiceEndpointMethodIdentifier{}
		CompoundIdFrom(&mid, str)

		if !IsIdentified(&mid) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Incomplete identifier",
				Detail:        "Endpoint method identifier is incomplete or malformed",
				AttributePath: path,
			}}
		} else {
			return diag.Diagnostics{}
		}
	} else {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Unexpected type",
			Detail:        fmt.Sprintf("Input should be string, but was %s", reflect.TypeOf(i)),
			AttributePath: path,
		}}
	}
}
