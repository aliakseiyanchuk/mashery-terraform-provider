package mashery

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryEndpoint() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   endpointRead,
		CreateContext: endpointCreate,
		UpdateContext: endpointUpdate,
		DeleteContext: endpointDelete,
		// Schema
		Schema: EndpointSchema,
		// Importer by ID
		Importer: &schema.ResourceImporter{
			StateContext: importMasheryServiceEndpoint,
		},
	}
}

func importMasheryServiceEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)

	endpKey := ServiceEndpointIdentifier{}
	endpKey.From(d.Id())

	serviceId := d.Get(MashSvcId).(string)
	if serviceId != endpKey.ServiceId {
		return nil, errors.New(fmt.Sprintf("Conflict between referring serviceId=%s and passed serviceId argument=%s", serviceId, endpKey.ServiceId))
	}

	if endp, err := mashV3Cl.GetEndpoint(ctx, endpKey.ServiceId, endpKey.EndpointId); err != nil {
		return []*schema.ResourceData{}, errwrap.Wrapf("Failed to import this endpoint: {{err}}", err)
	} else if endp == nil {
		return []*schema.ResourceData{}, errors.New("No such endpoint found in this service")
	} else {
		V3EndpointToResourceData(endp, d)
		return []*schema.ResourceData{d}, nil
	}
}

func endpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpKey := ServiceEndpointIdentifier{}
	endpKey.From(d.Id())

	v3Client := m.(v3client.Client)
	if rv, err := v3Client.GetEndpoint(ctx, endpKey.ServiceId, endpKey.EndpointId); err != nil {
		return diag.FromErr(err)
	} else {
		if rv != nil {
			doLogJson(fmt.Sprintf("Read the following endpoint data belonging to service %s endpoint %s", endpKey.ServiceId, endpKey.EndpointId), &rv)
			V3EndpointToResourceData(rv, d)
		} else {
			// Object no longer exists.
			doLogf("Service %s no longer contains endpoint %s", endpKey.ServiceId, endpKey.EndpointId)
			d.SetId("")
		}
	}

	return diag.Diagnostics{}
}

func endpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serviceId := d.Get(MashSvcId).(string)

	upsert, rvDiag := MashEndpointUpsertable(d)

	if len(rvDiag) > 0 {
		return rvDiag
	} else {
		v3Client := m.(v3client.Client)
		if rv, err := v3Client.CreateEndpoint(ctx, serviceId, upsert); err != nil {
			return diag.FromErr(err)
		} else {
			V3EndpointToResourceData(rv, d)

			compoundIdent := ServiceEndpointIdentifier{
				ServiceId:  serviceId,
				EndpointId: rv.Id,
			}

			d.SetId(compoundIdent.Id())
		}
	}

	return diag.Diagnostics{}
}

func endpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpKey := ServiceEndpointIdentifier{}
	endpKey.From(d.Id())

	upsert, rvDiag := MashEndpointUpsertable(d)
	upsert.Id = endpKey.EndpointId

	doLogJson(fmt.Sprintf("Updating service %s endpoint %s", endpKey.ServiceId, endpKey), &upsert)

	if len(rvDiag) > 0 {
		return rvDiag
	} else {
		v3Client := m.(v3client.Client)
		if rv, err := v3Client.UpdateEndpoint(ctx, endpKey.ServiceId, upsert); err != nil {
			return diag.FromErr(err)
		} else {
			V3EndpointToResourceData(rv, d)
		}
	}

	return diag.Diagnostics{}
}

func endpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpKey := ServiceEndpointIdentifier{}
	endpKey.From(d.Id())

	v3Client := m.(v3client.Client)

	if err := v3Client.DeleteEndpoint(ctx, endpKey.ServiceId, endpKey.EndpointId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
