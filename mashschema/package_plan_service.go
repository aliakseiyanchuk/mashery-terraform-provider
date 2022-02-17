package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PlanServiceIdentifier struct {
	PlanIdentifier
	ServiceId string
}

var PlanServiceMapper *PlanServiceMapperImpl

type PlanServiceMapperImpl struct {
	MapperImpl
}

func (psi *PlanServiceIdentifier) Self() interface{} {
	return psi
}

func (psm *PlanServiceMapperImpl) CreateIdentifierTyped() *PlanServiceIdentifier {
	return &PlanServiceIdentifier{}
}

func (sem *PlanServiceMapperImpl) GetIdentifier(d *schema.ResourceData) *PlanServiceIdentifier {
	rv := &PlanServiceIdentifier{}
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (psm *PlanServiceMapperImpl) UpsertableTyped(d *schema.ResourceData) masherytypes.MasheryPlanService {
	if d.Id() == "" {
		planId := PlanIdentifier{}
		CompoundIdFrom(&planId, ExtractString(d, MashPlanId, ""))

		return masherytypes.MasheryPlanService{
			PackageId: planId.PackageId,
			PlanId:    planId.PlanId,
			ServiceId: ExtractString(d, MashSvcId, ""),
		}
	} else {
		ident := PlanServiceIdentifier{}
		CompoundIdFrom(&ident, d.Id())

		return masherytypes.MasheryPlanService{
			PackageId: ident.PackageId,
			PlanId:    ident.PlanId,
			ServiceId: ident.ServiceId,
		}
	}
}

func init() {
	PlanServiceMapper = &PlanServiceMapperImpl{
		MapperImpl{
			schema: map[string]*schema.Schema{
				MashPlanId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Plan to which the service needs to be attached",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &PlanIdentifier{}
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
		},
	}
}
