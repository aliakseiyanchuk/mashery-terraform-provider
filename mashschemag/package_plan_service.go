package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

type PackagePlanServiceParam struct {
	ServiceIdentifier masherytypes.ServiceIdentifier
}

var PackagePlanServiceResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.PackagePlanIdentifier,
	masherytypes.PackagePlanServiceIdentifier,
	PackagePlanServiceParam]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.PackagePlanServiceIdentifier]{
		IdentityFunc: func() masherytypes.PackagePlanServiceIdentifier {
			return masherytypes.PackagePlanServiceIdentifier{}
		},
	})

// Parent package identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.PackagePlanIdentifier]{
		Key: mashschema.MashPlanRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Package plan reference, to which this service belongs",
		},
		IdentityFunc: func() masherytypes.PackagePlanIdentifier {
			return masherytypes.PackagePlanIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.PackagePlanIdentifier) bool {
			return len(inp.PackageId) > 0 && len(inp.PlanId) > 0
		},
	}

	PackagePlanServiceResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	mapper := tfmapper.JsonIdentityFieldMapper[masherytypes.ServiceIdentifier, PackagePlanServiceParam]{
		FieldMapperBase: tfmapper.FieldMapperBase[PackagePlanServiceParam]{
			Key: mashschema.MashSvcRef,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service to expose in this plan",
			},
		},
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
		ValidateIdentFunc: func(identifier masherytypes.ServiceIdentifier) bool {
			return len(identifier.ServiceId) > 0
		},
		Locator: func(in *PackagePlanServiceParam) *masherytypes.ServiceIdentifier {
			return &in.ServiceIdentifier
		},
	}

	PackagePlanServiceResourceSchemaBuilder.Add(mapper.PrepareMapper())
}
