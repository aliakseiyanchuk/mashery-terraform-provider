package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashPlanServiceId = "plan_service_id"
)

var PlanServiceEndpointMapper *PLanServiceEndpointMapperImpl

type PLanServiceEndpointMapperImpl struct {
	ResourceMapperImpl
}

func (pem *PLanServiceEndpointMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.PackagePlanServiceEndpointIdentifier, V3ObjectIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}

	plnService := masherytypes.PackagePlanServiceIdentifier{}
	if !CompoundIdFrom(&plnService, ExtractString(d, MashPlanServiceId, "")) {
		rvd = append(rvd, pem.lackingIdentificationDiagnostic(MashPlanServiceId))
	}

	endpointIdent := masherytypes.ServiceEndpointIdentifier{}
	if !CompoundIdFrom(&endpointIdent, ExtractString(d, MashEndpointId, "")) {
		rvd = append(rvd, pem.lackingIdentificationDiagnostic(MashEndpointId))
	}

	rv := masherytypes.PackagePlanServiceEndpointIdentifier{
		PackagePlanIdentifier:     plnService.PackagePlanIdentifier,
		ServiceEndpointIdentifier: endpointIdent,
	}

	return rv, nil, rvd
}

func init() {
	PlanServiceEndpointMapper = &PLanServiceEndpointMapperImpl{
		ResourceMapperImpl{
			v3ObjectName: "package plan service endpoint",
			schema: map[string]*schema.Schema{
				MashPlanServiceId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Plan service",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.PackagePlanServiceIdentifier{}
						})
					},
				},
				MashEndpointId: {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Endpoint to include",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.ServiceEndpointIdentifier{}
						})
					},
				},
			},

			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv, _, rvd := PlanServiceEndpointMapper.UpsertableTyped(d)
				return rv, rvd
			},

			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return PlanServiceEndpointMapper.UpsertableTyped(d)
			},
		},
	}
}
