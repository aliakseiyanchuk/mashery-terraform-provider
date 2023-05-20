package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var RoleResourceSchemaBuilder = tfmapper.NewSchemaBuilder[tfmapper.Orphan, string, masherytypes.Role]().
	Identity(&tfmapper.StringIdentityFieldMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *string {
			return &in.Id
		},
	}).
	RootIdentity(&tfmapper.RootParentIdentity{})

// Add created and updated.
func init() {
	RoleResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
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
	RoleResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
			Key: mashschema.MashObjName,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role name",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *string {
			return &in.Description
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
			Key: mashschema.MashSvcDescription, // TODO make a common constant
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role description",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *bool {
			return &in.Assignable
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
			Key: mashschema.MashRoleAssignableRole,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not this role is assignable",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *bool {
			return &in.OrgRole
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
			Key: mashschema.MashRoleOrgRole,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not this role is an org role",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Role]{
		Locator: func(in *masherytypes.Role) *bool {
			return &in.Predefined
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Role]{
			Key: mashschema.MashRolePredefined,
			// Perhaps simple strings can be optimized.
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this role is pre-defined",
			},
		},
	})
}
