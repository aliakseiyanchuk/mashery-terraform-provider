package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestCreatingServiceOAuthWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceOAuthResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashSvcOAuthAccessTokenTtl, "1h")
	h.givenStateFieldSetTo(t, mashschema.MashSvcOAuthGrantTypes, []string{"client-credentials"})

	givenCreatingServiceOAuthSucceeds(h)
	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceIdentifier) {
		assert.Equal(t, "abc", id.ServiceId)
	})
}

func givenCreatingServiceOAuthSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]) {
	rv := masherytypes.MasheryOAuth{
		SecureTokensEnabled: true,
	}
	h.mockClientWill().
		On("CreateServiceOAuthSecurityProfile", mock.Anything, mock.Anything).
		Return(&rv, nil).
		Once()
}
