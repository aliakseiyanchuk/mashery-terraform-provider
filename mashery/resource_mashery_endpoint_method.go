package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryEndpointMethod() *schema.Resource {
	return &schema.Resource{
		ReadContext:   endpointMethodRead,
		CreateContext: endpointMethodCreate,
		UpdateContext: endpointMethodUpdate,
		DeleteContext: endpointMethodDelete,
		Schema:        EndpointMethodSchema,
	}
}

func endpointMethodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := ServiceEndpointMethodIdentifier{}
	emid.From(d.Id())

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetEndpointMethod(ctx, emid.ServiceId, emid.EndpointId, emid.MethodId); err != nil {
		return diag.FromErr(err)
	} else {
		return V3EndpointMethodToResourceState(rv, d)
	}
}

func endpointMethodCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	eid := ServiceEndpointIdentifier{}
	eid.From(extractString(d, MashEndpointId, ""))

	if !eid.IsIdentified() {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Incomplete identifier",
			Detail:        "Endpoint identifier supplies incomplete data, or is malformed",
			AttributePath: cty.GetAttrPath(MashEndpointId),
		}}
	}

	v3cl := m.(v3client.Client)
	upsert := MashEndpointMethodUpsertable(d)

	if rv, err := v3cl.CreateEndpointMethod(ctx, eid.ServiceId, eid.EndpointId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		mid := ServiceEndpointMethodIdentifier{
			ServiceEndpointIdentifier: ServiceEndpointIdentifier{
				ServiceId:  eid.ServiceId,
				EndpointId: eid.EndpointId,
			},
			MethodId: rv.Id,
		}

		d.SetId(mid.Id())
		return V3EndpointMethodToResourceState(rv, d)
	}
}

func endpointMethodUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := ServiceEndpointMethodIdentifier{}
	emid.From(d.Id())

	v3cl := m.(v3client.Client)
	upsert := MashEndpointMethodUpsertable(d)

	if rv, err := v3cl.UpdateEndpointMethod(ctx, emid.ServiceId, emid.EndpointId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		return V3EndpointMethodToResourceState(rv, d)
	}
}

func endpointMethodDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := ServiceEndpointMethodIdentifier{}
	emid.From(d.Id())

	v3cl := m.(v3client.Client)
	if err := v3cl.DeleteEndpointMethod(ctx, emid.ServiceId, emid.EndpointId, emid.MethodId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
