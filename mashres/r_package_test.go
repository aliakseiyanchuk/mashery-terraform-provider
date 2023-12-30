package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"testing"
)

func TestCreatingPackageWillSucceed(t *testing.T) {
	h := CreateTestResource(PackageResource)

	expIdent := masherytypes.PackageIdentifier{
		PackageId: "pack-abc",
	}

	h.givenStateFieldSetTo(t, mashschema.MashObjName, "sample-package")
	givenCreatingPackageSucceeds(h, expIdent)
	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackageIdentifier) {
		assert.True(t, reflect.DeepEqual(expIdent, id))
	})
}

func givenCreatingPackageSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.PackageIdentifier, masherytypes.Package], ident masherytypes.PackageIdentifier) {
	rv := masherytypes.Package{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   ident.PackageId,
			Name: "Created Package name",
		},
		Description:                 "This is terraform package",
		NotifyDeveloperPeriod:       "",
		NotifyDeveloperNearQuota:    false,
		NotifyDeveloperOverQuota:    false,
		NotifyDeveloperOverThrottle: false,
		NotifyAdminPeriod:           "",
		NotifyAdminNearQuota:        false,
		NotifyAdminOverQuota:        false,
		NotifyAdminOverThrottle:     false,
		NotifyAdminEmails:           "",
		NearQuotaThreshold:          nil,
		Eav:                         nil,
		KeyAdapter:                  "",
		KeyLength:                   nil,
		SharedSecretLength:          nil,
		Plans:                       nil,
	}

	h.mockClientWill().
		On("CreatePackage", mock.Anything, mock.Anything).
		Return(rv, nil).
		Once()
}
