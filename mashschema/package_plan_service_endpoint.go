package mashschema

import (
	"context"
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

type PlanServiceEndpointIdentifier struct {
	PlanServiceIdentifier
	EndpointId string
}

func (pei *PlanServiceEndpointIdentifier) Self() interface{} {
	return pei
}

var PlanServiceEndpointMapper *PLanServiceEndpointMapperImpl

type PLanServiceEndpointMapperImpl struct {
	MapperImpl
}

func (pem *PLanServiceEndpointMapperImpl) SetIdentifier(d *schema.ResourceData) {
	psiIdent := pem.GetPlanServiceIdentifier(d)
	endIdent := pem.GetServiceEndpointIdentifier(d)

	ident := &PlanServiceEndpointIdentifier{
		PlanServiceIdentifier: PlanServiceIdentifier{
			PlanIdentifier: PlanIdentifier{
				PackageIdentifier: PackageIdentifier{
					PackageId: psiIdent.PackageId,
				},
				PlanId: psiIdent.PlanId,
			},
			ServiceId: psiIdent.ServiceId,
		},

		EndpointId: endIdent.EndpointId,
	}

	d.SetId(CompoundId(ident))
}

func (pem *PLanServiceEndpointMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.MasheryPlanServiceEndpoint, diag.Diagnostics) {
	refPlanService := pem.GetPlanServiceIdentifier(d)

	planEndpoint := pem.GetServiceEndpointIdentifier(d)

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

func (pem *PLanServiceEndpointMapperImpl) GetServiceEndpointIdentifier(d *schema.ResourceData) *ServiceEndpointIdentifier {
	planEndpoint := &ServiceEndpointIdentifier{}
	CompoundIdFrom(&planEndpoint, ExtractString(d, MashEndpointId, ""))
	return planEndpoint
}

func (pem *PLanServiceEndpointMapperImpl) GetPlanServiceIdentifier(d *schema.ResourceData) *PlanServiceIdentifier {
	refPlanService := &PlanServiceIdentifier{}
	CompoundIdFrom(&refPlanService, ExtractString(d, MashPlanServiceId, ""))
	return refPlanService
}

func init() {
	PlanServiceEndpointMapper = &PLanServiceEndpointMapperImpl{
		MapperImpl{
			schema: map[string]*schema.Schema{
				MashPlanServiceId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Plan service",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &PlanServiceIdentifier{}
						})
					},
				},
				MashEndpointId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Endpoint to include",
				},
			},
		},
	}

	PlanServiceEndpointMapper.identifier = func() interface{} {
		return &PlanServiceEndpointIdentifier{}
	}
	PlanServiceEndpointMapper.upsertFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return PlanServiceEndpointMapper.UpsertableTyped(d)
	}
}
