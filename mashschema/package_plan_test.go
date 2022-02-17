package mashschema_test

import (
	"testing"
)

func TestPlanIdentifier_From(t *testing.T) {
	//start := mashschema.PlanIdentifier{
	//	PackageId: "pack",
	//	PlanId:    "plan",
	//}
	//
	//check := mashschema.PlanIdentifier{}
	//check.From(start.Id())
	//
	//if check.PackageId != start.PackageId {
	//	t.Errorf("Package Ids differ")
	//}
	//if check.PlanId != start.PlanId {
	//	t.Errorf("Plan ids differ")
	//}
}

func TestV3PlanUpsertableFromMinimalInputs(t *testing.T) {
	//d := NewResourceData(&mashschema.PlanSchema)
	//
	//packId := "packId"
	//planName := "planName"

	//mashery.assertOk(t, d.Set(mashschema.MashPackagekId, packId))
	//mashery.assertOk(t, d.Set(mashschema.MashPlanName, planName))

	//upsert := mashschema.V3PlanUpsertable(d)
	//mashery.assertSameString(t, "ParentPackageId", &upsert.ParentPackageId, &packId)
	//mashery.assertSameString(t, "Name", &upsert.Name, &planName)
}

func TestV3PlanUpsertable_OnNilInputs(t *testing.T) {
	//d := NewResourceData(&mashschema.PlanSchema)
	//
	//diags := mashschema.V3PlanToResourceData(&masherytypes.MasheryPlan{}, d)
	//LogErrorDiagnostics(t, "V3 with minimal data", &diags)

	//mashery.assertResourceDoesNotHaveKey(t, d, mashschema.MashPlanEmailTemplateSetId)
}

func TestV3PlanUpsertable(t *testing.T) {
	//tm := masherytypes.MasheryJSONTime(time.Now())
	//
	//var qpsCeiling int64 = 4
	//var rateCeiling int64 = 5

	//orig := masherytypes.MasheryPlan{
	//	AddressableV3Object: masherytypes.AddressableV3Object{
	//		Id:      "planId",
	//		Name:    "planName",
	//		Created: &tm,
	//		Updated: &tm,
	//	},
	//	Description:                       "planDesc",
	//	Eav:                               &masherytypes.EAV{"a": "b"},
	//	SelfServiceKeyProvisioningEnabled: false,
	//	AdminKeyProvisioningEnabled:       false,
	//	Notes:                             "notes",
	//	MaxNumKeysAllowed:                 2,
	//	NumKeysBeforeReview:               3,
	//	QpsLimitCeiling:                   &qpsCeiling,
	//	QpsLimitExempt:                    false,
	//	QpsLimitKeyOverrideAllowed:        true,
	//	RateLimitCeiling:                  &rateCeiling,
	//	RateLimitExempt:                   false,
	//	RateLimitKeyOverrideAllowed:       true,
	//	RateLimitPeriod:                   "day",
	//	ResponseFilterOverrideAllowed:     true,
	//	Status:                            "active",
	//	EmailTemplateSetId:                "emailSetId",
	//}

	//reverse := ExchangeViaResourceData(t, &mashschema.PlanSchema, "package::planId",
	//	func(d *schema.ResourceData) diag.Diagnostics {
	//		return mashschema.V3PlanToResourceData(&orig, d)
	//	},
	//	func(d *schema.ResourceData) interface{} {
	//		return mashschema.V3PlanUpsertable(d)
	//	}).(masherytypes.MasheryPlan)

	//mashery.assertSameString(t, "Id", &orig.Id, &reverse.Id)
	//mashery.assertSameString(t, "Name", &orig.Name, &reverse.Name)
	//mashery.assertSameString(t, "Description", &orig.Description, &reverse.Description)
	//
	//mashery.assertDeepEqual(t, "Eav", orig.Eav, reverse.Eav)
	//
	//mashery.assertSameBool(t, "SelfServiceKeyProvisioningEnabled", &orig.SelfServiceKeyProvisioningEnabled, &reverse.SelfServiceKeyProvisioningEnabled)
	//mashery.assertSameBool(t, "AdminKeyProvisioningEnabled", &orig.AdminKeyProvisioningEnabled, &reverse.AdminKeyProvisioningEnabled)
	//
	//mashery.assertSameString(t, "Notes", &orig.Notes, &reverse.Notes)
	//
	//mashery.assertSameInt(t, "MaxNumKeysAllowed", &orig.MaxNumKeysAllowed, &reverse.MaxNumKeysAllowed)
	//mashery.assertSameInt(t, "NumKeysBeforeReview", &orig.NumKeysBeforeReview, &reverse.NumKeysBeforeReview)
	//
	//mashery.assertSameInt64(t, "QpsLimitCeiling", orig.QpsLimitCeiling, reverse.QpsLimitCeiling)
	//mashery.assertSameBool(t, "QpsLimitExempt", &orig.QpsLimitExempt, &reverse.QpsLimitExempt)
	//mashery.assertSameBool(t, "QpsLimitKeyOverrideAllowed", &orig.QpsLimitKeyOverrideAllowed, &reverse.QpsLimitKeyOverrideAllowed)
	//
	//mashery.assertSameInt64(t, "RateLimitCeiling", orig.RateLimitCeiling, reverse.RateLimitCeiling)
	//mashery.assertSameBool(t, "QpsLimitExempt", &orig.RateLimitExempt, &reverse.RateLimitExempt)
	//mashery.assertSameBool(t, "RateLimitKeyOverrideAllowed", &orig.RateLimitKeyOverrideAllowed, &reverse.RateLimitKeyOverrideAllowed)
	//
	//mashery.assertSameString(t, "RateLimitPeriod", &orig.RateLimitPeriod, &reverse.RateLimitPeriod)
	//mashery.assertSameBool(t, "ResponseFilterOverrideAllowed", &orig.ResponseFilterOverrideAllowed, &reverse.ResponseFilterOverrideAllowed)
	//
	//mashery.assertSameString(t, "EmailTemplateSetId", &orig.EmailTemplateSetId, &reverse.EmailTemplateSetId)
}
