package mashery

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryErrorSet() *schema.Resource {
	return &schema.Resource{
		Schema:        ServiceErrorSetSchema,
		ReadContext:   serviceErrorSetRead,
		CreateContext: serviceErrorSetCreate,
		UpdateContext: serviceErrorSetUpdate,
		DeleteContext: serviceErrorSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importServiceErrorSet,
		},
	}
}

func importServiceErrorSet(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)

	setId := ErrorSetIdentifier{}
	setId.From(d.Id())

	if errSet, err := mashV3Cl.GetErrorSet(ctx, setId.ServiceId, setId.ErrorSetId); err != nil {
		return []*schema.ResourceData{}, errwrap.Wrapf("Failed to import this errSet: {{err}}", err)
	} else if errSet == nil {
		return []*schema.ResourceData{}, errors.New("no such error set")
	} else {
		V3ErrorSetToResourceData(errSet, d)
		return []*schema.ResourceData{d}, nil
	}
}

func serviceErrorSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	setId := ErrorSetIdentifier{}
	setId.From(d.Id())

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetErrorSet(ctx, setId.ServiceId, setId.ErrorSetId); err != nil {
		return diag.FromErr(err)
	} else {
		return V3ErrorSetToResourceData(rv, d)
	}
}

func serviceErrorSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serviceId := d.Get(MashSvcId).(string)

	upsert := V3ErrorSetUpsertable(d)
	custMessages := V3ErrorSetMessages(d)

	retVal := diag.Diagnostics{}

	v3cl := m.(v3client.Client)
	if errSet, err := v3cl.CreateErrorSet(ctx, serviceId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		setId := ErrorSetIdentifier{
			ServiceId:  serviceId,
			ErrorSetId: errSet.Id,
		}
		d.SetId(setId.Id())

		if len(custMessages) > 0 {
			for _, msg := range custMessages {
				doLogJson("Sending error message", msg)

				if _, err := v3cl.UpdateErrorSetMessage(ctx, serviceId, errSet.Id, msg); err != nil {
					retVal = append(retVal, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Could not update error message",
						Detail:   fmt.Sprintf("%s", err),
					})
					doLogf("Failed to update error message: %s", err)
				}
			}

			if refreshed, err := v3cl.GetErrorSet(ctx, serviceId, errSet.Id); err != nil {
				retVal = append(retVal, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to refresh error set after updating messages",
					Detail:   fmt.Sprintf("%s", err),
				})
			} else {
				errSet = refreshed
			}
			setDiags := V3ErrorSetToResourceData(errSet, d)
			if len(setDiags) > 0 {
				retVal = append(retVal, setDiags...)
			}
		} else {
			// Save the basic fields
			rvDiags := V3ErrorSetToResourceData(errSet, d)
			if len(rvDiags) > 0 {
				retVal = append(retVal, rvDiags...)
			}
		}

		return retVal
	}
}

func modifiedMessages(inp *masherytypes.MasheryErrorSet, rawMsg []masherytypes.MasheryErrorMessage) []masherytypes.MasheryErrorMessage {
	var rv []masherytypes.MasheryErrorMessage

	for _, m := range rawMsg {
		if inp.ErrorMessages != nil {
			curVal := findMessageById(inp.ErrorMessages, m.Id)
			if curVal != nil &&
				curVal.DetailHeader == m.DetailHeader &&
				curVal.Status == m.Status &&
				curVal.Code == m.Code &&
				curVal.ResponseBody == m.ResponseBody {
				continue
			}
		}

		rv = append(rv, m)
	}

	return rv
}

func serviceErrorSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	setId := ErrorSetIdentifier{}
	setId.From(d.Id())

	v3cl := m.(v3client.Client)

	if err := v3cl.DeleteErrorSet(ctx, setId.ServiceId, setId.ErrorSetId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}

func serviceErrorSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	setId := ErrorSetIdentifier{}
	setId.From(d.Id())

	upsert := V3ErrorSetUpsertable(d)
	custMessages := V3ErrorSetMessages(d)

	retVal := diag.Diagnostics{}
	v3cl := m.(v3client.Client)

	var updInst *masherytypes.MasheryErrorSet
	var err error

	if updInst, err = v3cl.UpdateErrorSet(ctx, setId.ServiceId, upsert); err != nil {
		return diag.FromErr(err)
	} else if len(custMessages) > 0 && d.HasChanges(MashSvcErrorSetMessage) {
		for _, val := range modifiedMessages(updInst, custMessages) {
			if _, err := v3cl.UpdateErrorSetMessage(ctx, setId.ServiceId, setId.ErrorSetId, val); err != nil {
				retVal = append(retVal, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not update error message",
					Detail:   fmt.Sprintf("%s", err),
				})
			}
		}

		if refreshed, err := v3cl.GetErrorSet(ctx, setId.ServiceId, setId.ErrorSetId); err != nil {
			retVal = append(retVal, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "could not refresh error set after emitting custom error messagaes",
				Detail:   fmt.Sprintf("%s", err),
			})
		} else {
			updInst = refreshed
		}
	}

	setDiags := V3ErrorSetToResourceData(updInst, d)
	if len(setDiags) > 0 {
		retVal = append(retVal, setDiags...)
	}

	return retVal
}
