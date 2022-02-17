package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PlanEndpointMethodIdentifier struct {
	PlanServiceEndpointIdentifier
	MethodId string
}

type PackagePlanEndpointMethodFilterIdentifier struct {
	PlanEndpointMethodIdentifier
	FilterId string
}

const PackagePlanEndpointRef = "package_plan_service_endpoint_id"
const PlanEndpointMethodFilterId = "service_endpoint_method_id"
const PlanEndpointMethodId = "service_endpoint_method_filter_id"

func (pemi *PlanEndpointMethodIdentifier) Self() interface{} {
	return pemi
}

func (pemi *PackagePlanEndpointMethodFilterIdentifier) Self() interface{} {
	return pemi
}

var PackagePlanServiceEndpointMethodMapper *PackagePlanServiceEndpointMethodMapperImpl

type PackagePlanServiceEndpointMethodMapperImpl struct {
	MapperImpl
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) SetIdentifier(endpointIdent *masherytypes.MasheryPlanServiceEndpoint, meth *masherytypes.MasheryMethod, d *schema.ResourceData) {
	rid := &PlanEndpointMethodIdentifier{
		PlanServiceEndpointIdentifier: PlanServiceEndpointIdentifier{
			PlanServiceIdentifier: PlanServiceIdentifier{
				PlanIdentifier: PlanIdentifier{
					PackageIdentifier: PackageIdentifier{
						PackageId: endpointIdent.PackageId,
					},
					PlanId: endpointIdent.PlanId,
				},
				ServiceId: endpointIdent.ServiceId,
			},
			EndpointId: endpointIdent.EndpointId,
		},
		MethodId: meth.Id,
	}

	d.SetId(CompoundId(rid))
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) GetIdentifiers(d *schema.ResourceData) (*masherytypes.MasheryPlanServiceEndpoint, *ServiceEndpointMethodIdentifier, *ServiceEndpointMethodFilterIdentifier, diag.Diagnostics) {

	rvDiags := diag.Diagnostics{}

	endpointIdent := psem.GetPackageEndpointIdentifier(d)
	methIdent := psem.ServiceMethodIdentifier(d)

	// Detect a collision
	if methIdent.EndpointId != endpointIdent.EndpointId {
		rvDiags = append(rvDiags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Mismatching object",
			Detail:        "Method does not belong to the same endpoint",
			AttributePath: cty.GetAttrPath(ServiceEndpointMethodRef),
		})
	}

	filterIdentStr := ExtractString(d, ServiceEndpointMethodFilterRef, "")
	var uFilter *ServiceEndpointMethodFilterIdentifier = nil

	if len(filterIdentStr) > 0 {
		filterIdent := ServiceEndpointMethodFilterIdentifier{}
		CompoundIdFrom(filterIdent, filterIdentStr)

		if methIdent.MethodId != filterIdent.MethodId {
			rvDiags = append(rvDiags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Mismatching object",
				Detail:        "Filter does not belong to the containing method",
				AttributePath: cty.GetAttrPath(ServiceEndpointMethodFilterRef),
			})
		}

		uFilter = &filterIdent
	}

	return endpointIdent, methIdent, uFilter, rvDiags

}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) HasFilterChange(d *schema.ResourceData) bool {
	return d.HasChange(ServiceEndpointMethodFilterRef)
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) GetFilterChange(d *schema.ResourceData) (*ServiceEndpointMethodFilterIdentifier, *ServiceEndpointMethodFilterIdentifier) {
	before, after := d.GetChange(ServiceEndpointMethodFilterRef)

	rvBefore := &ServiceEndpointMethodFilterIdentifier{}
	rvAfter := &ServiceEndpointMethodFilterIdentifier{}

	CompoundIdFrom(rvBefore, before.(string))
	CompoundIdFrom(rvAfter, after.(string))

	return rvBefore, rvAfter
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) ClearFilter(d *schema.ResourceData) {
	d.Set(ServiceEndpointMethodFilterRef, "")
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) SetServiceFilterIdent(endpIdent *masherytypes.MasheryPlanServiceEndpointMethod, filter *masherytypes.MasheryResponseFilter, d *schema.ResourceData) {
	rv := &ServiceEndpointMethodFilterIdentifier{
		ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
			ServiceEndpointIdentifier: ServiceEndpointIdentifier{
				ServiceIdentifier: ServiceIdentifier{
					ServiceId: endpIdent.ServiceId},
				EndpointId: endpIdent.EndpointId,
			},
			MethodId: endpIdent.MethodId,
		},
		FilterId: filter.Id,
	}
	d.Set(ServiceEndpointMethodFilterRef, CompoundId(rv))
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) GetPackageEndpointIdentifier(d *schema.ResourceData) *masherytypes.MasheryPlanServiceEndpoint {
	ident := PlanServiceEndpointIdentifier{}
	CompoundIdFrom(&ident, ExtractString(d, PackagePlanEndpointRef, ""))

	return &masherytypes.MasheryPlanServiceEndpoint{
		MasheryPlanService: masherytypes.MasheryPlanService{
			PackageId: ident.PackageId,
			PlanId:    ident.PlanId,
			ServiceId: ident.ServiceId,
		},
		EndpointId: ident.EndpointId,
	}
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) GetIdentifier(d *schema.ResourceData) (masherytypes.MasheryPlanServiceEndpointMethod, diag.Diagnostics) {
	ident := PlanEndpointMethodIdentifier{}
	CompoundIdFrom(&ident, d.Id())

	v3Id := masherytypes.MasheryPlanServiceEndpointMethod{
		MasheryPlanServiceEndpoint: masherytypes.MasheryPlanServiceEndpoint{
			MasheryPlanService: masherytypes.MasheryPlanService{
				PackageId: ident.PackageId,
				PlanId:    ident.PlanId,
				ServiceId: ident.ServiceId,
			},
			EndpointId: ident.EndpointId,
		},
		MethodId: ident.MethodId,
	}

	rvd := diag.Diagnostics{}
	if !IsIdentified(&ident) {
		rvd = diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Incomplete identifier",
			Detail:        "Plan method identifier is incomplete",
			AttributePath: cty.GetAttrPath("id"),
		}}
	}

	return v3Id, rvd
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) ServiceMethodFilterIdentifier(d *schema.ResourceData) (*masherytypes.MasheryServiceMethodFilter, diag.Diagnostics) {
	ident := ServiceEndpointMethodFilterIdentifier{}
	CompoundIdFrom(&ident, ExtractString(d, ServiceEndpointMethodFilterRef, ""))

	v3Id := masherytypes.MasheryServiceMethodFilter{
		MasheryServiceMethod: masherytypes.MasheryServiceMethod{
			MasheryServiceEndpoint: masherytypes.MasheryServiceEndpoint{
				ServiceId:  ident.ServiceId,
				EndpointId: ident.EndpointId,
			},
			MethodId: ident.MethodId,
		},
		FilterId: ident.FilterId,
	}

	var rvd diag.Diagnostics = nil

	if !IsIdentified(&ident) {
		rvd = diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Incomplete identifier",
			Detail:        "Service method filter identifier is incomplete",
			AttributePath: cty.GetAttrPath("id"),
		}}
	}

	return &v3Id, rvd
}

