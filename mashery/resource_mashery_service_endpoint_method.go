package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryEndpointMethod() *schema.Resource {
	return &schema.Resource{
		ReadContext:   endpointMethodRead,
		CreateContext: endpointMethodCreate,
		UpdateContext: endpointMethodUpdate,
		DeleteContext: endpointMethodDelete,
		Schema:        mashschema.ServiceEndpointMethodMapper.TerraformSchema(),
	}
}

func endpointMethodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := mashschema.ServiceEndpointMethodMapper.CreateIdentifier(d)

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetEndpointMethod(ctx, emid.ServiceId, emid.EndpointId, emid.MethodId); err != nil {
		return diag.FromErr(err)
	} else {
		return mashschema.ServiceEndpointMethodMapper.PersistTyped(ctx, rv, d)
	}
}

func endpointMethodCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if eid, dg := mashschema.ServiceEndpointMethodMapper.EndpointIdentifier(d); len(dg) > 0 {
		return dg
	} else {

		v3cl := m.(v3client.Client)
		upsert, _ := mashschema.ServiceEndpointMethodMapper.UpsertableTyped(d)

		if rv, err := v3cl.CreateEndpointMethod(ctx, eid.ServiceId, eid.EndpointId, *upsert); err != nil {
			return diag.FromErr(err)
		} else {
			mid := mashschema.ServiceEndpointMethodMapper.CreateIdentifierTyped()
			mid.ServiceId = eid.ServiceId
			mid.EndpointId = eid.EndpointId
			mid.MethodId = rv.Id

			d.SetId(mashschema.CompoundId(mid))
			return mashschema.ServiceEndpointMethodMapper.PersistTyped(ctx, rv, d)
		}
	}
}

func endpointMethodUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := mashschema.ServiceEndpointMethodMapper.CreateIdentifier(d)

	v3cl := m.(v3client.Client)
	upsert, _ := mashschema.ServiceEndpointMethodMapper.UpsertableTyped(d)

	if rv, err := v3cl.UpdateEndpointMethod(ctx, emid.ServiceId, emid.EndpointId, *upsert); err != nil {
		return diag.FromErr(err)
	} else {
		return mashschema.ServiceEndpointMethodMapper.PersistTyped(ctx, rv, d)
	}
}

func endpointMethodDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := mashschema.ServiceEndpointMethodMapper.CreateIdentifier(d)

	v3cl := m.(v3client.Client)
	if err := v3cl.DeleteEndpointMethod(ctx, emid.ServiceId, emid.EndpointId, emid.MethodId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
