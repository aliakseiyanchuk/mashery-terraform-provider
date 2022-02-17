package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func dataSourceMasheryPublicDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: readPublicDomains,
		Schema:      mashschema.DomainsMapper.TerraformSchema(),
	}
}

func readPublicDomains(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetPublicDomains(ctx); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId("public_domains")
		doLogf("received %d public domains", len(rv))
		return mashschema.DomainsMapper.PersistTyped(ctx, rv, d)
	}
}
