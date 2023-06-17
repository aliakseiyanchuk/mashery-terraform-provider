package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var ServiceEndpointMethodFilterResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ServiceEndpointMethodFilterIdentifier]{
		IdentityFunc: func() masherytypes.ServiceEndpointMethodFilterIdentifier {
			return masherytypes.ServiceEndpointMethodFilterIdentifier{}
		},
	})

// Parent service endpoint method identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ServiceEndpointMethodIdentifier]{
		Key: mashschema.ServiceEndpointMethodRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Method reference, to which this filter belongs",
		},
		IdentityFunc: func() masherytypes.ServiceEndpointMethodIdentifier {
			return masherytypes.ServiceEndpointMethodIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ServiceEndpointMethodIdentifier) bool {
			return len(inp.ServiceId) > 0 && len(inp.EndpointId) > 0
		},
	}

	ServiceEndpointMethodFilterResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

// Read-only properties
func init() {
	ServiceEndpointMethodFilterResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
			Key: mashschema.MashSvcEndpointMethodFilterId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 identifier of this filter",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
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
	ServiceEndpointMethodFilterResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
			Key: mashschema.MashObjName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Filter name",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *string {
			return &in.Notes
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
			Key: mashschema.MashServiceEndpointMethodFilterNotes,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *string {
			return &in.XmlFilterFields
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
			Key: mashschema.MashServiceEndpointMethodFilterXmlFields,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "XML sample response",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ServiceEndpointMethodFilter]{
		Locator: func(in *masherytypes.ServiceEndpointMethodFilter) *string {
			return &in.JsonFilterFields
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ServiceEndpointMethodFilter]{
			Key: mashschema.MashServiceEndpointMethodFilterJsonFields,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON sample response",
			},
		},
	})
}
