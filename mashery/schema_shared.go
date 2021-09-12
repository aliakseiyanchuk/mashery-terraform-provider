package mashery

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
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

const CompoundIdSeparator = "::"

func DataSourceBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// --------------------------------------------------
		// Inputs
		MashDataSourceSearch: {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "V3 search criteria",
			Elem:        stringElem(),
		},
		MashDataSourceRequired: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "If true (default), then a service satisfying the search condition must exist. If such service doesn't exist, the error is generated",
		},
		MashDataSourceUnique: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "By default, where multiple matches would exist, any object will returned. When set to true, requires at most one matching object",
		},
		MashDataSourceNameFilterRegexp: {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Regular expression for service name",
			//ValidateDiagFunc: validateRegularExpressionSet,
			Elem: stringElem(),
		},
	}
}

// Validate that the set of strings contains valid regular expressions.
func validateRegularExpressionSet(i interface{}, path cty.Path) diag.Diagnostics {
	opts := schemaSetToStringArray(i)

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

// CreateCompoundId Create compound identifier, denoting traversal
func CreateCompoundId(comps ...string) string {
	return strings.Join(comps, CompoundIdSeparator)
}

// Creates a type-string schema used in the Elem mappings to save repetitive lines of code.
// TODO: Needs to be referenced from the code.
func stringElem() *schema.Schema {
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

func ParseCompoundId(id string, dest ...*string) {
	arr := strings.Split(IdWithoutComment(id), CompoundIdSeparator)
	for i := 0; i < minOf(len(arr), len(dest)); i++ {
		*dest[i] = arr[i]
	}
}

func IdWithComment(id, comment string) string {
	return fmt.Sprintf("%s # %s", id, comment)
}

func IdWithoutComment(id string) string {
	return strings.TrimSpace(compoundCommentRegex.FindString(id))
}

func ValidateCompoundIdent(i interface{}, path cty.Path, comp int) diag.Diagnostics {
	if str, ok := i.(string); ok {
		arr := strings.Split(str, CompoundIdSeparator)
		if len(arr) != comp {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Unexpected number of components in compound identifier",
				Detail:        fmt.Sprintf("Expected %d elements of compound Id, but got %d", comp, len(arr)),
				AttributePath: path,
			}}
		}

		// We should receive at non-empty components
		for idx, v := range arr {
			if len(v) == 0 {
				return diag.Diagnostics{diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Missing ident component",
					Detail:        fmt.Sprintf("Iden component at position %d has zero length", idx),
					AttributePath: path,
				}}
			}
		}

		return diag.Diagnostics{}
	} else {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Not a valid input",
			Detail:        fmt.Sprintf("Input to this field should be string, but got %s", reflect.TypeOf(i)),
			AttributePath: path,
		}}
	}
}
