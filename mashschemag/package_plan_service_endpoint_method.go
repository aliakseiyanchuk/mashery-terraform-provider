package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

type PackagePlanServiceEndpointMethodParam struct {
	ServiceEndpointMethod               masherytypes.ServiceEndpointMethodIdentifier
	ServiceEndpointMethodFilterDesired  masherytypes.ServiceEndpointMethodFilterIdentifier
	ServiceEndpointMethodFilterPrevious masherytypes.ServiceEndpointMethodFilterIdentifier
	ServiceEndpointMethodFilterChanged  bool
}

var PackagePlanServiceEndpointMethodResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.PackagePlanServiceEndpointIdentifier,
	masherytypes.PackagePlanServiceEndpointMethodIdentifier,
	PackagePlanServiceEndpointMethodParam]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.PackagePlanServiceEndpointMethodIdentifier]{
		IdentityFunc: func() masherytypes.PackagePlanServiceEndpointMethodIdentifier {
			return masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
		},
	})

// Parent package service identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.PackagePlanServiceEndpointIdentifier]{
		Key: mashschema.PackagePlanEndpointRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Package plan service Id, to which this endpoint belongs",
		},
		IdentityFunc: func() masherytypes.PackagePlanServiceEndpointIdentifier {
			return masherytypes.PackagePlanServiceEndpointIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.PackagePlanServiceEndpointIdentifier) bool {
			return len(inp.PackageId) > 0 &&
				len(inp.PlanId) > 0 &&
				len(inp.ServiceId) > 0 &&
				len(inp.EndpointId) > 0
		},
	}

	PackagePlanServiceEndpointMethodResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	mapper := tfmapper.JsonIdentityFieldMapper[masherytypes.ServiceEndpointMethodIdentifier, PackagePlanServiceEndpointMethodParam]{
		FieldMapperBase: tfmapper.FieldMapperBase[PackagePlanServiceEndpointMethodParam]{
			Key: mashschema.ServiceEndpointMethodRef,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service endpoint method to include in the plan",
			},
		},
		IdentityFunc: func() masherytypes.ServiceEndpointMethodIdentifier {
			return masherytypes.ServiceEndpointMethodIdentifier{}
		},
		ValidateIdentFunc: func(identifier masherytypes.ServiceEndpointMethodIdentifier) bool {
			return len(identifier.ServiceId) > 0 && len(identifier.EndpointId) > 0 &&
				len(identifier.MethodId) > 0
		},
		Locator: func(in *PackagePlanServiceEndpointMethodParam) *masherytypes.ServiceEndpointMethodIdentifier {
			return &in.ServiceEndpointMethod
		},
	}

	PackagePlanServiceEndpointMethodResourceSchemaBuilder.Add(mapper.PrepareMapper())
}

func init() {
	mapper := tfmapper.JsonIdentityFieldMapper[masherytypes.ServiceEndpointMethodFilterIdentifier, PackagePlanServiceEndpointMethodParam]{
		FieldMapperBase: tfmapper.FieldMapperBase[PackagePlanServiceEndpointMethodParam]{
			Key: mashschema.ServiceEndpointMethodFilterRef,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service endpoint method filter to include in the plan",
			},
			ModificationConsumer: func(a *PackagePlanServiceEndpointMethodParam, b bool) {
				a.ServiceEndpointMethodFilterChanged = b
			},
		},
		// This mapper is writeable. In case the filter would be changed e.g. on the console, the script
		// should be able to detect it.
		NullFunction: func(in *PackagePlanServiceEndpointMethodParam) bool {
			return in != nil &&
				len(in.ServiceEndpointMethodFilterDesired.FilterId) > 0 &&
				len(in.ServiceEndpointMethodFilterDesired.MethodId) > 0 &&
				len(in.ServiceEndpointMethodFilterDesired.EndpointId) > 0 &&
				len(in.ServiceEndpointMethodFilterDesired.ServiceId) > 0
		},
		PreviousLocator: func(in *PackagePlanServiceEndpointMethodParam) *masherytypes.ServiceEndpointMethodFilterIdentifier {
			return &in.ServiceEndpointMethodFilterPrevious
		},
		IdentityFunc: func() masherytypes.ServiceEndpointMethodFilterIdentifier {
			return masherytypes.ServiceEndpointMethodFilterIdentifier{}
		},
		ValidateIdentFunc: func(identifier masherytypes.ServiceEndpointMethodFilterIdentifier) bool {
			return len(identifier.ServiceId) > 0 && len(identifier.EndpointId) > 0
		},
		Locator: func(in *PackagePlanServiceEndpointMethodParam) *masherytypes.ServiceEndpointMethodFilterIdentifier {
			return &in.ServiceEndpointMethodFilterDesired
		},
	}

	PackagePlanServiceEndpointMethodResourceSchemaBuilder.Add(mapper.PrepareMapper())
}
