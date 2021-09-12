package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestV3PackageKeyToResourceDataFromMinimalInput(t *testing.T) {
	d := NewResourceData(&mashery.PackageKeySchema)

	expPackId := "packId"
	expPlanId := "planId"
	//expAppId := "appId"

	assertOk(t, d.Set(mashery.MashPlanId, "packId::planId"))
	assertOk(t, d.Set(mashery.MashAppId, "memberId::userName::appId"))

	upsert := mashery.MashPackageKeyUpsertable(d)
	assertSameString(t, "package.id", &expPackId, &upsert.Package.Id)
	assertSameString(t, "plan.id", &expPlanId, &upsert.Plan.Id)
}

func TestV3PackageKeyToResourceData(t *testing.T) {
	res := schema.Resource{
		Schema: mashery.PackageKeySchema,
	}
	var tm v3client.MasheryJSONTime = v3client.MasheryJSONTime(time.Now())

	apiKey := "apiKey"
	secret := "secret"

	var qps int64 = 10
	var rate int64 = 20

	orig := v3client.MasheryPackageKey{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "packId",
			Created: &tm,
			Updated: &tm,
		},
		Apikey:           &apiKey,
		Secret:           &secret,
		RateLimitCeiling: &rate,
		RateLimitExempt:  false,
		QpsLimitCeiling:  &qps,
		QpsLimitExempt:   true,
		Status:           "active",
		Limits: &[]v3client.Limit{
			{
				Period:  "second",
				Source:  "plan",
				Ceiling: 10,
			},
			{
				Period:  "day",
				Source:  "plan",
				Ceiling: 20,
			},
		},
	}

	d := res.TestResourceData()
	chk := mashery.V3PackageKeyToResourceData(&orig, d)
	if len(chk) != 0 {
		t.Errorf("setting data encountered %d problems", len(chk))
		for _, v := range chk {
			t.Errorf("%s: %s. %s", v.AttributePath, v.Summary, v.Detail)
		}
	}
	d.SetId("packId")
	_ = d.Set(mashery.MashPlanId, "packageId::planId")

	reverse := mashery.MashPackageKeyUpsertable(d)

	assertSameString(t, "id", &orig.Id, &reverse.Id)
	assertSameString(t, "apikey", orig.Apikey, reverse.Apikey)
	assertSameString(t, "secret", orig.Secret, reverse.Secret)
	assertSameInt64(t, "RateLimitCeiling", orig.RateLimitCeiling, reverse.RateLimitCeiling)
	assertSameBool(t, "RateLimitExempt", &orig.RateLimitExempt, &reverse.RateLimitExempt)

	assertSameInt64(t, "QpsLimitCeiling", orig.QpsLimitCeiling, reverse.QpsLimitCeiling)
	assertSameBool(t, "QpsLimitExempt", &orig.QpsLimitExempt, &reverse.QpsLimitExempt)

	assertSameString(t, "Status", &orig.Status, &reverse.Status)

	lims := d.Get(mashery.MashPackageKeyLimits).([]interface{})
	assertSameLimit(t, (*orig.Limits)[0], lims[0].(map[string]interface{}))
	assertSameLimit(t, (*orig.Limits)[1], lims[1].(map[string]interface{}))
}

func assertSameLimit(t *testing.T, limit v3client.Limit, tf map[string]interface{}) {
	mashPeriod := tf[mashery.MashPackageKeyLimitPeriod].(string)
	mashSource := tf[mashery.MashPackageKeyLimitSource].(string)
	mashCeiling := int64(tf[mashery.MashPackageKeyLimitCeiling].(int))

	assertSameString(t, "period", &limit.Period, &mashPeriod)
	assertSameString(t, "source", &limit.Source, &mashSource)
	assertSameInt64(t, "ceiling", &limit.Ceiling, &mashCeiling)
}
