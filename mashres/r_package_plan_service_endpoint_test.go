package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/mashschemag"
	"testing"
)

func TestCreatingPackagePlanServiceEndpointSucceeds(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceEndpointResource)

	expectedIdent := masherytypes.PackagePlanServiceEndpointIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "packageId",
			},
			PlanId: "planId",
		},
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
			ServiceIdentifier: masherytypes.ServiceIdentifier{
				ServiceId: "serviceId",
			},
			EndpointId: "endpointId",
		},
	}

	parentIdent := masherytypes.PackagePlanServiceIdentifier{
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

	ident := masherytypes.PackagePlanServiceEndpointIdentifier{
		ServiceEndpointIdentifier: expectedIdent.ServiceEndpointIdentifier,
		PackagePlanIdentifier:     expectedIdent.PackagePlanIdentifier,
	}

	h.givenParentIdentity(t, parentIdent)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashServiceEndpointRef, expectedIdent.ServiceEndpointIdentifier)
	givenCreatingPackagePlanServiceEndpointSucceeds(h, ident)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackagePlanServiceEndpointIdentifier) {
		assert.Equal(t, "packageId", id.PackageId)
		assert.Equal(t, "planId", id.PlanId)
		assert.Equal(t, "serviceId", id.ServiceId)
		assert.Equal(t, "endpointId", id.EndpointId)
	})
}

func TestCreatingPackagePlanServiceEndpointWithConflictingParamsWillErr(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceEndpointResource)

	parentIdent := masherytypes.PackagePlanServiceIdentifier{
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

	ident := masherytypes.ServiceEndpointIdentifier{
		EndpointId: "endpointId",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "conflictingServiceId",
		},
	}

	h.givenParentIdentity(t, parentIdent)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashServiceEndpointRef, ident)

	h.thenExecutingCreateWillYieldDiagnostic(t, "object cannot be created due to conflict in the parameters")

}

func givenCreatingPackagePlanServiceEndpointSucceeds(h *ResourceTemplateMockHelper[masherytypes.PackagePlanServiceIdentifier, masherytypes.PackagePlanServiceEndpointIdentifier, mashschemag.PackagePlanServiceEndpointParam],
	expectedIdent masherytypes.PackagePlanServiceEndpointIdentifier) {

	rv := masherytypes.AddressableV3Object{
		Id: expectedIdent.ServiceId,
	}

	h.mockClientWill().
		On("CreatePlanEndpoint", mock.Anything, expectedIdent).
		Return(&rv, nil).
		Once()
}
