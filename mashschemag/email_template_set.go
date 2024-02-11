package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var EmailTemplateSetResourceSchemaBuilder = tfmapper.NewSchemaBuilder[tfmapper.Orphan, string, masherytypes.EmailTemplateSet]().
	Identity(&tfmapper.StringIdentityFieldMapper[masherytypes.EmailTemplateSet]{
		Locator: func(in *masherytypes.EmailTemplateSet) *string {
			return &in.Id
		},
	}).
	RootIdentity(&tfmapper.RootParentIdentity{})

// Add created and updated.
func init() {
	EmailTemplateSetResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.EmailTemplateSet]{
		Locator: func(in *masherytypes.EmailTemplateSet) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.EmailTemplateSet]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.EmailTemplateSet]{
		Locator: func(in *masherytypes.EmailTemplateSet) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.EmailTemplateSet]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	})
}

func init() {
	EmailTemplateSetResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.EmailTemplateSet]{
		Locator: func(in *masherytypes.EmailTemplateSet) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.EmailTemplateSet]{
			Key: mashschema.MashObjName,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email data set name",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.EmailTemplateSet]{
		Locator: func(in *masherytypes.EmailTemplateSet) *string {
			return &in.Type
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.EmailTemplateSet]{
			Key: mashschema.MashEmailTemplateSetType,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email template set type",
			},
		},
	})
}
