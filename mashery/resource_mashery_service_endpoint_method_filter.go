package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mashschema "terraform-provider-mashery/mashschema"
)

func resourceMasheryEndpointMethodFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext:   endpointMethodFilterRead,
		CreateContext: endpointMethodFilterCreate,
		UpdateContext: endpointMethodFilterUpdate,
		DeleteContext: endpointMethodFilterDelete,
		Schema:        mashschema.ServiceEndpointMethodFilterMapper.TerraformSchema(),
	}
}

func endpointMethodFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emid := mashschema.ServiceEndpointMethodFilterMapper.GetIdentifier(d)

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetEndpointMethodFilter(ctx, emid.ServiceId, emid.EndpointId, emid.MethodId, emid.FilterId); err != nil {
		return diag.FromErr(err)
	} else {
		return mashschema.ServiceEndpointMethodFilterMapper.PersistTyped(ctx, rv, d)
	}
}

func endpointMethodFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	eid, dg := mashschema.ServiceEndpointMethodFilterMapper.GetMethodIdentifier(d)
	if len(dg) > 0 {
		return dg
	}

	v3cl := m.(v3client.Client)
	upsert, dg := mashschema.ServiceEndpointMethodFilterMapper.UpsertableTyped(d)
	if len(dg) > 0 {
		return dg
	}

	if rv, err := v3cl.CreateEndpointMethodFilter(ctx, eid.ServiceId, eid.EndpointId, eid.MethodId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.ServiceEndpointMethodFilterMapper.SetIdentifier(eid, rv, d)
		return mashschema.ServiceEndpointMethodFilterMapper.PersistTyped(ctx, rv, d)
	}
}

func endpointMethodFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emfi := mashschema.ServiceEndpointMethodFilterMapper.GetIdentifier(d)

	v3cl := m.(v3client.Client)
	upsert, dg := mashschema.ServiceEndpointMethodFilterMapper.UpsertableTyped(d)

	if len(dg) > 0 {
		return dg
	}

	if rv, err := v3cl.UpdateEndpointMethodFilter(ctx, emfi.ServiceId, emfi.EndpointId, emfi.MethodId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		return mashschema.ServiceEndpointMethodFilterMapper.PersistTyped(ctx, rv, d)
	}
}

func endpointMethodFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emfi := mashschema.ServiceEndpointMethodFilterMapper.GetIdentifier(d)

	v3cl := m.(v3client.Client)
	if err := v3cl.DeleteEndpointMethodFilter(ctx, emfi.ServiceId, emfi.EndpointId, emfi.MethodId, emfi.FilterId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
