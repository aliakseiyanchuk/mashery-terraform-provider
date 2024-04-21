package mashres

import (
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestImportServiceEndpointWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceEndpointResource)
	assert.NotNil(t, h.template.Importer)

	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	endpointId := masherytypes.ServiceEndpointIdentifier{
		ServiceIdentifier: serviceId,
		EndpointId:        "hijklmnop",
	}

	h.givenParentIdentity(t, serviceId)
	givenReadServiceEndpointSucceeds(h, endpointId)

	h.thenExecutingImport(t, "/services/abcdefg/endpoints/hijklmnop", func(data []*schema.ResourceData) error {
		if len(data) != 1 {
			return errors.New("expected 1 result")
		}

		// Verify the identity of the object created
		id, err := h.template.Mapper.Identity(data[0])
		assert.Nil(t, err)
		assert.Equal(t, "abcdefg", id.ServiceId)
		assert.Equal(t, "hijklmnop", id.EndpointId)

		return nil
	})
	// It just works.
}

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
		Return(retVal, nil).
		Once()
}

func givenReadServiceEndpointSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint], ident masherytypes.ServiceEndpointIdentifier) {
	mashTime := masherytypes.MasheryJSONTime(time.Now())

	returnedServiceEndpoint := masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:        ident.ServiceId,
			Name:      "service-name",
			Created:   &mashTime,
			Updated:   &mashTime,
			Retrieved: time.Time{},
		},
		Cache: nil,
	}
	h.mockClientWill().
		On("GetEndpoint", mock.Anything, ident).
		Return(returnedServiceEndpoint, true, nil).
		Once()
}
