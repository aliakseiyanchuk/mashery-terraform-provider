package mashschema_test

import (
	"testing"
)

func TestV3MasheryPlanEndpointUpsertable(t *testing.T) {
	//d := NewResourceData(&mashschema.PlanEndpointSchema)

	//mashery.assertOk(t, d.Set(mashschema.MashEndpointId, "serviceId::endpointId"))
	//mashery.assertOk(t, d.Set(mashschema.MashPlanServiceId, "packageId::planId::serviceId"))

	//upsert, diag := mashschema.V3MasheryPlanEndpointUpsertable(d)
	//LogErrorDiagnostics(t, "plan endpoint upsert calc", &diag)

	//pckId := "packageId"
	//plnId := "planId"
	//serviceId := "serviceId"
	//endpointId := "endpointId"

	//mashery.assertSameString(t, "PackageId", &upsert.PackageId, &pckId)
	//mashery.assertSameString(t, "PlanId", &upsert.PlanId, &plnId)
	//mashery.assertSameString(t, "ServiceId", &upsert.ServiceId, &serviceId)
	//mashery.assertSameString(t, "EndpointId", &upsert.EndpointId, &endpointId)
}

// Mashery plan should include only endpoints that reside in the same service plan. It it technically possible
// to specify offending endpoints. This would be detected and an error would be returned.
func TestV3MasheryPlanEndpointUpsertableWithServiceConflict(t *testing.T) {
	//d := NewResourceData(&mashschema.PlanEndpointSchema)

	//mashery.assertOk(t, d.Set(mashschema.MashEndpointId, "offendingServiceId::endpointId"))
	//mashery.assertOk(t, d.Set(mashschema.MashPlanServiceId, "packageId::planId::serviceId"))

	//_, diag := mashschema.V3MasheryPlanEndpointUpsertable(d)
	//if len(diag) == 1 {
	//	//expSummary := "Incompatible V3 object hierarchy"
	//	//mashery.assertSameString(t, "Summary", &diag[0].Summary, &expSummary)
	//} else {
	//	t.Errorf("Diagnostics should have found an error")
	//}
}
