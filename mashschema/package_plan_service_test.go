package mashschema_test

import (
	"testing"
)

func TestPlanServiceIdentifier_From(t *testing.T) {
	//start := mashschema.PlanServiceIdentifier{
	//	PlanIdentifier: mashschema.PlanIdentifier{
	//		PackageId: "pack",
	//		PlanId:    "plan",
	//	},
	//	ServiceId: "service",
	//}
	//
	//check := mashschema.PlanServiceIdentifier{}
	//check.From(start.Id())
	//
	//if check.PackageId != start.PackageId {
	//	t.Errorf("Mismatching package id")
	//}
	//
	//if check.PlanId != start.PlanId {
	//	t.Errorf("Mismatching plan id")
	//}
	//
	//if check.ServiceId != start.ServiceId {
	//	t.Errorf("Mismatching service id")
	//}
}

//
//func TestV3MasheryPlanServiceUpsertable_fromId(t *testing.T) {
//	d := schema.ResourceData{}
//	d.SetId("a::b::c")
//
//	upsert := mashschema.V3MasheryPlanServiceUpsertable(&d)
//	assertMasheryPlanServiceUpsert(t, &upsert, "a", "b", "c")
//}
//
//func TestV3MasheryPlanServiceUpsertable_fromRefs(t *testing.T) {
//	res := schema.Resource{
//		schema: mashschema.PlanService,
//	}
//	d := res.TestResourceData()
//
//	_ = d.Set(mashschema.MashPlanId, "a::b")
//	_ = d.Set(mashschema.MashSvcId, "c")
//
//	upsert := mashschema.V3MasheryPlanServiceUpsertable(d)
//	assertMasheryPlanServiceUpsert(t, &upsert, "a", "b", "c")
//}
//
//func assertMasheryPlanServiceUpsert(t *testing.T, upsert *masherytypes.MasheryPlanService, packageId, planId, serviceId string) {
//	if upsert.PackageId != packageId {
//		t.Errorf("Unexpected package id")
//	}
//
//	if upsert.PlanId != planId {
//		t.Errorf("Unexpected plan id")
//	}
//
//	if upsert.ServiceId != serviceId {
//		t.Errorf("Unexpected service id")
//	}
//}
