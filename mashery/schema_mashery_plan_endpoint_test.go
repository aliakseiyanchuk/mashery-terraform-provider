package mashery_test

import (
	"terraform-provider-mashery/mashery"
	"testing"
)

func TestV3MasheryPlanEndpointUpsertable(t *testing.T) {
	d := NewResourceData(&mashery.PlanEndpointSchema)

	assertOk(t, d.Set(mashery.MashEndpointId, "serviceId::endpointId"))
	assertOk(t, d.Set(mashery.MashPlanServiceId, "packageId::planId::serviceId"))

	upsert, diag := mashery.V3MasheryPlanEndpointUpsertable(d)
	LogErrorDiagnostics(t, "plan endpoint upsert calc", &diag)

	pckId := "packageId"
	plnId := "planId"
	serviceId := "serviceId"
	endpointId := "endpointId"

	assertSameString(t, "PackageId", &upsert.PackageId, &pckId)
	assertSameString(t, "PlanId", &upsert.PlanId, &plnId)
	assertSameString(t, "ServiceId", &upsert.ServiceId, &serviceId)
	assertSameString(t, "EndpointId", &upsert.EndpointId, &endpointId)
}

// Mashery plan should include only endpoints that reside in the same service plan. It it technically possible
// to specify offending endpoints. This would be detected and an error would be returned.
func TestV3MasheryPlanEndpointUpsertableWithServiceConflict(t *testing.T) {
	d := NewResourceData(&mashery.PlanEndpointSchema)

	assertOk(t, d.Set(mashery.MashEndpointId, "offendingServiceId::endpointId"))
	assertOk(t, d.Set(mashery.MashPlanServiceId, "packageId::planId::serviceId"))

	_, diag := mashery.V3MasheryPlanEndpointUpsertable(d)
	if len(diag) == 1 {
		expSummary := "Incompatible V3 object hierarchy"
		assertSameString(t, "Summary", &diag[0].Summary, &expSummary)
	} else {
		t.Errorf("Diagnostics should have found an error")
	}
}
