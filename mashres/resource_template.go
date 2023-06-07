package mashres

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/tfmapper"
)

type ReaderFunc[Ident any, MType any] func(context.Context, v3client.Client, Ident) (*MType, error)
type CreatorFunc[ParentIdent any, Ident any, MType any] func(context.Context, v3client.Client, ParentIdent, MType) (*MType, *Ident, error)
type UpdaterFunc[Ident any, MType any] func(context.Context, v3client.Client, Ident, MType) (*MType, error)
type DeleterFunc[Ident any] func(context.Context, v3client.Client, Ident) error
type OffendersCounterFunc[Ident any] func(context.Context, v3client.Client, Ident) (int64, error)

type ResourceTemplate[ParentIdent any, Ident any, MType any] struct {
	Schema map[string]*schema.Schema
	Mapper *tfmapper.Mapper[ParentIdent, Ident, MType]

	UpsertableFunc func() MType

	// Validate function: if defined, checks the consistency of the parameters with each other.
	ValidateFunc func(parent ParentIdent, upsertable MType) string

	DoRead   ReaderFunc[Ident, MType]
	DoCreate CreatorFunc[ParentIdent, Ident, MType]
	DoUpdate UpdaterFunc[Ident, MType]
	DoDelete DeleterFunc[Ident]

	DoCountOffending OffendersCounterFunc[Ident]
}

// ResourceSchema returns the Terraform data source schema
func (rt *ResourceTemplate[ParentIdent, Ident, MType]) ResourceSchema() *schema.Resource {

	// Update function will be enabled only in case the schema contains non-computed elements as well as
	//
	updateCtx := rt.Update
	if rt.IsUpdateSuperfluous() {
		updateCtx = nil
	}

	return &schema.Resource{
		ReadContext:   rt.Read,
		CreateContext: rt.Create,
		UpdateContext: updateCtx,
		DeleteContext: rt.Delete,
		Schema:        rt.Schema,
	}
}

func (rt *ResourceTemplate[ParentIdent, Ident, MType]) IsUpdateSuperfluous() bool {
	for _, v := range rt.Schema {
		if !v.ForceNew && !v.Computed {
			return false
		}
	}

	return true
}

func (rt *ResourceTemplate[ParentIdent, Ident, MType]) TestState() *schema.ResourceData {
	res := schema.Resource{
		Schema: rt.Schema,
	}

	return res.TestResourceData()
}

func (rt *ResourceTemplate[ParentIdent, Ident, MType]) Read(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	ident, err := rt.Mapper.Identity(state)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("invalid resource identifier"),
			Detail:   fmt.Sprintf("attempt to parse identified returned this error: %s", err.Error()),
		}}
	}

	v3Client := m.(v3client.Client)
	curState, err := rt.DoRead(ctx, v3Client, ident)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unexpected error returned from Mashery V3 api"),
			Detail:   fmt.Sprintf("read operation on Mashery V3 api has returned the following error: %s", err.Error()),
		}}
	}

	rv := diag.Diagnostics{}

	if curState != nil {
		if setErrors := rt.Mapper.RemoteToSchema(curState, state); len(setErrors) > 0 {
			rv = append(rv, setErrors...)
		}
	} else {
		rt.Mapper.ResetIdentity(state)
	}

	return rv
}

func (rt *ResourceTemplate[ParentIdent, Ident, MType]) Create(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	parentIdent, err := rt.Mapper.ParentIdentity(state)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid reference to parent object",
			Detail:   fmt.Sprintf("attempt to parse parent identity has returned the following error: %s", err.Error()),
		}}
	}

	if validationDiags := rt.Mapper.IsStateValid(state); len(validationDiags) > 0 {
		return validationDiags
	}

	rv := diag.Diagnostics{}

	upsertable := rt.UpsertableFunc()
	rt.Mapper.SchemaToRemote(state, &upsertable)

	if rt.ValidateFunc != nil {
		if msg := rt.ValidateFunc(parentIdent, upsertable); len(msg) > 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "object cannot be created due to conflict in the parameters",
				Detail:   msg,
			}}
		}
	}

	v3Client := m.(v3client.Client)
	readBack, ident, err := rt.DoCreate(ctx, v3Client, parentIdent, upsertable)

	if ident != nil {
		if idRV := rt.Mapper.AssignIdentity(*ident, state); idRV != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Detail:   "was unable to assign the identity to this resource",
				Summary:  fmt.Sprintf("assinging identity has returned error: %s", idRV.Error()),
			})
		}
	}

	if setRV := rt.Mapper.RemoteToSchema(readBack, state); len(setRV) > 0 {
		rv = append(rv, setRV...)
	}

	if err != nil {
		rv = append(rv, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unexpected error returned from Mashery V3 api during creating object"),
			Detail:   fmt.Sprintf("create operation on Mashery V3 api has returned the following error: %s", err.Error()),
		})
	} else if ident == nil {
		rv = append(rv, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("create function returned a nil identity"),
			Detail:   fmt.Sprintf("create functions must return a non-nil identifier if create operation has been successful"),
		})
	} else if readBack == nil {
		rv = append(rv, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("create function returned a nil object"),
			Detail:   fmt.Sprintf("create functions must return a non-nil object if create operation has been successful"),
		})
	}

	return rv
}

func (rt *ResourceTemplate[ParentIdent, Ident, MType]) Update(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	ident, err := rt.Mapper.Identity(state)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid identity object",
			Detail:   fmt.Sprintf("attempt to parse identity has returned the following error: %s", err.Error()),
		}}
	}

	if validationDiags := rt.Mapper.IsStateValid(state); len(validationDiags) > 0 {
		return validationDiags
	}

	upsertable := rt.UpsertableFunc()
	rt.Mapper.SchemaToRemote(state, &upsertable)

	v3Client := m.(v3client.Client)
	readBack, err := rt.DoUpdate(ctx, v3Client, ident, upsertable)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unexpected error returned from Mashery V3 api during update object"),
			Detail:   fmt.Sprintf("update operation on Mashery V3 api has returned the following error: %s", err.Error()),
		}}
	}

	rv := diag.Diagnostics{}
	if setRV := rt.Mapper.RemoteToSchema(readBack, state); len(setRV) > 0 {
		rv = append(rv, setRV...)
	}

	return rv
}

func (rt *ResourceTemplate[ParentIdent, Ident, MType]) Delete(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	ident, err := rt.Mapper.Identity(state)
	if err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid identity object",
			Detail:   fmt.Sprintf("attempt to parse identity has returned the following error: %s", err.Error()),
		}}
	}

	v3Client := m.(v3client.Client)
	if rt.DoCountOffending != nil {
		if cnt, offendingCountErr := rt.DoCountOffending(ctx, v3Client, ident); offendingCountErr != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("query for offending objects of has returned an error"),
				Detail:   fmt.Sprintf("querying offeding object has returned error: %s", offendingCountErr.Error()),
			}}
		} else if cnt > 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "offending objets prevent deletion",
				Detail:   fmt.Sprintf("there are %d other objects preventing deletion", cnt),
			}}
		}
	}

	if err := rt.DoDelete(ctx, v3Client, ident); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "delete has returned an error",
			Detail:   fmt.Sprintf("attempt to delete has encountered error: %s", err.Error()),
		}}
	} else {
		rt.Mapper.ResetIdentity(state)
		return nil
	}
}
