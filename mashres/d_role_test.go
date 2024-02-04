package mashres

import (
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"testing"
	"time"
)

func TestQueryingRolesWillBeSavedOnUniqueMatch(t *testing.T) {
	params := map[string]string{"name": "abc"}
	h := CreateTestDatasource[tfmapper.Orphan, string, masherytypes.Role](RoleDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenSingleRoleMatchIsReturned(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.willHaveStateId(t, "roleId")
}

func TestQueryingRoleWhereNoMatchReturnedWithRequiredWillGiveDiagnostics(t *testing.T) {
	params := map[string]string{"name": "abc"}
	h := CreateTestDatasource[tfmapper.Orphan, string, masherytypes.Role](RoleDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	h.givenStateFieldSetTo(t, mashschema.MashDataSourceRequired, true)
	givenNoRoleMatchIsReturned(h, params)

	h.thenExecutingDataSourceQueryWillYieldDiagnostic(t, "no matching object was found, however the configuration requires a match")
}

func TestQueryingRoleWhereAPICallFWillGiveDiagnostics(t *testing.T) {
	params := map[string]string{"name": "abc"}
	h := CreateTestDatasource[tfmapper.Orphan, string, masherytypes.Role](RoleDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenErrorReturnedWhenQueryingRoles(h, params)

	h.thenExecutingDataSourceQueryWillYieldDiagnostic(t, "query has returned an error: sample rejection")
}

func TestQueryingRoleWithOptionalFlagWitchNoMatchesSucceeds(t *testing.T) {
	params := map[string]string{"name": "abc"}
	h := CreateTestDatasource[tfmapper.Orphan, string, masherytypes.Role](RoleDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	h.givenStateFieldSetTo(t, mashschema.MashDataSourceRequired, false)
	givenNoRoleMatchIsReturned(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.willHaveStateId(t, "")
	h.willHaveFieldSetTo(t, mashschema.MashObjName, "")
}

func givenSingleRoleMatchIsReturned(h *DatasourceTemplateMockHelper[tfmapper.Orphan, string, masherytypes.Role], params map[string]string) {
	mashTime := masherytypes.MasheryJSONTime(time.Now())

	returnedRoles := []masherytypes.Role{
		{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "roleId",
				Name:      params["name"],
				Created:   &mashTime,
				Updated:   &mashTime,
				Retrieved: time.Now(),
			},
			Description: "Desc",
			Predefined:  false,
			OrgRole:     false,
			Assignable:  true,
		},
	}
	h.mockClientWill().
		On("ListRolesFiltered", mock.Anything, params).
		Return(returnedRoles, nil)
}
func givenNoRoleMatchIsReturned(h *DatasourceTemplateMockHelper[tfmapper.Orphan, string, masherytypes.Role], params map[string]string) {
	h.mockClientWill().
		On("ListRolesFiltered", mock.Anything, params).
		Return([]masherytypes.Role{}, nil)
}

func givenErrorReturnedWhenQueryingRoles(h *DatasourceTemplateMockHelper[tfmapper.Orphan, string, masherytypes.Role], params map[string]string) {
	h.mockClientWill().
		On("ListRolesFiltered", mock.Anything, params).
		Return([]masherytypes.Role{}, errors.New("sample rejection"))
}
