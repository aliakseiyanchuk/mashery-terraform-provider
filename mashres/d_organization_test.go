package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestQueryingEOrganizationsWillSucceedOnUniqueMatch(t *testing.T) {
	params := map[string]string{"name": "org-name"}
	h := CreateTestDatasource(OrganizationDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenListOrganizationsFilteredSucceeds(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.willHaveStateId(t, "org-id")
}

func givenListOrganizationsFilteredSucceeds(h *DatasourceTemplateMockHelper, qs map[string]string) {
	rv := masherytypes.Organization{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "org-id",
			Name: "org-name",
		},
	}
	h.mockClientWill().
		On("ListOrganizationsFiltered", mock.Anything, qs).
		Return([]masherytypes.Organization{rv}, nil).
		Once()
}
