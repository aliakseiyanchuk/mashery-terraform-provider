package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

// DatasourceTemplate a template for performing the query operations. The template accepts functions performing
// Mashery querying and converting from Mashery V3 to the Terraform schema.
type DatasourceTemplate struct {
	Mapper mashschema.DataSourceMapper

	// RequireUnique whether the schema of the data query operation requires a unique match.
	RequireUnique bool

	// QueryExtractor function that that will build the query parameters. {@link DefaultQueryExtractor}
	// provides a sensible default
	QueryExtractor func(*schema.ResourceData) (map[string]string, error)

	// Query Execute a V3 query
	Query func(ctx context.Context, cl v3client.Client, query map[string]string) ([]interface{}, error)
}

// TFDataSourceSchema returns the Terraform data source schema
func (t *DatasourceTemplate) TFDataSourceSchema() *schema.Resource {
	// Panic if necessary functions were not supplied
	if t.Mapper == nil || t.Query == nil {
		panic(fmt.Sprintf("Unsatisfied initialization for object %s", t.Mapper.V3ObjectName()))
	}

	return &schema.Resource{
		ReadContext: t.TemplateQuery,
		Schema:      t.Mapper.TerraformSchema(),
	}
}

func (t *DatasourceTemplate) isMatchRequired(d *schema.ResourceData) bool {
	return mashschema.ExtractBool(d, mashschema.MashDataSourceRequired, true)
}

func (t *DatasourceTemplate) isUniqueMatchRequired() bool {
	return t.RequireUnique
}

func (t *DatasourceTemplate) actualQuery(d *schema.ResourceData) (map[string]string, error) {
	if t.QueryExtractor != nil {
		return t.QueryExtractor(d)
	} else {
		return map[string]string{}, nil
	}
}

// DataSourceDefaultQueryExtractor default query extract for the data sources that will use
// {@link mashschema.MashDataSourceSearch} constant as a field.
func DataSourceDefaultQueryExtractor(d *schema.ResourceData) (map[string]string, error) {
	return mashschema.ExtractStringMap(d, mashschema.MashDataSourceSearch), nil
}

func (t *DatasourceTemplate) TemplateQuery(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	if query, queryError := t.actualQuery(d); queryError != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("invalid query for %s", t.Mapper.V3ObjectName()),
			Detail:   fmt.Sprintf("qeurying %s is not possible: %s", t.Mapper.V3ObjectName(), queryError.Error()),
		}}
	} else {
		if rv, err := t.Query(ctx, v3cl, query); err != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("querying %s failed", t.Mapper.V3ObjectName()),
				Detail:   fmt.Sprintf("querying %ss using parameters %s has failed: %s", t.Mapper.V3ObjectName(), query, err.Error()),
			}}
		} else {
			switch {
			case t.isMatchRequired(d) && len(rv) == 0:
				return diag.Diagnostics{diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("no %ss mached the query", t.Mapper.V3ObjectName()),
					Detail:   fmt.Sprintf("querying %ss using parameters %s has retuned no matches where results are expected", t.Mapper.V3ObjectName(), query),
				}}
			case t.isUniqueMatchRequired() && len(rv) != 1:
				return diag.Diagnostics{diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("multiple %ss mached the query", t.Mapper.V3ObjectName()),
					Detail:   fmt.Sprintf("querying %ss using parameters %s has retuned %d matches were exactly one required", t.Mapper.V3ObjectName(), query, len(rv)),
				}}
			case t.isUniqueMatchRequired() && len(rv) == 1:
				return t.Mapper.SetStateOf(rv[0], d)
			case t.isMatchRequired(d) && len(rv) > 0:
				return t.Mapper.SetState(rv, d)
			case !t.isMatchRequired(d) && len(rv) == 0:
				d.SetId("")
				return nil
			default:
				return diag.Diagnostics{diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("undefined query result of %s", t.Mapper.V3ObjectName()),
					Detail:   fmt.Sprintf("querying %s returned data set that matches non supported conditions to persis in Terraform schema", t.Mapper.V3ObjectName()),
				}}
			}
		}
	}
}
