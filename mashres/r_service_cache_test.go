package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestCreatingServiceCacheWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceCacheResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashSvcCacheTtl, "1h")

	givenCreatingServiceCacheSucceeds(h, serviceIdent)
	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceIdentifier) {
		assert.Equal(t, "abc", id.ServiceId)
	})
}

func givenCreatingServiceCacheSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache], serviceIdent masherytypes.ServiceIdentifier) {
	rv := masherytypes.ServiceCache{
		CacheTtl: float64(45),
	}
	h.mockClientWill().
		On("CreateServiceCache", mock.Anything, serviceIdent, mock.Anything).
		Return(rv, nil).
		Once()
}
