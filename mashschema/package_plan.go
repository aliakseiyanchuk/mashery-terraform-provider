package mashschema

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashPlanId                                = "plan_id"
	MashPlanCreated                           = "created"
	MashPlanUpdated                           = "updated"
	MashPlanName                              = "name"
	MashPlanDescription                       = "description"
	MashPlanEAV                               = "extended_attribute_values"
	MashPlanSelfServiceKeyProvisioningEnabled = "self_service_provisioning"
	MashPlanAdminKeyProvisioningEnabled       = "admin_provisioning"
	MashPlanNotes                             = "notes"
	MashPlanMaxNumKeysAllowed                 = "max_keys"
	MashPlanNumKeysBeforeReview               = "keys_before_review"

	MashPlanQpsLimitCeiling            = "qps"
	MashPlanQpsLimitExempt             = "unlimited_qps"
	MashPlanQpsLimitKeyOverrideAllowed = "qps_override"

	MashPlanRateLimitCeiling            = "quota"
	MashPlanRateLimitExempt             = "unlimited_quota"
	MashPlanRateLimitKeyOverrideAllowed = "quota_override"
	MashPlanRateLimitPeriod             = "quota_period"

	MashPlanResponseFilterOverrideAllowed = "response_filter_override"
	MashPlanStatus                        = "status"
	MashPlanEmailTemplateSetId            = "email_template_set"
)

type PlanIdentifier struct {
	PackageIdentifier
	PlanId string
}

func (pi *PlanIdentifier) Self() interface{} {
	return pi
}

var PlanMapper *PlanMapperImpl

type PlanMapperImpl struct {
	MapperImpl
}

var rateLimitPeriodEnum = []string{MashDurationMinute, MashDurationHourly, MashDurationDay, MashDurationMonth}

