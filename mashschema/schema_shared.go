package mashschema

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strings"
)

// TODO: What do we need to keep here???

const (
	MashDataSourceSearch           = "search"
	MashDataSourceFilter           = "filter"
	MashDataSourceNameFilterRegexp = "filter_name"

	// MashDataSourceRequired Whether a foreign object must exist at the moment the query is issued.
	MashDataSourceRequired = "required"
	// MashDataSourceUnique Whether in the result of the matching sequence, there should be a single object
	// left.
	MashDataSourceUnique = "require_unique"

	// MashObjId Universal field name for the created timestamp
	MashObjId      = "id"
	MashObjCreated = "created"
	// MashObjUpdated Universal field for the created timestamp
	MashObjUpdated = "updated"
	// MashObjName Universal field for the name
	MashObjName        = "name"
	MashObjDescription = "description"
	// MashObjNamePrefix Universal field for the name prefix.
	MashObjNamePrefix = "name_prefix"

	// TODO: Search and replace the usage of these fields.
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

func minOf(opts ...int) int {
	rv := opts[0]
	for i := 1; i < len(opts); i++ {
		if opts[i] < rv {
			rv = opts[i]
		}
	}

	return rv
}

var compoundCommentRegex = regexp.MustCompile("^[a-zA-Z0-9_\\-:]*")

func IdWithComment(id, comment string) string {
	return fmt.Sprintf("%s # %s", id, comment)
}

func IdWithoutComment(id string) string {
	return strings.TrimSpace(compoundCommentRegex.FindString(id))
}

type Supplier func() interface{}
