package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema definition for Mashery role.

const (
	MashRolePredefined     = "predefined_role"
	MashRoleOrgRole        = "org_role"
	MashRoleAssignableRole = "assignable_role"
	MashReadRolePermission = "read_permission"

	MashRoleAction = "action"
)

func V3MashRoleToResourceData(inp *masherytypes.MasheryRole, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashObjCreated:         inp.Created.ToString(),
		MashObjUpdated:         inp.Updated.ToString(),
		MashObjName:            inp.Name,
		MashObjDescription:     inp.Description,
		MashRolePredefined:     inp.Predefined,
		MashRoleOrgRole:        inp.OrgRole,
		MashRoleAssignableRole: inp.Assignable,
		MashReadRolePermission: V3RoleRefenceWithPermissionToTerraform(inp, "read"),
	}

	return SetResourceFields(data, d)
}

func V3RoleRefenceWithPermissionToTerraform(inp *masherytypes.MasheryRole, perm string) map[string]interface{} {
	return map[string]interface{}{
		MashObjId:      inp.Id,
		MashRoleAction: perm,
	}
}

func V3MashRolePermissionToTerraform(permission masherytypes.MasheryRolePermission) map[string]interface{} {
	return map[string]interface{}{
		MashObjId:      permission.Id,
		MashRoleAction: permission.Action,
	}
}

func V3RolesPermissionsToTerraform(inp []masherytypes.MasheryRolePermission) []map[string]interface{} {
	rv := make([]map[string]interface{}, len(inp))
	for idx, v := range inp {
		rv[idx] = V3MashRolePermissionToTerraform(v)
	}

	return rv
}

func V3RolePermissionUpsertable(inp map[string]interface{}) masherytypes.MasheryRolePermission {
	return masherytypes.MasheryRolePermission{
		MasheryRole: masherytypes.MasheryRole{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: inp[MashObjId].(string),
			},
		},
		Action: inp[MashRoleAction].(string),
	}
}

var RolePermissionReferenceSchema = map[string]*schema.Schema{}
var DataSourceRoleSchema = map[string]*schema.Schema{
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
}

func initRoleSchemaBoilerplate() {
	addComputedString(&DataSourceRoleSchema, MashObjCreated, "Date/time object was created")
	addComputedString(&DataSourceRoleSchema, MashObjUpdated, "Date/time object was updated")
	addComputedString(&DataSourceRoleSchema, MashObjName, "Role name")
	addOptionalString(&DataSourceRoleSchema, MashObjDescription, "Role description")

	addComputedBoolean(&DataSourceRoleSchema, MashRolePredefined, "Whether role is pre-defined")
	addComputedBoolean(&DataSourceRoleSchema, MashRoleOrgRole, "Whether this role is an org-role")
	addComputedBoolean(&DataSourceRoleSchema, MashRoleAssignableRole, "Whether this role is assignable")
}

func initRoleReference() {
	addRequiredString(&RolePermissionReferenceSchema, MashObjId, "Role Id")
	addRequiredString(&RolePermissionReferenceSchema, MashRoleAction, "Action granted to this role")
}

func init() {
	initRoleSchemaBoilerplate()
	initRoleReference()
}
