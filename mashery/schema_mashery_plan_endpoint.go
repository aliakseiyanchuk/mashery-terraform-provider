package mashery

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashPlanServiceId = "plan_service_id"

	// Method-related constants
	MashPlanMethodBlock = "method"
	MashPlanMethodId    = "method_id" // TODO: Shared, to be moved

	//MashPlanMethodFilterName       = "name"
	MashPlanMethodXmlFilterFields  = "xml_filter_fields"
	MashPlanMethodJsonFilterFields = "json_filter_fields"
)

var PlanEndpointSchema = map[string]*schema.Schema{
	MashPlanServiceId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Plan service",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateCompoundIdent(i, path, 3)
		},
	},
	MashEndpointId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Endpoint to include",
	},
}

type PlanEndpointIdentifier struct {
	PlanServiceIdentifier
	EndpointId string
}

func (pei *PlanEndpointIdentifier) Id() string {
	return CreateCompoundId(pei.PackageId, pei.PlanId, pei.ServiceId, pei.EndpointId)
}

func (pei *PlanEndpointIdentifier) From(id string) {
	ParseCompoundId(id, &pei.PackageId, &pei.PlanId, &pei.ServiceId, &pei.EndpointId)
}

func (pei *PlanEndpointIdentifier) IsIdentified() bool {
	return pei.PlanServiceIdentifier.IsIdentified() && len(pei.EndpointId) > 0
}

func V3MasheryPlanEndpointUpsertable(d *schema.ResourceData) (masherytypes.MasheryPlanServiceEndpoint, diag.Diagnostics) {
	refPlanService := PlanServiceIdentifier{}
	refPlanService.From(extractString(d, MashPlanServiceId, ""))

	planEndpoint := ServiceEndpointIdentifier{}
	planEndpoint.From(extractString(d, MashEndpointId, ""))

	if refPlanService.ServiceId != planEndpoint.ServiceId {
		return masherytypes.MasheryPlanServiceEndpoint{}, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Incompatible V3 object hierarchy",
			Detail: fmt.Sprintf("Endpoint %s belonging to service %s cannot be added to plan %s service %s",
				planEndpoint.EndpointId,
				planEndpoint.ServiceId,
				refPlanService.PlanId,
				refPlanService.ServiceId),
			AttributePath: nil,
		}}
	}

	return masherytypes.MasheryPlanServiceEndpoint{
		MasheryPlanService: masherytypes.MasheryPlanService{
			PackageId: refPlanService.PackageId,
			PlanId:    refPlanService.PlanId,
			ServiceId: refPlanService.ServiceId,
		},
		EndpointId: planEndpoint.EndpointId,
	}, diag.Diagnostics{}
}
