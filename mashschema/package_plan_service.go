package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var PlanServiceMapper *PlanServiceMapperImpl

type PlanServiceMapperImpl struct {
	ResourceMapperImpl
}

func (psm *PlanServiceMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.PackagePlanServiceIdentifier, V3ObjectIdentifier, diag.Diagnostics) {
	dg := diag.Diagnostics{}
	planIdent := masherytypes.PackagePlanIdentifier{}

	if !CompoundIdFrom(&planIdent, ExtractString(d, MashPlanId, "")) {
		dg = append(dg, psm.lackingIdentificationDiagnostic(MashPlanId))
	}

	rv := masherytypes.PackagePlanServiceIdentifier{
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: ExtractString(d, MashSvcId, ""),
		},
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: planIdent.PackageId,
			},
			PlanId: planIdent.PlanId,
		},
	}

	return rv, nil, dg
}

func init() {
	PlanServiceMapper = &PlanServiceMapperImpl{
		ResourceMapperImpl{
			v3ObjectName: "package plan service",
			schema: map[string]*schema.Schema{
				MashPlanId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Plan to which the service needs to be attached",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.PackagePlanIdentifier{}
						})
					},
				},
				MashSvcId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Service to expose in this plan",
				},
			},
			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv, _, dg := PlanServiceMapper.UpsertableTyped(d)
				return rv, dg
			},
			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return PlanServiceMapper.UpsertableTyped(d)
			},
			persistFunc: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				return PlanServiceMapper.persistMap(rv, nil, d)
			},
		},
	}
}
