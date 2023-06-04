package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var ServiceEndpointMethodResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ServiceEndpointMethodIdentifier]{
		IdentityFunc: func() masherytypes.ServiceEndpointMethodIdentifier {
			return masherytypes.ServiceEndpointMethodIdentifier{}
		},
	})

// Parent service endpoint identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ServiceEndpointIdentifier]{
		Key: mashschema.MashEndpointId,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Endpoint Id, to which this method belongs",
		},
		IdentityFunc: func() masherytypes.ServiceEndpointIdentifier {
			return masherytypes.ServiceEndpointIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ServiceEndpointIdentifier) bool {
			return len(inp.ServiceId) > 0 && len(inp.EndpointId) > 0
		},
	}

	ServiceEndpointMethodResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

// Read-only properties
func init() {
	ServiceEndpointMethodResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethod]{
		Locator: func(in *masherytypes.ServiceEndpointMethod) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethod]{
			Key: mashschema.MashServiceEndpointMethodId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 identifier of this method",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ServiceEndpointMethod]{
		Locator: func(in *masherytypes.ServiceEndpointMethod) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethod]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ServiceEndpointMethod]{
		Locator: func(in *masherytypes.ServiceEndpointMethod) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethod]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	})
}

// Fields
func init() {
	ServiceEndpointMethodResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethod]{
		Locator: func(in *masherytypes.ServiceEndpointMethod) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethod]{
			Key: mashschema.MashObjName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Method name, as it would be detected by Mashery",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethod]{
		Locator: func(in *masherytypes.ServiceEndpointMethod) *string {
			return &in.SampleJsonResponse
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethod]{
			Key: mashschema.MashServiceEndpointMethodSampleJson,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sample JSON response",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethod]{
		Locator: func(in *masherytypes.ServiceEndpointMethod) *string {
			return &in.SampleXmlResponse
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethod]{
			Key: mashschema.MashServiceEndpointMethodSampleXml,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sample JSON response",
			},
		},
	})
}
