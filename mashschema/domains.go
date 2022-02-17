package mashschema

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
)

const MashDomains = "domains"

var DomainsMapper *DomainsMapperImpl

type DomainsMapperImpl struct {
	MapperImpl
}

func (dmi *DomainsMapperImpl) PersistTyped(ctx context.Context, domains []string, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashDomains: domains,
	}

	return dmi.SetResourceFields(ctx, data, d)
}

func init() {
	DomainsMapper = &DomainsMapperImpl{
		MapperImpl{
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
		},
	}

	DomainsMapper.persistFunc = func(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics {
		return DomainsMapper.PersistTyped(ctx, rv.([]string), d)
	}
}
