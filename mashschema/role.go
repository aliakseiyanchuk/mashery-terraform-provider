package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// schema definition for Mashery role.

const (
	MashRolePredefined     = "predefined_role"
	MashRoleOrgRole        = "org_role"
	MashRoleAssignableRole = "assignable_role"
	MashReadRolePermission = "read_permission"

	MashRoleAction = "action"
)

var RoleMapper *RoleMapperImpl

type RoleMapperImpl struct {
	DataSourceMapperImpl
}

func (rm *RoleMapperImpl) PersistTyped(inp masherytypes.Role, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashObjCreated:         inp.Created.ToString(),
		MashObjUpdated:         inp.Updated.ToString(),
		MashObjName:            inp.Name,
		MashObjDescription:     inp.Description,
		MashRolePredefined:     inp.Predefined,
		MashRoleOrgRole:        inp.OrgRole,
		MashRoleAssignableRole: inp.Assignable,
		MashReadRolePermission: rm.rolePermission(&inp, "read"),
	}

	return SetResourceFields(data, d)
}

func (rm *RoleMapperImpl) rolePermission(inp *masherytypes.Role, perm string) map[string]interface{} {
	return map[string]interface{}{
		MashObjId: inp.Id,
		//MashObjName:    inp.Name,
		MashRoleAction: perm,
	}
}

func (rm *RoleMapperImpl) initSchemaBoilerplate() {
	rm.SchemaBuilder().
		AddComputedString(MashObjCreated, "Date/time object was created").
		AddComputedString(MashObjUpdated, "Date/time object was updated").
		AddComputedString(MashObjName, "Role name").
		AddOptionalString(MashObjDescription, "Role description").
		AddComputedBoolean(MashRolePredefined, "Whether role is pre-defined").
		AddComputedBoolean(MashRoleOrgRole, "Whether this role is an org-role").
		AddComputedBoolean(MashRoleAssignableRole, "Whether this role is assignable")
}

func init() {
	RoleMapper = &RoleMapperImpl{
		DataSourceMapperImpl: DataSourceMapperImpl{
			v3ObjectName: "role",
			schema: map[string]*schema.Schema{
				MashDataSourceSearch: {
					Type:        schema.TypeMap,
					Required:    true,
					Description: "Search conditions for this role, typically name = value",
					Elem:        stringElem(),
				},
				MashDataSourceRequired: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "If true (default), then role satisfying the search condition must exist. If an element doesn't exist, the error is generated",
				},
				MashReadRolePermission: {
					Type:        schema.TypeMap,
					Computed:    true,
					Description: "Role reference object usable in service I/O docs and package plan portal access",
				},
			},

			persistOne: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				return RoleMapper.PersistTyped(rv.(masherytypes.Role), d)
			},
		},
	}

	RoleMapper.initSchemaBoilerplate()
}