var PlanSchema = map[string]*schema.Schema{
	MashPackagekId: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Package Id, to which this plan belongs",
	},
	MashPlanId: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Plan id",
	},
	MashPlanCreated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date/time the object was created",
	},
	MashPlanUpdated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date/time the object was updated",
	},
	MashPlanName: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Plan name",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			rv := diag.Diagnostics{}
			str := i.(string)
			if len(str) == 0 {
				rv = append(rv, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid argument",
					Detail:        "Plan name cannot be empty",
					AttributePath: path,
				})
			}
			return rv
		},
	},
	MashPlanDescription: {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "Managed by Terraform",
	},
	MashPlanEAV: {
		Type:        schema.TypeMap,
		Optional:    true,
		Computed:    true,
		Description: "Extended attribute values",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashPlanSelfServiceKeyProvisioningEnabled: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: false,
	},
	MashPlanAdminKeyProvisioningEnabled: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	MashPlanNotes: {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	MashPlanMaxNumKeysAllowed: {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  2,
	},
	MashPlanNumKeysBeforeReview: {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	},
	MashPlanQpsLimitCeiling: {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  2,
	},
	MashPlanQpsLimitExempt: {
		Type:          schema.TypeBool,
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{MashPlanQpsLimitCeiling},
	},
	MashPlanQpsLimitKeyOverrideAllowed: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	MashPlanRateLimitCeiling: {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  5000,
	},
	MashPlanRateLimitExempt: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	MashPlanRateLimitKeyOverrideAllowed: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	MashPlanRateLimitPeriod: {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return validateStringValueInSet(i, path, &rateLimitPeriodEnum)
		},
	},
	MashPlanResponseFilterOverrideAllowed: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	MashPlanStatus: {
		Type:     schema.TypeString,
		Computed: true,
	},
	MashPlanEmailTemplateSetId: {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func (pmi *PlanMapperImpl) SetIdentifier(rv *masherytypes.MasheryPlan, d *schema.ResourceData) {
	ident := PlanIdentifier{
		PackageIdentifier: PackageIdentifier{
			PackageId: pmi.GetExplicitPackageIdentifier(d),
		},
		PlanId: rv.Id,
	}
	d.SetId(CompoundId(&ident))
}

func (pmi *PlanMapperImpl) GetIdentifier(d *schema.ResourceData) *PlanIdentifier {
	rv := &PlanIdentifier{}
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (pmi *PlanMapperImpl) GetPackageIdentifier(ident *PlanIdentifier, d *schema.ResourceData) string {
	if len(ident.PackageId) > 0 {
		return ident.PackageId
	} else {
		return pmi.GetExplicitPackageIdentifier(d)
	}
}

func (pmi *PlanMapperImpl) GetExplicitPackageIdentifier(d *schema.ResourceData) string {
	return ExtractString(d, MashPackagekId, "")
}

func (pmi *PlanMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.MasheryPlan, diag.Diagnostics) {
	plnIdent := pmi.GetIdentifier(d)

	rv := masherytypes.MasheryPlan{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   plnIdent.PlanId,
			Name: ExtractString(d, MashPlanName, "Default"),
		},
		Description:                       ExtractString(d, MashPlanDescription, ""),
		Eav:                               extractEAVPointer(d, MashPlanEAV),
		SelfServiceKeyProvisioningEnabled: extractBool(d, MashPlanSelfServiceKeyProvisioningEnabled, false),
		AdminKeyProvisioningEnabled:       extractBool(d, MashPlanAdminKeyProvisioningEnabled, false),
		Notes:                             ExtractString(d, MashPlanNotes, ""),
		MaxNumKeysAllowed:                 extractInt(d, MashPlanMaxNumKeysAllowed, 2),
		NumKeysBeforeReview:               extractInt(d, MashPlanNumKeysBeforeReview, 1),
		QpsLimitCeiling:                   extractInt64Pointer(d, MashPlanQpsLimitCeiling, 0),
		QpsLimitExempt:                    extractBool(d, MashPlanQpsLimitExempt, false),
		QpsLimitKeyOverrideAllowed:        extractBool(d, MashPlanQpsLimitKeyOverrideAllowed, false),
		RateLimitCeiling:                  extractInt64Pointer(d, MashPlanRateLimitCeiling, 0),
		RateLimitExempt:                   extractBool(d, MashPlanRateLimitExempt, false),
		RateLimitKeyOverrideAllowed:       extractBool(d, MashPlanRateLimitKeyOverrideAllowed, false),
		RateLimitPeriod:                   ExtractString(d, MashPlanRateLimitPeriod, ""),
		ResponseFilterOverrideAllowed:     extractBool(d, MashPlanResponseFilterOverrideAllowed, false),
		EmailTemplateSetId:                ExtractString(d, MashPlanEmailTemplateSetId, ""),

		ParentPackageId: pmi.GetPackageIdentifier(plnIdent, d),
	}

	return rv, nil
}

func (pmi *PlanMapperImpl) PersistTyped(ctx context.Context, pln *masherytypes.MasheryPlan, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashPlanId:          pln.Id,
		MashPlanCreated:     pln.Created.ToString(),
		MashPlanUpdated:     pln.Updated.ToString(),
		MashPlanName:        pln.Name,
		MashPlanDescription: pln.Description,
		MashPlanEAV:         pln.Eav,

		MashPlanSelfServiceKeyProvisioningEnabled: pln.SelfServiceKeyProvisioningEnabled,
		MashPlanAdminKeyProvisioningEnabled:       pln.AdminKeyProvisioningEnabled,
		MashPlanNotes:                             pln.Notes,

		MashPlanMaxNumKeysAllowed:   pln.MaxNumKeysAllowed,
		MashPlanNumKeysBeforeReview: pln.NumKeysBeforeReview,

		MashPlanQpsLimitCeiling:            pln.QpsLimitCeiling,
		MashPlanQpsLimitExempt:             pln.QpsLimitExempt,
		MashPlanQpsLimitKeyOverrideAllowed: pln.QpsLimitKeyOverrideAllowed,

		MashPlanRateLimitCeiling:            pln.RateLimitCeiling,
		MashPlanRateLimitExempt:             pln.RateLimitExempt,
		MashPlanRateLimitKeyOverrideAllowed: pln.RateLimitKeyOverrideAllowed,
		MashPlanRateLimitPeriod:             pln.RateLimitPeriod,

		MashPlanResponseFilterOverrideAllowed: pln.ResponseFilterOverrideAllowed,
		MashPlanStatus:                        pln.Status,
		MashPlanEmailTemplateSetId:            nullForEmptyString(pln.EmailTemplateSetId),
	}

	return pmi.SetResourceFields(ctx, data, d)
}

func init() {
	PlanMapper = &PlanMapperImpl{
		MapperImpl{
			schema: PlanSchema,
		},
	}

	PlanMapper.identifier = func() interface{} {
		return &PlanIdentifier{}
	}
	PlanMapper.upsertFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return PlanMapper.UpsertableTyped(d)
	}
	PlanMapper.persistFunc = func(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics {
		return PlanMapper.PersistTyped(ctx, rv.(*masherytypes.MasheryPlan), d)
	}
}
