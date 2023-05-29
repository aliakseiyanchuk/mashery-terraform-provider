package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestCreatingPackagePlanWillSucceed(t *testing.T) {
	h := CreateTestResource(PackagePlanResource)

	expectedIdent := masherytypes.PackagePlanIdentifier{
		PackageIdentifier: masherytypes.PackageIdentifier{
			PackageId: "packageId",
		},
		PlanId: "planId",
	}

	h.givenParentIdentity(t, expectedIdent.PackageIdentifier)
	h.givenStateFieldSetTo(t, mashschema.MashObjName, "package plan")
	givenCreatingPackagePlanSucceeds(h, expectedIdent)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackagePlanIdentifier) {
		assert.True(t, reflect.DeepEqual(expectedIdent, id))
	})
}

func givenCreatingPackagePlanSucceeds(h *ResourceTemplateMockHelper[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan],
	ident masherytypes.PackagePlanIdentifier,
) {

	retVal := masherytypes.Plan{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   ident.PlanId,
			Name: "packagePlan",
		},
		ParentPackageId: ident.PackageIdentifier,
	}

	h.mockClientWill().
		On("CreatePlan", mock.Anything, ident.PackageIdentifier, mock.Anything).
		Return(&retVal, nil).
		Once()
}
