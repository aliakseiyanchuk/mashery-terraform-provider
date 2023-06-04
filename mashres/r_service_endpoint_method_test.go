package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestCreateServiceEndpointMethodWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceEndpointMethodResource)

	serviceIdent := masherytypes.ServiceEndpointIdentifier{}
	serviceIdent.ServiceId = "serviceId"
	serviceIdent.EndpointId = "endpointId"

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashObjName, "sample-method")

	givenCreatingEndpointMethodSucceeds(h, serviceIdent)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceEndpointMethodIdentifier) {
		assert.Equal(t, "serviceId", id.ServiceId)
		assert.Equal(t, "endpointId", id.EndpointId)
		assert.Equal(t, "methId", id.MethodId)
	})
}

func givenCreatingEndpointMethodSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod],
	ident masherytypes.ServiceEndpointIdentifier,
) {
	retVal := masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "methId",
				Name:      "objectName",
				Created:   nil,
				Updated:   nil,
				Retrieved: time.Time{},
			},
			SampleJsonResponse: "abc",
			SampleXmlResponse:  "def",
		},
		ParentEndpointId: ident,
	}

	h.mockClientWill().
		On("CreateEndpointMethod", mock.Anything, ident, mock.Anything).
		Return(&retVal, nil).
		Once()
}
