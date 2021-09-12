package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryEndpointMethodFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext:   endpointMethodFilterRead,
		CreateContext: endpointMethodFilterCreate,
		UpdateContext: endpointMethodFilterUpdate,
		DeleteContext: endpointMethodFilterDelete,
		Schema:        EndpointMethodFilterSchema,
	}
}

func endpointMethodFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := ServiceEndpointMethodFilterIdentifier{}
	emid.From(d.Id())

	if !emid.IsIdentified() {
		return emid.MalformedDiagnostic("id")
	}

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetEndpointMethodFilter(ctx, emid.ServiceId, emid.EndpointId, emid.MethodId, emid.FilterId); err != nil {
		return diag.FromErr(err)
	} else {
		return V3EndpointMethodFilterToResourceData(rv, d)
	}
}

func endpointMethodFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	eid := ServiceEndpointMethodIdentifier{}
	eid.From(extractString(d, MashServiceEndpointMethodRef, ""))

	if !eid.IsIdentified() {
		return eid.MalformedDiagnostic(MashServiceEndpointMethodRef)
	}

	v3cl := m.(v3client.Client)
	upsert, dg := MashEndpointMethodFilterUpsertable(d)
	if len(dg) > 0 {
		return dg
	}

	if rv, err := v3cl.CreateEndpointMethodFilter(ctx, eid.ServiceId, eid.EndpointId, eid.MethodId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		mid := ServiceEndpointMethodFilterIdentifier{
			ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
				ServiceEndpointIdentifier: ServiceEndpointIdentifier{
					ServiceId:  eid.ServiceId,
					EndpointId: eid.EndpointId,
				},
				MethodId: eid.MethodId,
			},
			FilterId: rv.Id,
		}

		d.SetId(mid.Id())
		return V3EndpointMethodFilterToResourceData(rv, d)
	}
}

func endpointMethodFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emfi := ServiceEndpointMethodFilterIdentifier{}
	emfi.From(d.Id())

	if !emfi.IsIdentified() {
		return emfi.MalformedDiagnostic("id")
	}

	v3cl := m.(v3client.Client)
	upsert, dg := MashEndpointMethodFilterUpsertable(d)

	if len(dg) > 0 {
		return dg
	}

	if rv, err := v3cl.UpdateEndpointMethodFilter(ctx, emfi.ServiceId, emfi.EndpointId, emfi.MethodId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		return V3EndpointMethodFilterToResourceData(rv, d)
	}
}

func endpointMethodFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emfi := ServiceEndpointMethodFilterIdentifier{}
	emfi.From(d.Id())

	if !emfi.IsIdentified() {
		return emfi.MalformedDiagnostic("id")
	}

	v3cl := m.(v3client.Client)
	if err := v3cl.DeleteEndpointMethodFilter(ctx, emfi.ServiceId, emfi.EndpointId, emfi.MethodId, emfi.FilterId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
