package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PlanEndpointMethodIdentifier struct {
	PlanEndpointIdentifier
	MethodId string
}

type PlanEndpointMethodFilterIdentifier struct {
	PlanEndpointMethodIdentifier
	FilterId string
}

const PlanEndpointRef = "package_endpoint_id"
const PlanEndpointMethodFilterId = "package_filter_id"
const PlanEndpointMethodId = "package_method_id"

func (pemi *PlanEndpointMethodIdentifier) Id() string {
	return CreateCompoundId(pemi.PackageId, pemi.PlanId, pemi.ServiceId, pemi.EndpointId, pemi.MethodId)
}

func (pemi *PlanEndpointMethodIdentifier) From(id string) {
	ParseCompoundId(id, &pemi.PackageId, &pemi.PlanId, &pemi.ServiceId, &pemi.EndpointId, &pemi.MethodId)
}

func (pemi *PlanEndpointMethodIdentifier) IsIdentified() bool {
	return pemi.PlanEndpointIdentifier.IsIdentified() && len(pemi.MethodId) > 0
}

func MashPlanServiceEndpointUpsert(d *schema.ResourceData) v3client.MasheryPlanServiceEndpoint {
	ident := PlanEndpointIdentifier{}
	ident.From(extractString(d, PlanEndpointRef, ""))

	return v3client.MasheryPlanServiceEndpoint{
		MasheryPlanService: v3client.MasheryPlanService{
			PackageId: ident.PackageId,
			PlanId:    ident.PlanId,
			ServiceId: ident.ServiceId,
		},
		EndpointId: ident.EndpointId,
	}
}

func MashPlanMethodIdentifier(d *schema.ResourceData) (v3client.MasheryPlanServiceEndpointMethod, diag.Diagnostics) {
	ident := PlanEndpointMethodIdentifier{}
	ident.From(d.Id())

	v3Id := v3client.MasheryPlanServiceEndpointMethod{
		MasheryPlanServiceEndpoint: v3client.MasheryPlanServiceEndpoint{
			MasheryPlanService: v3client.MasheryPlanService{
				PackageId: ident.PackageId,
				PlanId:    ident.PlanId,
				ServiceId: ident.ServiceId,
			},
			EndpointId: ident.EndpointId,
		},
		MethodId: ident.MethodId,
	}

	rvd := diag.Diagnostics{}
	if !ident.IsIdentified() {
		rvd = diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Incomplete identifier",
			Detail:        "Plan method identifier is incomplete",
			AttributePath: cty.GetAttrPath("id"),
		}}
	}

	return v3Id, rvd
}

func MashServiceMethodIdentifier(d *schema.ResourceData) (v3client.MasheryServiceMethodFilter, diag.Diagnostics) {
	ident := ServiceEndpointMethodFilterIdentifier{}
	ident.From(extractString(d, MashServiceEndpointMethodFilterRef, ""))

	v3Id := v3client.MasheryServiceMethodFilter{
		MasheryServiceMethod: v3client.MasheryServiceMethod{
			MasheryServiceEndpoint: v3client.MasheryServiceEndpoint{
				ServiceId:  ident.ServiceId,
				EndpointId: ident.EndpointId,
			},
			MethodId: ident.MethodId,
		},
		FilterId: ident.FilterId,
	}

	rvd := diag.Diagnostics{}
	if !ident.IsIdentified() {
		rvd = diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Incomplete identifier",
			Detail:        "Service method filter identifier is incomplete",
			AttributePath: cty.GetAttrPath("id"),
		}}
	}

	return v3Id, rvd
}

var PlanMethodSchema = map[string]*schema.Schema{
	PlanEndpointRef: {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			ident := PlanEndpointIdentifier{}
			ident.From(i.(string))

			rv := diag.Diagnostics{}

			if !ident.IsIdentified() {
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
	MashServiceEndpointMethodRef: {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			ident := ServiceEndpointMethodIdentifier{}
			ident.From(i.(string))
			if !ident.IsIdentified() {
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
	MashServiceEndpointMethodFilterRef: {
		Type:     schema.TypeString,
		Optional: true,
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			filterId := i.(string)

			if len(filterId) > 0 {
				ident := ServiceEndpointMethodFilterIdentifier{}

				ident.From(filterId)
				if !ident.IsIdentified() {
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
}

func initPlanMethodSchemaBoilerplate() {
	addComputedString(&PlanMethodSchema, MashObjId, "Object identifier")
	addComputedString(&PlanMethodSchema, MashObjCreated, "Date/time the object was created")
	addComputedString(&PlanMethodSchema, MashObjUpdated, "Date/time the object was created")

	addComputedString(&PlanMethodSchema, PlanEndpointMethodFilterId, "Package filter UUID identifier")
	addComputedString(&PlanMethodSchema, PlanEndpointMethodId, "Package method UUID identifier")
}

func init() {
	initPlanMethodSchemaBoilerplate()
}
