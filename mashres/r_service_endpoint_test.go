package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
)

// TestCreatingEndpointWillSucceed verifies the creation of the composite identifier for this endpoint
func TestCreatingEndpointWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceEndpointResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashObjName, "sample-endpoint")
	givenCreatingEndpointSucceeds(h, serviceIdent)

	h.thenExecutingCreate(t)
	// It succeeds

	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceEndpointIdentifier) {
		assert.Equal(t, "abc", id.ServiceId)
		assert.Equal(t, "endpointId", id.EndpointId)
	})

}

// ---------------------------------------------------------------------------------------------------------
// GIVEN mocks

func givenCreatingEndpointSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint],
	ident masherytypes.ServiceIdentifier,
) {

	retVal := masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "endpointId",
			Name: "endpointName",
		},
		ParentServiceId: ident,
	}

	h.mockClientWill().
		On("CreateEndpoint", mock.Anything, ident, mock.Anything).
		Return(&retVal, nil).
		Once()
}
