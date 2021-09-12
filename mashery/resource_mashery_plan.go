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

func resourceMasheryPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: planCreate,
		ReadContext:   planRead,
		UpdateContext: planUpdate,
		DeleteContext: planDelete,
		Schema:        PlanSchema,
		Importer: &schema.ResourceImporter{
			StateContext: importMasheryPackagePlan,
		},
	}
}

func importMasheryPackagePlan(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)

	planKey := PlanIdentifier{}
	planKey.From(d.Id())

	serviceId := d.Get(MashPackagekId).(string)
	if serviceId != planKey.PackageId {
		return nil, errors.New(fmt.Sprintf("Conflict between referring packageId=%s and passed packageId argument=%s", serviceId, planKey.PackageId))
	}

	if plan, err := mashV3Cl.GetPlan(ctx, planKey.PackageId, planKey.PlanId); err != nil {
		return []*schema.ResourceData{}, errwrap.Wrapf("Failed to import this plan: {{err}}", err)
	} else if plan == nil {
		return []*schema.ResourceData{}, errors.New("No such plan found in this package")
	} else {
		V3PlanToResourceData(plan, d)
		return []*schema.ResourceData{d}, nil
	}
}

func planCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	packId := d.Get(MashPackagekId).(string)
	upsert := V3PlanUpsertable(d)

	v3Client := m.(v3client.Client)
	if rv, err := v3Client.CreatePlan(ctx, packId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		doLogJson("Received successful response from create-plan", &rv)
		V3PlanToResourceData(rv, d)

		planIdent := PlanIdentifier{
			PackageId: packId,
			PlanId:    rv.Id,
		}
		d.SetId(planIdent.Id())
		return diag.Diagnostics{}
	}
}

func planRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	plnIdent := PlanIdentifier{}
	plnIdent.From(d.Id())

	doLogf("Reading package %s plan %s", plnIdent.PackageId, plnIdent.PlanId)

	v3Client := m.(v3client.Client)
	if rv, err := v3Client.GetPlan(ctx, plnIdent.PackageId, plnIdent.PlanId); err != nil {
		return diag.FromErr(err)
	} else {
		if rv != nil {
			doLogJson("Received plan object", rv)
			V3PlanToResourceData(rv, d)
		} else {
			doLogf("Package %s and plan received nil plan object.", plnIdent.PackageId)
			d.SetId("")
		}

		return diag.Diagnostics{}
	}
}

func planUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	upsert := V3PlanUpsertable(d)

	v3Client := m.(v3client.Client)
	if rv, err := v3Client.UpdatePlan(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		if rv != nil {
			V3PlanToResourceData(rv, d)
		} else {
			d.SetId("")
		}

		return diag.Diagnostics{}
	}
}

func planDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	plnIdent := PlanIdentifier{}
	plnIdent.From(d.Id())

	v3Client := m.(v3client.Client)
	if err := v3Client.DeletePlan(ctx, plnIdent.PackageId, plnIdent.PlanId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
