package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mashschema "terraform-provider-mashery/mashschema"
)

func dataSourceMasheryRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: readRole,
		Schema:      mashschema.RoleMapper.TerraformSchema(),
	}
}

func readRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	query := mashschema.RoleMapper.GetQuery(d)

	if rv, err := v3cl.ListRolesFiltered(ctx, query, mashschema.EmptyStringArray); err != nil {
		return diag.FromErr(err)
	} else if rv != nil && len(rv) > 0 {
		mashschema.RoleMapper.PersistTyped(ctx, &rv[0], d)
		d.SetId(rv[0].Id)
	} else {
		d.SetId("")
		if mashschema.RoleMapper.IsMatchRequired(d) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "required role was not found",
				Detail:   "no role matching this query exists",
			}}
		}
	}

	return nil
}
