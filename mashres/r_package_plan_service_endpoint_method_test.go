package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/mashschemag"
	"testing"
)

func TestCreatingPackagePlanServiceEndpointMethodSucceeds(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceEndpointMethodResource)

	stateParentIdent, stateMethodIdent, stateFilterIdent, callExpectedMethodIdent, callExpectedFilterIdent := givenMethodAndFilterIdentifiers()

	h.givenParentIdentity(t, stateParentIdent)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.ServiceEndpointMethodRef, stateMethodIdent)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.MashSvcEndpointMethodFilterRef, stateFilterIdent)

	givenCreatingPackagePlanServiceEndpointMethodSucceeds(h, callExpectedMethodIdent)
	givenCreatingPackagePlanServiceEndpointMethodFilterSucceeds(h, callExpectedFilterIdent)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackagePlanServiceEndpointMethodIdentifier) {
		id.PackageId = stateParentIdent.PackageId
		id.PlanId = stateParentIdent.PlanId
		id.ServiceId = stateParentIdent.ServiceId
		id.EndpointId = stateParentIdent.EndpointId
		id.MethodId = stateMethodIdent.MethodId
	})
}

func TestCreatingPackagePlanServiceEndpointMethodWithoutFilterSucceeds(t *testing.T) {
	h := CreateTestResource(PackagePlanServiceEndpointMethodResource)

	stateParentIdent, stateMethodIdent, _, callExpectedMethodIdent, _ := givenMethodAndFilterIdentifiers()

	h.givenParentIdentity(t, stateParentIdent)
	h.givenStateFieldSetToWrappedJSON(t, mashschema.ServiceEndpointMethodRef, stateMethodIdent)

	givenCreatingPackagePlanServiceEndpointMethodSucceeds(h, callExpectedMethodIdent)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackagePlanServiceEndpointMethodIdentifier) {
		id.PackageId = stateParentIdent.PackageId
		id.PlanId = stateParentIdent.PlanId
		id.ServiceId = stateParentIdent.ServiceId
		id.EndpointId = stateParentIdent.EndpointId
		id.MethodId = stateMethodIdent.MethodId
	})
}

// ---------------------------------------------------------------------------------------------------------
// Given conditions

func givenMethodAndFilterIdentifiers() (masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier) {
	stateParentIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	stateParentIdent.PackageId = "packId"
	stateParentIdent.PlanId = "planId"
	stateParentIdent.ServiceId = "serviceId"
	stateParentIdent.EndpointId = "endpointId"

	stateMethodIdent := masherytypes.ServiceEndpointMethodIdentifier{}
	stateMethodIdent.ServiceEndpointIdentifier = stateParentIdent.ServiceEndpointIdentifier
	stateMethodIdent.MethodId = "methodId"

	stateFilterIdent := masherytypes.ServiceEndpointMethodFilterIdentifier{}
	stateFilterIdent.ServiceEndpointMethodIdentifier = stateMethodIdent
	stateFilterIdent.FilterId = "filterId"

	callExpectedMethodIdent := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
	callExpectedMethodIdent.PackagePlanIdentifier = stateParentIdent.PackagePlanIdentifier
	callExpectedMethodIdent.ServiceEndpointMethodIdentifier = stateMethodIdent

	callExpectedFilterIdent := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{}
	callExpectedFilterIdent.PackagePlanIdentifier = stateParentIdent.PackagePlanIdentifier
	callExpectedFilterIdent.ServiceEndpointMethodFilterIdentifier = stateFilterIdent
	return stateParentIdent, stateMethodIdent, stateFilterIdent, callExpectedMethodIdent, callExpectedFilterIdent
}

func givenCreatingPackagePlanServiceEndpointMethodSucceeds(h *ResourceTemplateMockHelper[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, mashschemag.PackagePlanServiceEndpointMethodParam],
	expectedMethodIdent masherytypes.PackagePlanServiceEndpointMethodIdentifier) {

	rv := masherytypes.PackagePlanServiceEndpointMethod{}
	rv.Id = "methId"
	rv.Name = "method-name"

	h.mockClientWill().
		On("CreatePackagePlanMethod", mock.Anything, expectedMethodIdent).
		Return(&rv, nil).
		Once()
}

func givenCreatingPackagePlanServiceEndpointMethodFilterSucceeds(h *ResourceTemplateMockHelper[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, mashschemag.PackagePlanServiceEndpointMethodParam],
	expectedIdent masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier) {

	rv := masherytypes.PackagePlanServiceEndpointMethodFilter{}
	rv.Id = "filterId"
	rv.Name = "filterName"
	rv.PackagePlanServiceEndpointMethod = masherytypes.PackagePlanServiceEndpointMethodIdentifier{
		ServiceEndpointMethodIdentifier: expectedIdent.ServiceEndpointMethodIdentifier,
		PackagePlanIdentifier:           expectedIdent.PackagePlanIdentifier,
	}

	h.mockClientWill().
		On("CreatePackagePlanMethodFilter", mock.Anything, expectedIdent).
		Return(&rv, nil).
		Once()
}
