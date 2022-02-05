package mashery

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

func (psi *PlanServiceIdentifier) IsIdentified() bool {
	return psi.PlanIdentifier.IsIdentified() && len(psi.ServiceId) > 0
}

func (psi *PlanServiceIdentifier) Id() string {
	return CreateCompoundId(psi.PackageId, psi.PlanId, psi.ServiceId)
}

func (psi *PlanServiceIdentifier) From(id string) {
	ParseCompoundId(id, &psi.PackageId, &psi.PlanId, &psi.ServiceId)
}

var PlanService = map[string]*schema.Schema{
	MashPlanId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Plan to which the service needs to be attached",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateCompoundIdent(i, path, 2)
		},
	},
	MashSvcId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Service to expose in this plan",
	},
}

func V3MasheryPlanServiceUpsertable(d *schema.ResourceData) masherytypes.MasheryPlanService {
	if d.Id() == "" {
		planId := PlanIdentifier{}
		planId.From(extractString(d, MashPlanId, ""))

		return masherytypes.MasheryPlanService{
			PackageId: planId.PackageId,
			PlanId:    planId.PlanId,
			ServiceId: extractString(d, MashSvcId, ""),
		}
	} else {
		ident := PlanServiceIdentifier{}
		ident.From(d.Id())

		return masherytypes.MasheryPlanService{
			PackageId: ident.PackageId,
			PlanId:    ident.PlanId,
			ServiceId: ident.ServiceId,
		}
	}
}