func (psem *PackagePlanServiceEndpointMethodMapperImpl) ServiceMethodIdentifier(d *schema.ResourceData) *ServiceEndpointMethodIdentifier {
	ident := &ServiceEndpointMethodIdentifier{}
	CompoundIdFrom(&ident, ExtractString(d, ServiceEndpointMethodRef, ""))

	return ident
}

func initPlanMethodSchemaBoilerplate() {
	addComputedString(&PackagePlanServiceEndpointMethodMapper.schema, MashObjId, "Object identifier")
	addComputedString(&PackagePlanServiceEndpointMethodMapper.schema, MashObjCreated, "Date/time the object was created")
	addComputedString(&PackagePlanServiceEndpointMethodMapper.schema, MashObjUpdated, "Date/time the object was created")

	addComputedString(&PackagePlanServiceEndpointMethodMapper.schema, PlanEndpointMethodFilterId, "Package filter UUID identifier")
	addComputedString(&PackagePlanServiceEndpointMethodMapper.schema, PlanEndpointMethodId, "Package method UUID identifier")
}

func init() {
	PackagePlanServiceEndpointMethodMapper = &PackagePlanServiceEndpointMethodMapperImpl{
		MapperImpl{
			schema: map[string]*schema.Schema{
				PackagePlanEndpointRef: {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						ident := PlanServiceEndpointIdentifier{}
						CompoundIdFrom(&ident, i.(string))

						rv := diag.Diagnostics{}

						if !IsIdentified(&ident) {
							rv = append(rv, diag.Diagnostic{
								Severity:      diag.Error,
								Summary:       "Incomplete identifier",
								Detail:        "Endpoint reference is incomplete or malformed",
								AttributePath: path,
							})
						}

						return rv
					},
				},
				ServiceEndpointMethodRef: {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						ident := ServiceEndpointMethodIdentifier{}
						CompoundIdFrom(&ident, i.(string))
						if !IsIdentified(&ident) {
							return diag.Diagnostics{diag.Diagnostic{
								Severity:      diag.Error,
								Summary:       "Malformed identifier",
								Detail:        "Package pan endpoint identifier is incomplete or malformed. Ensure that you are passing in the id attribute of mashery_package_endpoint resource.",
								AttributePath: path,
							}}
						} else {
							return diag.Diagnostics{}
						}
					},
				},
				ServiceEndpointMethodFilterRef: {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						filterId := i.(string)

						if len(filterId) > 0 {
							ident := ServiceEndpointMethodFilterIdentifier{}
							CompoundIdFrom(&ident, filterId)

							if !IsIdentified(&ident) {
								return diag.Diagnostics{diag.Diagnostic{
									Severity:      diag.Error,
									Summary:       "Malformed identifier",
									Detail:        "Endpoint method identifier is incomplete or malformed. Ensure that you are passing in the id attribute of mashery_endpoint_method_filter resource.",
									AttributePath: path,
								}}
							} else {
								return diag.Diagnostics{}
							}
						} else {
							return diag.Diagnostics{}
						}
					},
				},
			},
		},
	}

	initPlanMethodSchemaBoilerplate()
}
