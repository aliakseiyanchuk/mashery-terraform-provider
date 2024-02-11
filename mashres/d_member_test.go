package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"testing"
)

func TestQueryingMembersWillSucceedOnUniqueMatch(t *testing.T) {
	params := map[string]string{"name": "member-name"}
	h := CreateTestDatasource[tfmapper.Orphan, masherytypes.MemberIdentifier, masherytypes.Member](MemberDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenListMembersFilteredSucceeds(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.MemberIdentifier) {
		assert.Equal(t, "member-id", id.MemberId)
		assert.Equal(t, "member-user-name", id.Username)
	})
}

func givenListMembersFilteredSucceeds(h *DatasourceTemplateMockHelper[tfmapper.Orphan, masherytypes.MemberIdentifier, masherytypes.Member], qs map[string]string) {
	rv := masherytypes.Member{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "member-id",
			Name: "member-name",
		},
		Username: "member-user-name",
	}
	h.mockClientWill().
		On("ListMembersFiltered", mock.Anything, qs).
		Return([]masherytypes.Member{rv}, nil).
		Once()
}
