package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/mashschemag"
	"testing"
)

func TestQueryingApplicationsWillSucceedOnUniqueMatch(t *testing.T) {
	params := map[string]string{"name": "app-name"}
	h := CreateTestDatasource[masherytypes.MemberIdentifier,
		mashschemag.ApplicationOfMemberIdentifier,
		masherytypes.Application](ApplicationDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenListApplicationsFilteredSucceeds(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id mashschemag.ApplicationOfMemberIdentifier) {
		assert.Equal(t, "app-id", id.ApplicationId)
	})
}

func givenListApplicationsFilteredSucceeds(h *DatasourceTemplateMockHelper[masherytypes.MemberIdentifier,
	mashschemag.ApplicationOfMemberIdentifier,
	masherytypes.Application], qs map[string]string) {
	rv := masherytypes.Application{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "app-id",
			Name: "app-name",
		},
		Username: "app-user-name",
	}
	h.mockClientWill().
		On("ListApplicationsFiltered", mock.Anything, qs).
		Return([]masherytypes.Application{rv}, nil).
		Once()
}
