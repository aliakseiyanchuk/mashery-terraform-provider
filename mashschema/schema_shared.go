package mashschema

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

// TODO: What do we need to keep here???

const (
	MashDataSourceSearch = "search"

	// MashDataSourceRequired Whether a foreign object must exist at the moment the query is issued.
	MashDataSourceRequired = "required"

	MashObjCreated = "created"
	// MashObjUpdated Universal field for the created timestamp
	MashObjUpdated = "updated"
	// MashObjName Universal field for the name
	MashObjName        = "name"
	MashObjDescription = "description"
	// MashObjNamePrefix Universal field for the name prefix.
	MashObjNamePrefix = "name_prefix"
)

var EmptyStringArray []string

// Validate that the set of strings contains valid regular expressions.
func validateRegularExpressionSet(i interface{}, path cty.Path) diag.Diagnostics {
	opts := SchemaSetToStringArray(i)

	rv := diag.Diagnostics{}

	for _, str := range opts {
		if err := isValidRegexp(str); err != nil {
			rv = append(rv, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid regexp",
				Detail:        fmt.Sprintf("Error in regular expression: %s", err.Error()),
				AttributePath: path,
			})
		}
	}

	return rv
}

func isValidRegexp(str string) error {
	_, err := regexp.Compile(str)
	return err
}

// Creates a type-string mashschema used in the Elem mappings to save repetitive lines of code.
// TODO: Needs to be referenced from the code.
func StringElem() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeString,
	}
}
