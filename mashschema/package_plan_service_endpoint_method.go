package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const PackagePlanEndpointRef = "package_plan_service_endpoint_id"
const PackagePlanEndpointMethodRef = "package_plan_service_endpoint_method_id"
const PlanEndpointMethodFilterId = "service_endpoint_method_id"
const PlanEndpointMethodId = "service_endpoint_method_filter_id"

var PackagePlanServiceEndpointMethodMapper *PackagePlanServiceEndpointMethodMapperImpl

type PackagePlanServiceEndpointMethodMapperImpl struct {
	ResourceMapperImpl
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.PackagePlanServiceEndpointMethodIdentifier, V3ObjectIdentifier, diag.Diagnostics) {
	main, rvd := psem.V3Identity(d)
	return main, nil, rvd
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) V3Identity(d *schema.ResourceData) (masherytypes.PackagePlanServiceEndpointMethodIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}

	// Get the passed service endpoint method
	methIdent := masherytypes.ServiceEndpointMethodIdentifier{}
	if !CompoundIdFrom(&methIdent, ExtractString(d, ServiceEndpointMethodRef, "")) {
		rvd = append(rvd, psem.lackingIdentificationDiagnostic(ServiceEndpointMethodRef))
	}

	planEndpointIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	if !CompoundIdFrom(&planEndpointIdent, ExtractString(d, PackagePlanEndpointRef, "")) {
		rvd = append(rvd, psem.lackingIdentificationDiagnostic(PackagePlanEndpointRef))
	}

	main := masherytypes.PackagePlanServiceEndpointMethodIdentifier{
		PackagePlanIdentifier:           planEndpointIdent.PackagePlanIdentifier,
		ServiceEndpointMethodIdentifier: methIdent,
	}

	return main, rvd
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) HasFilterChange(d *schema.ResourceData) bool {
	return d.HasChange(ServiceEndpointMethodFilterRef)
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) HasFilter(d *schema.ResourceData) bool {
	return d.Get(ServiceEndpointMethodFilterRef) != nil
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) GetFilterChange(d *schema.ResourceData) (*masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier, *masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier) {
	before, after := d.GetChange(ServiceEndpointMethodFilterRef)

	rvBefore := &masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{}
	rvAfter := &masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{}

	CompoundIdFrom(rvBefore, before.(string))
	CompoundIdFrom(rvAfter, after.(string))

	return rvBefore, rvAfter
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) GetFilterIdentity(d *schema.ResourceData) (masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}

	// Get the passed service endpoint method
	filterIdent := masherytypes.ServiceEndpointMethodFilterIdentifier{}
	if !CompoundIdFrom(&filterIdent, ExtractString(d, ServiceEndpointMethodFilterRef, "")) {
		rvd = append(rvd, psem.lackingIdentificationDiagnostic(ServiceEndpointMethodFilterRef))
	}

	planEndpointIdent := masherytypes.PackagePlanServiceEndpointIdentifier{}
	if !CompoundIdFrom(&planEndpointIdent, ExtractString(d, PackagePlanEndpointRef, "")) {
		rvd = append(rvd, psem.lackingIdentificationDiagnostic(PackagePlanEndpointRef))
	}

	main := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{
		PackagePlanIdentifier:                 planEndpointIdent.PackagePlanIdentifier,
		ServiceEndpointMethodFilterIdentifier: filterIdent,
	}

	return main, rvd
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) ClearFilter(d *schema.ResourceData) {
	d.Set(ServiceEndpointMethodFilterRef, "")
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) SetServiceFilterIdent(filter *masherytypes.PackagePlanServiceEndpointMethodFilter, d *schema.ResourceData) {
	d.Set(ServiceEndpointMethodFilterRef, CompoundId(filter))
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) initPlanMethodSchemaBoilerplate() {
	psem.SchemaBuilder().
		AddComputedString(MashObjId, "Object v3Identity").
		AddComputedString(MashObjCreated, "Date/time the object was created").
		AddComputedString(MashObjUpdated, "Date/time the object was created").
		AddComputedString(PlanEndpointMethodFilterId, "Package filter UUID v3Identity").
		AddComputedString(PlanEndpointMethodId, "Package method UUID v3Identity")
}

func init() {
	PackagePlanServiceEndpointMethodMapper = &PackagePlanServiceEndpointMethodMapperImpl{
		ResourceMapperImpl{
			v3ObjectName: "package plan service endpoint method",
			schema: map[string]*schema.Schema{
				PackagePlanEndpointRef: {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.PackagePlanServiceEndpointIdentifier{}
						})
					},
				},
				ServiceEndpointMethodRef: {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.ServiceEndpointMethodIdentifier{}
						})
					},
				},
				ServiceEndpointMethodFilterRef: {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &masherytypes.ServiceEndpointMethodFilterIdentifier{}
						})
					},
				},
			},

			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				return PackagePlanServiceEndpointMethodMapper.V3Identity(d)
			},
		},
	}

	PackagePlanServiceEndpointMethodMapper.initPlanMethodSchemaBoilerplate()
}
