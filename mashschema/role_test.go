package mashschema_test

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestV3MashRoleToResourceData(t *testing.T) {
	d := mashschema.RoleMapper.TestResourceData()

	now := masherytypes.MasheryJSONTime(time.Now())
	orig := masherytypes.Role{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "roleId",
			Name:    "Name",
			Created: &now,
			Updated: &now,
		},
		Description: "Desc",
		Predefined:  true,
		OrgRole:     true,
		Assignable:  true,
	}

	setDiag := mashschema.RoleMapper.SetStateOf(orig, d)
	LogErrorDiagnostics(t, "setting mashery role to resource state", &setDiag)

	name := d.Get(mashschema.MashObjName).(string)
	desc := d.Get(mashschema.MashObjDescription).(string)
	predefined := d.Get(mashschema.MashRolePredefined).(bool)
	orgRole := d.Get(mashschema.MashRoleOrgRole).(bool)
	assignable := d.Get(mashschema.MashRoleAssignableRole).(bool)

	assert.Equal(t, orig.Name, name)
	assert.Equal(t, orig.Description, desc)
	assert.Equal(t, orig.Predefined, predefined)
	assert.Equal(t, orig.OrgRole, orgRole)
	assert.Equal(t, orig.Assignable, assignable)

	roleRef := d.Get(mashschema.MashReadRolePermission).(map[string]interface{})
	fmt.Println(roleRef)
	assert.Equal(t, "roleId", roleRef[mashschema.MashObjId])
	assert.Equal(t, "Name", roleRef[mashschema.MashObjName])
	assert.Equal(t, "read", roleRef[mashschema.MashRoleAction])
}
