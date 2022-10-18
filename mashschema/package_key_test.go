package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"testing"
)

func TestV3PackageKeyToResourceDataFromMinimalInput(t *testing.T) {
	//d := TestResourceData(&mashschema.PackageKeySchema)
	//
	//expPackId := "packId"
	//expPlanId := "planId"
	////expAppId := "appId"
	//
	//mashery.assertOk(t, d.Set(mashschema.MashPlanId, "packId::planId"))
	//mashery.assertOk(t, d.Set(mashschema.MashAppId, "memberId::userName::appId"))
	//
	//upsert := mashschema.MashPackageKeyUpsertable(d)
	//mashery.assertSameString(t, "package.id", &expPackId, &upsert.Package.Id)
	//mashery.assertSameString(t, "plan.id", &expPlanId, &upsert.Plan.Id)
}

func TestV3PackageKeyToResourceData(t *testing.T) {
	//res := schema.Resource{
	//	schema: mashschema.PackageKeySchema,
	//}
	//var tm masherytypes.MasheryJSONTime = masherytypes.MasheryJSONTime(time.Now())
	//
	//apiKey := "apiKey"
	//secret := "secret"
	//
	//var qps int64 = 10
	//var rate int64 = 20
	//
	//orig := masherytypes.MasheryPackageKey{
	//	AddressableV3Object: masherytypes.AddressableV3Object{
	//		Id:      "packId",
	//		Created: &tm,
	//		Updated: &tm,
	//	},
	//	Apikey:           &apiKey,
	//	Secret:           &secret,
	//	RateLimitCeiling: &rate,
	//	RateLimitExempt:  false,
	//	QpsLimitCeiling:  &qps,
	//	QpsLimitExempt:   true,
	//	Status:           "active",
	//	Limits: &[]masherytypes.Limit{
	//		{
	//			Period:  "second",
	//			Source:  "plan",
	//			Ceiling: 10,
	//		},
	//		{
	//			Period:  "day",
	//			Source:  "plan",
	//			Ceiling: 20,
	//		},
	//	},
	//}
	//
	//d := res.TestResourceData()
	//chk := mashschema.V3PackageKeyToResourceData(&orig, d)
	//if len(chk) != 0 {
	//	t.Errorf("setting data encountered %d problems", len(chk))
	//	for _, v := range chk {
	//		t.Errorf("%s: %s. %s", v.AttributePath, v.Summary, v.Detail)
	//	}
	//}
	//d.SetId("packId")
	//_ = d.Set(mashschema.MashPlanId, "packageId::planId")
	//
	////reverse := mashschema.MashPackageKeyUpsertable(d)
	//
	////mashery.assertSameString(t, "id", &orig.Id, &reverse.Id)
	////mashery.assertSameString(t, "apikey", orig.Apikey, reverse.Apikey)
	////mashery.assertSameString(t, "secret", orig.Secret, reverse.Secret)
	////mashery.assertSameInt64(t, "RateLimitCeiling", orig.RateLimitCeiling, reverse.RateLimitCeiling)
	////mashery.assertSameBool(t, "RateLimitExempt", &orig.RateLimitExempt, &reverse.RateLimitExempt)
	////
	////mashery.assertSameInt64(t, "QpsLimitCeiling", orig.QpsLimitCeiling, reverse.QpsLimitCeiling)
	////mashery.assertSameBool(t, "QpsLimitExempt", &orig.QpsLimitExempt, &reverse.QpsLimitExempt)
	////
	////mashery.assertSameString(t, "Status", &orig.Status, &reverse.Status)
	//
	//lims := d.Get(mashschema.MashPackageKeyLimits).([]interface{})
	//assertSameLimit(t, (*orig.Limits)[0], lims[0].(map[string]interface{}))
	//assertSameLimit(t, (*orig.Limits)[1], lims[1].(map[string]interface{}))
}

func assertSameLimit(t *testing.T, limit masherytypes.Limit, tf map[string]interface{}) {
	//mashPeriod := tf[mashschema.MashPackageKeyLimitPeriod].(string)
	//mashSource := tf[mashschema.MashPackageKeyLimitSource].(string)
	//mashCeiling := int64(tf[mashschema.MashPackageKeyLimitCeiling].(int))

	//mashery.assertSameString(t, "period", &limit.Period, &mashPeriod)
	//mashery.assertSameString(t, "source", &limit.Source, &mashSource)
	//mashery.assertSameInt64(t, "ceiling", &limit.Ceiling, &mashCeiling)
}
