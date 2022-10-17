package mashschema

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

var PackageKeyMapper *PackageKeyMapperImpl

type PackageKeyMapperImpl struct {
	ResourceMapperImpl
}

func (pkm *PackageKeyMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.PackageKey, V3ObjectIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}

	appIdent := masherytypes.ApplicationIdentifier{}
	if !CompoundIdFrom(&appIdent, ExtractString(d, MashAppId, "")) {
		rvd = append(rvd, pkm.lackingIdentificationDiagnostic(MashAppId))
	}

	plnIdent := masherytypes.PackagePlanIdentifier{}
	if !CompoundIdFrom(&plnIdent, ExtractString(d, MashPlanId, "")) {
		rvd = append(rvd, pkm.lackingIdentificationDiagnostic(MashPlanId))
	}

	return masherytypes.PackageKey{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: d.Id(),
		},
		Apikey: ExtractStringPointer(d, MashPackageKeyIdent),
		Secret: ExtractStringPointer(d, MashPackageKeySecret),

		RateLimitCeiling: extractInt64Pointer(d, MashPackageKeyRateLimitCeiling, 0),
		RateLimitExempt:  ExtractBool(d, MashPackageKeyRateLimitExempt, false),

		QpsLimitCeiling: extractInt64Pointer(d, MashPackageKeyQpsLimitCeiling, 0),
		QpsLimitExempt:  ExtractBool(d, MashPackageKeyQpsLimitExempt, false),

		Status: ExtractString(d, MashPackageKeyStatus, "waiting"),

		Package: &masherytypes.Package{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: plnIdent.PackageId,
			},
		},
		Plan: &masherytypes.Plan{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id: plnIdent.PlanId,
			},
		},
	}, appIdent, nil
}

func (pkm *PackageKeyMapperImpl) persistLimits(inp masherytypes.PackageKey) interface{} {
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

func (pkm *PackageKeyMapperImpl) PersistTyped(inp masherytypes.PackageKey, d *schema.ResourceData) diag.Diagnostics {
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

	return pkm.persistMap(inp.Identifier(), data, d)
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
		ResourceMapperImpl: ResourceMapperImpl{
			v3ObjectName: "package key",
			schema: map[string]*schema.Schema{
				MashPlanId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Plan to which this application needs to be attached",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.PackagePlanIdentifier{}
						})
					},
				},
				MashAppId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Application to which this key is attached",

					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.ApplicationIdentifier{}
						})
					},
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
			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv := masherytypes.PackageKeyIdentifier{PackageKeyId: d.Id()}

				rvd := diag.Diagnostics{}
				if len(rv.PackageKeyId) == 0 {
					rvd = append(rvd, PackageKeyMapper.lackingIdentificationDiagnostic("id"))
				}

				return rv, rvd
			},
			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return PackageKeyMapper.Upsertable(d)
			},
			persistFunc: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				ptr := rv.(*masherytypes.PackageKey)
				return PackageKeyMapper.PersistTyped(*ptr, d)
			},
		},
	}
	initPackageKeyBoilerplate()
}
