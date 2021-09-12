package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestPlanIdentifier_From(t *testing.T) {
	start := mashery.PlanIdentifier{
		PackageId: "pack",
		PlanId:    "plan",
	}

	check := mashery.PlanIdentifier{}
	check.From(start.Id())

	if check.PackageId != start.PackageId {
		t.Errorf("Package Ids differ")
	}
	if check.PlanId != start.PlanId {
		t.Errorf("Plan ids differ")
	}
}

func TestV3PlanUpsertableFromMinimalInputs(t *testing.T) {
	d := NewResourceData(&mashery.PlanSchema)

	packId := "packId"
	planName := "planName"

	assertOk(t, d.Set(mashery.MashPackagekId, packId))
	assertOk(t, d.Set(mashery.MashPlanName, planName))

	upsert := mashery.V3PlanUpsertable(d)
	assertSameString(t, "ParentPackageId", &upsert.ParentPackageId, &packId)
	assertSameString(t, "Name", &upsert.Name, &planName)
}

func TestV3PlanUpsertable_OnNilInputs(t *testing.T) {
	d := NewResourceData(&mashery.PlanSchema)

	diags := mashery.V3PlanToResourceData(&v3client.MasheryPlan{}, d)
	LogErrorDiagnostics(t, "V3 with minimal data", &diags)

	assertResourceDoesNotHaveKey(t, d, mashery.MashPlanEmailTemplateSetId)
}

func TestV3PlanUpsertable(t *testing.T) {
	tm := v3client.MasheryJSONTime(time.Now())

	var qpsCeiling int64 = 4
	var rateCeiling int64 = 5

	orig := v3client.MasheryPlan{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "planId",
			Name:    "planName",
			Created: &tm,
			Updated: &tm,
		},
		Description:                       "planDesc",
		Eav:                               &v3client.EAV{"a": "b"},
		SelfServiceKeyProvisioningEnabled: false,
		AdminKeyProvisioningEnabled:       false,
		Notes:                             "notes",
		MaxNumKeysAllowed:                 2,
		NumKeysBeforeReview:               3,
		QpsLimitCeiling:                   &qpsCeiling,
		QpsLimitExempt:                    false,
		QpsLimitKeyOverrideAllowed:        true,
		RateLimitCeiling:                  &rateCeiling,
		RateLimitExempt:                   false,
		RateLimitKeyOverrideAllowed:       true,
		RateLimitPeriod:                   "day",
		ResponseFilterOverrideAllowed:     true,
		Status:                            "active",
		EmailTemplateSetId:                "emailSetId",
	}

	reverse := ExchangeViaResourceData(t, &mashery.PlanSchema, "package::planId",
		func(d *schema.ResourceData) diag.Diagnostics {
			return mashery.V3PlanToResourceData(&orig, d)
		},
		func(d *schema.ResourceData) interface{} {
			return mashery.V3PlanUpsertable(d)
		}).(v3client.MasheryPlan)

	assertSameString(t, "Id", &orig.Id, &reverse.Id)
	assertSameString(t, "Name", &orig.Name, &reverse.Name)
	assertSameString(t, "Description", &orig.Description, &reverse.Description)

	assertDeepEqual(t, "Eav", orig.Eav, reverse.Eav)

	assertSameBool(t, "SelfServiceKeyProvisioningEnabled", &orig.SelfServiceKeyProvisioningEnabled, &reverse.SelfServiceKeyProvisioningEnabled)
	assertSameBool(t, "AdminKeyProvisioningEnabled", &orig.AdminKeyProvisioningEnabled, &reverse.AdminKeyProvisioningEnabled)

	assertSameString(t, "Notes", &orig.Notes, &reverse.Notes)

	assertSameInt(t, "MaxNumKeysAllowed", &orig.MaxNumKeysAllowed, &reverse.MaxNumKeysAllowed)
	assertSameInt(t, "NumKeysBeforeReview", &orig.NumKeysBeforeReview, &reverse.NumKeysBeforeReview)

	assertSameInt64(t, "QpsLimitCeiling", orig.QpsLimitCeiling, reverse.QpsLimitCeiling)
	assertSameBool(t, "QpsLimitExempt", &orig.QpsLimitExempt, &reverse.QpsLimitExempt)
	assertSameBool(t, "QpsLimitKeyOverrideAllowed", &orig.QpsLimitKeyOverrideAllowed, &reverse.QpsLimitKeyOverrideAllowed)

	assertSameInt64(t, "RateLimitCeiling", orig.RateLimitCeiling, reverse.RateLimitCeiling)
	assertSameBool(t, "QpsLimitExempt", &orig.RateLimitExempt, &reverse.RateLimitExempt)
	assertSameBool(t, "RateLimitKeyOverrideAllowed", &orig.RateLimitKeyOverrideAllowed, &reverse.RateLimitKeyOverrideAllowed)

	assertSameString(t, "RateLimitPeriod", &orig.RateLimitPeriod, &reverse.RateLimitPeriod)
	assertSameBool(t, "ResponseFilterOverrideAllowed", &orig.ResponseFilterOverrideAllowed, &reverse.ResponseFilterOverrideAllowed)

	assertSameString(t, "EmailTemplateSetId", &orig.EmailTemplateSetId, &reverse.EmailTemplateSetId)
}
