package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func dataSourceMasheryEmailTemplateSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: readEmailTemplateSet,
		Schema:      mashschema.EmailTemplateSetMapper.TerraformSchema(),
	}
}

func readEmailTemplateSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	query := mashschema.EmailTemplateSetMapper.GetQuery(d)

	if rv, err := v3cl.ListEmailTemplateSetsFiltered(ctx, query, mashschema.EmptyStringArray); err != nil {
		return diag.FromErr(err)
	} else if rv != nil && len(rv) > 0 {
		mashschema.EmailTemplateSetMapper.PersistTyped(ctx, &rv[0], d)
		d.SetId(rv[0].Id)
	} else {
		d.SetId("")
		if mashschema.EmailTemplateSetMapper.IsMatchRequired(d) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing required V3 object",
				Detail:   "No email template matching this query exists",
			}}
		}
	}

	return nil
}
