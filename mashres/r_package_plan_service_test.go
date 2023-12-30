package mashres

import (
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/mashschemag"
	"testing"
)

func TestCreatingPackagePlanServiceSucceeds(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceResource)

	expectedIdent := masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "packageId",
			},
			PlanId: "planId",
		},
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "serviceId",
		},
	}

	h.givenParentIdentity(t, expectedIdent.PackagePlanIdentifier)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashSvcRef, expectedIdent.ServiceIdentifier)
	givenCreatingPackagePlanServiceSucceeds(h, expectedIdent)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackagePlanServiceIdentifier) {
		assert.Equal(t, "serviceId", id.ServiceId)
		assert.Equal(t, "packageId", id.PackageId)
		assert.Equal(t, "planId", id.PlanId)
	})
}

func TestCreatingPackagePlanServiceWillReturnDiagnosticsIfErrs(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceResource)

	expectedIdent := masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "packageId",
			},
			PlanId: "planId",
		},
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "serviceId",
		},
	}

	h.givenParentIdentity(t, expectedIdent.PackagePlanIdentifier)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashSvcRef, expectedIdent.ServiceIdentifier)
	givenCreatingPackagePlanServiceFails(h, expectedIdent)

	h.thenExecutingCreateWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api during creating object")
}

func TestReadingPackagePlanServiceSucceeds(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceResource)

	expectedIdent := masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "packageId",
			},
			PlanId: "planId",
		},
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "serviceId",
		},
	}

	h.givenIdentity(t, expectedIdent)
	h.givenParentIdentity(t, expectedIdent.PackagePlanIdentifier)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashSvcRef, expectedIdent.ServiceIdentifier)
	givenCheckingIfPackagePlanServiceSucceeds(h, expectedIdent, true)

	h.thenExecutingRead(t)
	// The assigned ID stays in-place.
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackagePlanServiceIdentifier) {
		assert.Equal(t, "serviceId", id.ServiceId)
		assert.Equal(t, "packageId", id.PackageId)
		assert.Equal(t, "planId", id.PlanId)
	})
}

func TestReadingPackagePlanServiceMissResetsIdentifier(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceResource)

	expectedIdent := masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "packageId",
			},
			PlanId: "planId",
		},
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "serviceId",
		},
	}

	h.givenIdentity(t, expectedIdent)
	h.givenParentIdentity(t, expectedIdent.PackagePlanIdentifier)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashSvcRef, expectedIdent.ServiceIdentifier)
	givenCheckingIfPackagePlanServiceSucceeds(h, expectedIdent, false)

	h.thenExecutingRead(t)
	// The assigned ID is removed
	h.thenAssignedIdIsEmpty(t)
}

func TestReadingPackagePlanServiceWillReturnDiagnosticsIfErrs(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceResource)

	expectedIdent := masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "packageId",
			},
			PlanId: "planId",
		},
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "serviceId",
		},
	}

	h.givenIdentity(t, expectedIdent)
	h.givenParentIdentity(t, expectedIdent.PackagePlanIdentifier)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashSvcRef, expectedIdent.ServiceIdentifier)
	givenCheckingIfPackagePlanServiceFails(h, expectedIdent)

	h.thenExecutingReadWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api")
}

func givenCreatingPackagePlanServiceSucceeds(h *ResourceTemplateMockHelper[masherytypes.PackagePlanIdentifier, masherytypes.PackagePlanServiceIdentifier, mashschemag.PackagePlanServiceParam],
	expectedIdent masherytypes.PackagePlanServiceIdentifier) {

	rv := masherytypes.AddressableV3Object{
		Id: expectedIdent.ServiceId,
	}

	h.mockClientWill().
		On("CreatePlanService", mock.Anything, expectedIdent).
		Return(rv, nil).
		Once()
}

func givenCheckingIfPackagePlanServiceSucceeds(h *ResourceTemplateMockHelper[masherytypes.PackagePlanIdentifier, masherytypes.PackagePlanServiceIdentifier, mashschemag.PackagePlanServiceParam],
	expectedIdent masherytypes.PackagePlanServiceIdentifier, expResult bool) {
	h.mockClientWill().
		On("CheckPlanServiceExists", mock.Anything, expectedIdent).
		Return(expResult, nil).
		Once()
}

func givenCheckingIfPackagePlanServiceFails(h *ResourceTemplateMockHelper[masherytypes.PackagePlanIdentifier, masherytypes.PackagePlanServiceIdentifier, mashschemag.PackagePlanServiceParam],
	expectedIdent masherytypes.PackagePlanServiceIdentifier) {
	h.mockClientWill().
		On("CheckPlanServiceExists", mock.Anything, expectedIdent).
		Return(false, errors.New("unit test error in CheckPlanServiceExists method")).
		Once()
}

func givenCreatingPackagePlanServiceFails(h *ResourceTemplateMockHelper[masherytypes.PackagePlanIdentifier, masherytypes.PackagePlanServiceIdentifier, mashschemag.PackagePlanServiceParam],
	expectedIdent masherytypes.PackagePlanServiceIdentifier) {

	h.mockClientWill().
		On("CreatePlanService", mock.Anything, expectedIdent).
		Return(masherytypes.AddressableV3Object{}, errors.New("unit test failure during creating package plan service")).
		Once()
}
