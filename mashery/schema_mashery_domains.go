package mashery

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MashDomains = "domains"

var DomainsSchema = map[string]*schema.Schema{
	MashDomains: {
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "Domains associated with your area",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

func V3DomainsToResourceData(domains []string, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashDomains: domains,
	}

	return SetResourceFields(data, d)
}
