package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var OrganizationResourceSchemaBuilder = tfmapper.NewSchemaBuilder[tfmapper.Orphan, string, masherytypes.Organization]().
	Identity(&tfmapper.StringIdentityFieldMapper[masherytypes.Organization]{
		Locator: func(in *masherytypes.Organization) *string {
			return &in.Id
		},
	}).
	RootIdentity(&tfmapper.RootParentIdentity{})

// Add created and updated.
func init() {
	OrganizationResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.Organization]{
		Locator: func(in *masherytypes.Organization) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Organization]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Organization]{
		Locator: func(in *masherytypes.Organization) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Organization]{
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
	OrganizationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Organization]{
		Locator: func(in *masherytypes.Organization) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Organization]{
			Key: mashschema.MashObjName,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email data set name",
			},
		},
	})
}
