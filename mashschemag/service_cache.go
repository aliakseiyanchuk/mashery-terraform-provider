package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"time"
)

var ServiceCacheResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
	})

// Parent package identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		Key: mashschema.MashSvcId,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Service Id, to which this plan belongs",
		},
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ServiceIdentifier) bool {
			return len(inp.ServiceId) > 0
		},
	}

	ServiceCacheResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	ServiceCacheResourceSchemaBuilder.Add(&tfmapper.DurationFieldMapper[masherytypes.ServiceCache]{
		Locator: func(in *masherytypes.ServiceCache) *int64 {
			return &in.CacheTtl
		},
		Unit: time.Minute,
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceCache]{
			Key: mashschema.MashSvcCacheTtl,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Description:      "Time till which the data is stored in cache",
				ValidateDiagFunc: mashschema.ValidateDuration,
				Required:         true,
			},
		},
	})
}
