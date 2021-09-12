package mashery_test

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestV3MashRoleToResourceData(t *testing.T) {
	d := NewResourceData(&mashery.DataSourceRoleSchema)

	now := v3client.MasheryJSONTime(time.Now())
	orig := v3client.MasheryRole{
		AddressableV3Object: v3client.AddressableV3Object{
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

	setDiag := mashery.V3MashRoleToResourceData(&orig, d)
	LogErrorDiagnostics(t, "setting mashery role to resource staet", &setDiag)

	name := d.Get(mashery.MashObjName).(string)
	desc := d.Get(mashery.MashObjDescription).(string)
	predefined := d.Get(mashery.MashRolePredefined).(bool)
	orgRole := d.Get(mashery.MashRoleOrgRole).(bool)
	assignable := d.Get(mashery.MashRoleAssignableRole).(bool)

	assert.Equal(t, orig.Name, name)
	assert.Equal(t, orig.Description, desc)
	assert.Equal(t, orig.Predefined, predefined)
	assert.Equal(t, orig.OrgRole, orgRole)
	assert.Equal(t, orig.Assignable, assignable)

	roleRef := d.Get(mashery.MashReadRolePermission).(map[string]interface{})
	fmt.Println(roleRef)
	assert.Equal(t, "roleId", roleRef[mashery.MashObjId])
	assert.Equal(t, "Name", roleRef[mashery.MashObjName])
	assert.Equal(t, "read", roleRef[mashery.MashRoleAction])
}
