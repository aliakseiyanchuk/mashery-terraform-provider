package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

type PackagePlanServiceEndpointParam struct {
	ServiceEndpointIdentifier masherytypes.ServiceEndpointIdentifier
}

var PackagePlanServiceEndpointResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.PackagePlanServiceIdentifier,
	masherytypes.PackagePlanServiceEndpointIdentifier,
	PackagePlanServiceEndpointParam]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.PackagePlanServiceEndpointIdentifier]{
		IdentityFunc: func() masherytypes.PackagePlanServiceEndpointIdentifier {
			return masherytypes.PackagePlanServiceEndpointIdentifier{}
		},
	})

// Parent package service identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.PackagePlanServiceIdentifier]{
		Key: mashschema.MashPlanServiceId,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Package plan service Id, to which this endpoint belongs",
		},
		IdentityFunc: func() masherytypes.PackagePlanServiceIdentifier {
			return masherytypes.PackagePlanServiceIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.PackagePlanServiceIdentifier) bool {
			return len(inp.PackageId) > 0 && len(inp.PlanId) > 0 && len(inp.ServiceId) > 0
		},
	}

	PackagePlanServiceEndpointResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	mapper := tfmapper.JsonIdentityFieldMapper[masherytypes.ServiceEndpointIdentifier, PackagePlanServiceEndpointParam]{
		FieldMapperBase: tfmapper.FieldMapperBase[PackagePlanServiceEndpointParam]{
			Key: mashschema.MashEndpointId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service endpoint to expose in this plan",
			},
		},
		IdentityFunc: func() masherytypes.ServiceEndpointIdentifier {
			return masherytypes.ServiceEndpointIdentifier{}
		},
		ValidateIdentFunc: func(identifier masherytypes.ServiceEndpointIdentifier) bool {
			return len(identifier.ServiceId) > 0 && len(identifier.EndpointId) > 0
		},
		Locator: func(in *PackagePlanServiceEndpointParam) *masherytypes.ServiceEndpointIdentifier {
			return &in.ServiceEndpointIdentifier
		},
	}

	PackagePlanServiceEndpointResourceSchemaBuilder.Add(mapper.PrepareMapper())
}
