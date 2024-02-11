package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestCreateServiceEndpointFilterMethodWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceEndpointMethodFilterResource)

	serviceIdent := masherytypes.ServiceEndpointMethodIdentifier{}
	serviceIdent.ServiceId = "serviceId"
	serviceIdent.EndpointId = "endpointId"
	serviceIdent.MethodId = "methId"

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashObjName, "sample-method-filter")

	givenCreatingEndpointMethodFilterSucceeds(h, serviceIdent)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceEndpointMethodFilterIdentifier) {
		assert.Equal(t, "serviceId", id.ServiceId)
		assert.Equal(t, "endpointId", id.EndpointId)
		assert.Equal(t, "methId", id.MethodId)
		assert.Equal(t, "filterId", id.FilterId)
	})
}

func givenCreatingEndpointMethodFilterSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter],
	ident masherytypes.ServiceEndpointMethodIdentifier,
) {
	retVal := masherytypes.ServiceEndpointMethodFilter{
		ResponseFilter: masherytypes.ResponseFilter{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "filterId",
				Name:      "filterName",
				Created:   nil,
				Updated:   nil,
				Retrieved: time.Time{},
			},
			JsonFilterFields: "abc",
			XmlFilterFields:  "def",
			Notes:            "egh",
		},
		ServiceEndpointMethod: ident,
	}

	h.mockClientWill().
		On("CreateEndpointMethodFilter", mock.Anything, ident, mock.Anything).
		Return(retVal, nil).
		Once()
}
