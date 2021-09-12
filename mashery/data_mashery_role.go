package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMasheryRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: readRole,
		Schema:      DataSourceRoleSchema,
	}
}

func readRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	query := extractStringMap(d, MashDataSourceSearch)

	if rv, err := v3cl.ListRolesFiltered(ctx, query, EmptyStringArray); err != nil {
		return diag.FromErr(err)
	} else if rv != nil && len(rv) > 0 {
		V3MashRoleToResourceData(&rv[0], d)
		d.SetId(rv[0].Id)
	} else {
		d.SetId("")
		if extractBool(d, MashDataSourceRequired, true) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing required V3 object", // TODO: this is constant and needs to be moved to the file.
				Detail:   "No role matching this query exists",
			}}
		}
	}

	return diag.Diagnostics{}
}
