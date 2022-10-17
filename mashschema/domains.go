package mashschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
)

const MashDomains = "domains"

var DomainsMapper *DomainsMapperImpl

type DomainsMapperImpl struct {
	DataSourceMapperImpl
}

func (dmi *DomainsMapperImpl) PersistTyped(domains []string, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashDomains: domains,
	}

	return SetResourceFields(data, d)
}

func init() {
	DomainsMapper = &DomainsMapperImpl{
		DataSourceMapperImpl{
			schema: map[string]*schema.Schema{
				MashDomains: {
					Type:        schema.TypeSet,
					Computed:    true,
					Description: "Domains associated with your area",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Set: func(i interface{}) int {
						return hashcode.String(i.(string))
					},
				},
			},

			persistMany: func(rv []interface{}, d *schema.ResourceData) diag.Diagnostics {
				return DomainsMapper.PersistTyped(CoerceInterfaceArrayToStringArray(rv), d)
			},
		},
	}
}
