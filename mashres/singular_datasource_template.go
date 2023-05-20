package mashres

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

type QueryFunc[MType any] func(context.Context, v3client.Client, map[string]string) (*MType, error)

type SingularDatasourceTemplate[ParentIdent any, Ident any, MType any] struct {
	Schema map[string]*schema.Schema
	Mapper *tfmapper.Mapper[ParentIdent, Ident, MType]

	DoQuery QueryFunc[MType]
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) DataSourceSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: sdt.Query,
		Schema:      sdt.Schema,
	}
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) isMatchRequired(d *schema.ResourceData) bool {
	return mashschema.ExtractBool(d, mashschema.MashDataSourceRequired, true)
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) Query(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Cl := m.(v3client.Client)
	query := mashschema.ExtractStringMap(d, mashschema.MashDataSourceSearch)

	if obj, err := sdt.DoQuery(ctx, v3Cl, query); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("query has returned an error: %s", err.Error()),
		}}
	} else if obj == nil && sdt.isMatchRequired(d) {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "no matching object were found, however the configuration requires a match",
		}}
	} else {
		return sdt.Mapper.RemoteToSchema(obj, d)
	}
}

func CreateSingularDataSource[ParentIdent any, Ident any, MType any](builder *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], queryFunc QueryFunc[MType]) *SingularDatasourceTemplate[ParentIdent, Ident, MType] {
	mapperSchema := mashschema.CloneAsComputed(builder.ResourceSchema())
	mapperSchema[mashschema.MashDataSourceSearch] = &schema.Schema{
		Type:        schema.TypeMap,
		MinItems:    1,
		Required:    true,
		Description: "Search conditions for this email set, typically name = value",
		Elem:        mashschema.StringElem(),
	}

	mapperSchema[mashschema.MashDataSourceRequired] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "If true (default), then email template set must exist. If an element doesn't exist, the error is generated",
	}

	rv := SingularDatasourceTemplate[ParentIdent, Ident, MType]{
		Schema:  mapperSchema,
		Mapper:  builder.Mapper(),
		DoQuery: queryFunc,
	}

	return &rv
}
