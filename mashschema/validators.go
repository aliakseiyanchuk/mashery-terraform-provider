package mashschema

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"regexp"
	"time"
)

func ValidateNonEmptyString(i interface{}, path cty.Path) diag.Diagnostics {
	rv := diag.Diagnostics{}
	str := i.(string)
	if len(str) == 0 {
		rv = append(rv, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "invalid argument: this field string cannot be empty",
			AttributePath: path,
		})
	}
	return rv
}

func ValidateDuration(i interface{}, path cty.Path) diag.Diagnostics {
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

func ValidateRegExp(expr string) schema.SchemaValidateDiagFunc {
	regExp := regexp.MustCompile(expr)
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		rv := diag.Diagnostics{}
		strVal := i.(string)
		if !regExp.MatchString(strVal) {
			rv = append(rv, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "malformed input parameter",
				Detail:        fmt.Sprintf("value %s doesn't match regular expression %s", strVal, expr),
				AttributePath: path,
			})
		}

		return rv
	}
}

func ValidateZeroOrGreater(i interface{}, path cty.Path) diag.Diagnostics {
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
