package mashery

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type ReaderFunc func(context.Context, v3client.Client, mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error)
type CreatorFunc func(context.Context, v3client.Client, mashschema.Upsertable, mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error)
type UpdaterFunc func(context.Context, v3client.Client, mashschema.Upsertable) (mashschema.Upsertable, error)
type DeleterFunc func(context.Context, v3client.Client, mashschema.V3ObjectIdentifier) error
type OffendersCounterFunc func(context.Context, v3client.Client, mashschema.V3ObjectIdentifier) (int64, error)

type ResourceTemplate struct {
	Mapper mashschema.ResourceMapper

	DoRead   ReaderFunc
	DoCreate CreatorFunc
	DoUpdate UpdaterFunc
	DoDelete DeleterFunc

	DoCountOffending OffendersCounterFunc
}

// TFDataSourceSchema returns the Terraform data source schema
func (t *ResourceTemplate) TFDataSourceSchema() *schema.Resource {
	// Panic if necessary functions were not supplied
	if t.Mapper == nil || t.DoCreate == nil || t.DoDelete == nil {
		panic(fmt.Sprintf("Unsatisfied initialization: mapper and all CRUD method must be supplied"))
	}

	var readCtx schema.ReadContextFunc = nil
	var updateCtx schema.UpdateContextFunc = nil

	if t.DoRead != nil {
		readCtx = t.Create
	}
	if t.DoUpdate != nil {
		updateCtx = t.Update
	}

	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: t.Import,
		},
		ReadContext:   readCtx,
		CreateContext: t.Create,
		UpdateContext: updateCtx,
		DeleteContext: t.Delete,
		Schema:        t.Mapper.TerraformSchema(),
	}
}

func (rt *ResourceTemplate) Import(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)

	if v3Id, diags := rt.Mapper.V3Identity(d); len(diags) > 0 {
		return []*schema.ResourceData{}, errors.New(fmt.Sprintf("object %s is not fully identified", rt.Mapper.V3ObjectName()))
	} else if v3Obj, err := rt.DoRead(ctx, mashV3Cl, v3Id); err != nil {
		if v3Obj != nil {
			rt.Mapper.SetState(v3Obj, d)
			return []*schema.ResourceData{d}, nil
		} else {
			return []*schema.ResourceData{}, errors.New(fmt.Sprintf("%s with this idenitifcation is not found", rt.Mapper.V3ObjectName()))
		}
	} else {
		return []*schema.ResourceData{}, errwrap.Wrapf(fmt.Sprintf("query V3 object %s encountered error", rt.Mapper.V3ObjectName()), err)
	}

}

func (rt *ResourceTemplate) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Client := m.(v3client.Client)

	if v3Id, diags := rt.Mapper.V3Identity(d); len(diags) > 0 {
		return diags
	} else if rv, err := rt.DoRead(ctx, v3Client, v3Id); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("invalid query for %s", rt.Mapper.V3ObjectName()),
			Detail:   fmt.Sprintf("qeurying %s is not possible: %s", rt.Mapper.V3ObjectName(), err.Error()),
		}}
	} else {
		if rv != nil {
			return rt.Mapper.SetState(rv, d)
		} else {
			d.SetId("")
			return nil
		}
	}
}

func (rt *ResourceTemplate) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Client := m.(v3client.Client)

	if v3Obj, v3Ctx, diags := rt.Mapper.Upsertable(d); len(diags) > 0 {
		return diags
	} else if rv, err := rt.DoCreate(ctx, v3Client, v3Obj, v3Ctx); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("create of %s failed with an error", rt.Mapper.V3ObjectName()),
			Detail:   fmt.Sprintf("an error was encounterd during %s create: %s", rt.Mapper.V3ObjectName(), err.Error()),
		}}
	} else {
		if rv != nil {
			return rt.Mapper.SetState(rv, d)
		} else {
			d.SetId("")
			return nil
		}
	}
}

func (rt *ResourceTemplate) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Client := m.(v3client.Client)

	if v3Id, diags := rt.Mapper.V3Identity(d); len(diags) > 0 {
		return diags
	} else if rv, err := rt.DoUpdate(ctx, v3Client, v3Id); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("update query for %s has returned an error", rt.Mapper.V3ObjectName()),
			Detail:   fmt.Sprintf("an error was encountered during  %s update: %s", rt.Mapper.V3ObjectName(), err.Error()),
		}}
	} else {
		if rv != nil {
			return rt.Mapper.SetState(rv, d)
		} else {
			d.SetId("")
			return nil
		}
	}
}

func (rt *ResourceTemplate) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Client := m.(v3client.Client)

	if v3Id, diags := rt.Mapper.V3Identity(d); len(diags) > 0 {
		return diags
	} else {

		// If the template specifies the projection against deleting resources that are related to other
		// resource (e.g. no mapped in Terraform), the count is executed before deletion.
		if rt.DoCountOffending != nil {
			if cnt, offendingCountErr := rt.DoCountOffending(ctx, v3Client, v3Id); offendingCountErr != nil {
				return diag.Diagnostics{diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("query for offending objects of %s has returned an error", rt.Mapper.V3ObjectName()),
					Detail:   fmt.Sprintf("an error was encountered during query for object preventing the deletion of %s: %s", rt.Mapper.V3ObjectName(), offendingCountErr.Error()),
				}}
			} else if cnt > 0 {
				return diag.Diagnostics{diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "offending objets prevent deletion",
					Detail:   fmt.Sprintf("there are %d other objects preventing deletion of this %s", cnt, rt.Mapper.V3ObjectName()),
				}}
			}
		}

		if err := rt.DoDelete(ctx, v3Client, v3Id); err != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("delete of %s has returned an error", rt.Mapper.V3ObjectName()),
				Detail:   fmt.Sprintf("an error was encountered during  %s delete: %s", rt.Mapper.V3ObjectName(), err.Error()),
			}}
		} else {
			d.SetId("")
			return nil
		}
	}
}
