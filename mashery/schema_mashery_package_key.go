package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashPackageKeyIdent            = "api_key"
	MashPackageKeySecret           = "secret"
	MashPackageKeyCreated          = "created"
	MashPackageKeyUpdated          = "updated"
	MashPackageKeyRateLimitCeiling = "quota"
	MashPackageKeyRateLimitExempt  = "quota_exempt"
	MashPackageKeyQpsLimitCeiling  = "qps"
	MashPackageKeyQpsLimitExempt   = "qps_exempt"
	MashPackageKeyStatus           = "status"
	MashPackageKeyLimits           = "limits"

	MashPackageKeyLimitPeriod  = "period"
	MashPackageKeyLimitSource  = "source"
	MashPackageKeyLimitCeiling = "ceiling"
)

var mashPackageKeyStatusEnum = []string{"active", "waiting", "disabled"}

var PackageKeySchema = map[string]*schema.Schema{
	MashPlanId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Plan to which this application needs to be attached",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateCompoundIdent(i, path, 2)
		},
	},
	MashAppId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Application to which this key is attached",
	},
	MashPackageKeyStatus: {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return validateStringValueInSet(i, path, &mashPackageKeyStatusEnum)
		},
	},
	MashPackageKeyLimits: {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				MashPackageKeyLimitPeriod: {
					Computed: true,
					Type:     schema.TypeString,
				},
				MashPackageKeyLimitSource: {
					Computed: true,
					Type:     schema.TypeString,
				},
				MashPackageKeyLimitCeiling: {
					Computed: true,
					Type:     schema.TypeInt,
				},
			},
		},
	},
}

func MashPackageKeyUpsertable(d *schema.ResourceData) masherytypes.MasheryPackageKey {
	plnIdent := PlanIdentifier{}
	plnIdent.From(extractString(d, MashPlanId, ""))

	return masherytypes.MasheryPackageKey{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: d.Id(),
		},
		Apikey: ExtractStringPointer(d, MashPackageKeyIdent),
		Secret: ExtractStringPointer(d, MashPackageKeySecret),

		RateLimitCeiling: extractInt64Pointer(d, MashPackageKeyRateLimitCeiling, 0),
		RateLimitExempt:  extractBool(d, MashPackageKeyRateLimitExempt, false),

		QpsLimitCeiling: extractInt64Pointer(d, MashPackageKeyQpsLimitCeiling, 0),
		QpsLimitExempt:  extractBool(d, MashPackageKeyQpsLimitExempt, false),

		Status: extractString(d, MashPackageKeyStatus, "waiting"),

		Package: &masherytypes.MasheryPackage{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: plnIdent.PackageId,
			},
		},
		Plan: &masherytypes.MasheryPlan{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: plnIdent.PlanId,
			},
		},
	}
}

func v3LimitToTerraform(inp *masherytypes.MasheryPackageKey) interface{} {
	if inp.Limits != nil {
		rv := make([]interface{}, len(*inp.Limits))
		for idx, v := range *inp.Limits {
			obj := map[string]interface{}{
				MashPackageKeyLimitSource:  v.Source,
				MashPackageKeyLimitPeriod:  v.Period,
				MashPackageKeyLimitCeiling: v.Ceiling,
			}
			rv[idx] = obj
		}
		return rv
	} else {
		return nil
	}
}

func V3PackageKeyToResourceData(inp *masherytypes.MasheryPackageKey, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashPackageKeyIdent:            inp.Apikey,
		MashPackageKeySecret:           inp.Secret,
		MashPackageKeyRateLimitCeiling: inp.RateLimitCeiling,
		MashPackageKeyRateLimitExempt:  inp.RateLimitExempt,
		MashPackageKeyQpsLimitCeiling:  inp.QpsLimitCeiling,
		MashPackageKeyQpsLimitExempt:   inp.QpsLimitExempt,
		MashPackageKeyStatus:           inp.Status,
		MashPackageKeyLimits:           v3LimitToTerraform(inp),
	}

	return SetResourceFields(data, d)
}

func initPackageKeyBoilerplate() {
	addComputedOptionalString(&PackageKeySchema, MashPackageKeyIdent, "Package key value")
	addComputedOptionalString(&PackageKeySchema, MashPackageKeySecret, "Package key secret part")

	addComputedString(&PackageKeySchema, MashPackageKeyCreated, "Date/time the object was created")
	addComputedString(&PackageKeySchema, MashPackageKeyUpdated, "Date/time the object was last updated")

	addRequiredInt(&PackageKeySchema, MashPackageKeyRateLimitCeiling, "Quota limit")
	addOptionalBoolean(&PackageKeySchema, MashPackageKeyRateLimitExempt, "If set to true, package key can make unlimited number of calls")

	addRequiredInt(&PackageKeySchema, MashPackageKeyQpsLimitCeiling, "Qps Limit")
	addOptionalBoolean(&PackageKeySchema, MashPackageKeyQpsLimitExempt, "If set to true, package key can make unlimited number of calls per second")
}

func init() {
	initPackageKeyBoilerplate()
}
