package mashschema

import (
	"context"
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

type PackageKeyIdentifier struct {
	ApplicationIdentifier
	KeyId string
}

func (pki *PackageKeyIdentifier) Self() interface{} {
	return pki
}

var PackageKeyMapper *PackageKeyMapperImpl

type PackageKeyMapperImpl struct {
	MapperImpl
}

func (pkm *PackageKeyMapperImpl) GetIdentifier(d *schema.ResourceData) *PackageKeyIdentifier {
	rv := &PackageKeyIdentifier{}
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (pkm *PackageKeyMapperImpl) SetIdentifier(appId *ApplicationIdentifier, id string, d *schema.ResourceData) *PackageKeyIdentifier {
	rv := &PackageKeyIdentifier{
		ApplicationIdentifier: ApplicationIdentifier{
			MemberIdentifier: MemberIdentifier{
				MemberId: appId.MemberId,
				Username: appId.Username,
			},
			AppId: appId.AppId,
		},
		KeyId: id,
	}
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (pkm *PackageKeyMapperImpl) GetApplicationIdentifier(d *schema.ResourceData) *ApplicationIdentifier {
	rv := &ApplicationIdentifier{}
	CompoundIdFrom(rv, ExtractString(d, MashAppId, ""))

	return rv
}

func (pkm *PackageKeyMapperImpl) GetPlanIdentifier(d *schema.ResourceData) *PlanIdentifier {
	rv := &PlanIdentifier{}
	CompoundIdFrom(rv, ExtractString(d, MashPlanId, ""))

	return rv
}

func (pkm *PackageKeyMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.MasheryPackageKey, diag.Diagnostics) {
	plnIdent := pkm.GetPlanIdentifier(d)

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

		Status: ExtractString(d, MashPackageKeyStatus, "waiting"),

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
	}, nil
}

func (pkm *PackageKeyMapperImpl) persistLimits(inp *masherytypes.MasheryPackageKey) interface{} {
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

func (pkm *PackageKeyMapperImpl) PersistTyped(ctx context.Context, inp *masherytypes.MasheryPackageKey, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashPackageKeyIdent:            inp.Apikey,
		MashPackageKeySecret:           inp.Secret,
		MashPackageKeyRateLimitCeiling: inp.RateLimitCeiling,
		MashPackageKeyRateLimitExempt:  inp.RateLimitExempt,
		MashPackageKeyQpsLimitCeiling:  inp.QpsLimitCeiling,
		MashPackageKeyQpsLimitExempt:   inp.QpsLimitExempt,
		MashPackageKeyStatus:           inp.Status,
		MashPackageKeyLimits:           pkm.persistLimits(inp),
	}

	return pkm.SetResourceFields(ctx, data, d)
}

func initPackageKeyBoilerplate() {
	addComputedOptionalString(&PackageKeyMapper.schema, MashPackageKeyIdent, "Package key value")
	addComputedOptionalString(&PackageKeyMapper.schema, MashPackageKeySecret, "Package key secret part")

	addComputedString(&PackageKeyMapper.schema, MashPackageKeyCreated, "Date/time the object was created")
	addComputedString(&PackageKeyMapper.schema, MashPackageKeyUpdated, "Date/time the object was last updated")

	addRequiredInt(&PackageKeyMapper.schema, MashPackageKeyRateLimitCeiling, "Quota limit")
	addOptionalBoolean(&PackageKeyMapper.schema, MashPackageKeyRateLimitExempt, "If set to true, package key can make unlimited number of calls")

	addRequiredInt(&PackageKeyMapper.schema, MashPackageKeyQpsLimitCeiling, "Qps Limit")
	addOptionalBoolean(&PackageKeyMapper.schema, MashPackageKeyQpsLimitExempt, "If set to true, package key can make unlimited number of calls per second")
}

func init() {
	PackageKeyMapper = &PackageKeyMapperImpl{
		MapperImpl{
			schema: map[string]*schema.Schema{
				MashPlanId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Plan to which this application needs to be attached",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &PlanIdentifier{}
						})
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
			},
		},
	}
	initPackageKeyBoilerplate()

	PackageKeyMapper.identifier = func() interface{} {
		return &PackageKeyIdentifier{}
	}
	PackageKeyMapper.upsertFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return PackageKeyMapper.UpsertableTyped(d)
	}
	PackageKeyMapper.persistFunc = func(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics {
		return PackageKeyMapper.PersistTyped(ctx, rv.(*masherytypes.MasheryPackageKey), d)
	}
}
