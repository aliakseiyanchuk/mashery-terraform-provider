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

	// Pointless configuration no. 2
	processorCfg := map[string]interface{}{
		mashschema.MashEndpointProcessorAdapter:           "com.github.fake-adapter",
		mashschema.MashEndpointProcessorPreProcessEnabled: true,
	}

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashObjName, "sample-endpoint")
	h.givenStateFieldSetTo(t, mashschema.MashEndpointProcessor, []interface{}{processorCfg})
	givenCreatingEndpointSucceeds(h, serviceIdent)

	h.thenExecutingCreate(t)
	// It succeeds

	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceEndpointIdentifier) {
		assert.Equal(t, "abc", id.ServiceId)
		assert.Equal(t, "endpointId", id.EndpointId)
	})
}

func TestCreatingPointlessEndpointWillBeRejected(t *testing.T) {
	h := CreateTestResource(ServiceEndpointResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	// Pointless configuration no. 1
	processorCfg := map[string]interface{}{
		mashschema.MashEndpointProcessorAdapter:           "",
		mashschema.MashEndpointProcessorPreProcessEnabled: false,
	}

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashEndpointProcessor, []interface{}{processorCfg})

	h.thenExecutingCreateWillYieldDiagnostic(t, "invalid input for field processor")

	// Pointless configuration no. 2
	processorCfg = map[string]interface{}{
		mashschema.MashEndpointProcessorAdapter:           "com.github.fake-adapter",
		mashschema.MashEndpointProcessorPreProcessEnabled: false,
	}

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashEndpointProcessor, []interface{}{processorCfg})

	h.thenExecutingCreateWillYieldDiagnostic(t, "invalid input for field processor")
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
