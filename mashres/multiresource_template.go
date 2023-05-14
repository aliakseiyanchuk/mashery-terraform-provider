package mashres

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/funcsupport"
	"terraform-provider-mashery/tfmapper"
)

// MultiResourceTemplate for situations where a single Terraform script needs to be translated to multiple read operations
type MultiResourceTemplate[ParentIdent any, Ident any] struct {
	Schema    map[string]*schema.Schema
	Templates []*ResourceTemplate[ParentIdent, Ident, any]
}

func BuildMultiResourceTemplate[ParentIdent any, Ident any](template ...*ResourceTemplate[ParentIdent, Ident, any]) *MultiResourceTemplate[ParentIdent, Ident] {
	var schemas = make([]map[string]*schema.Schema, len(template))

	for i, t := range template {
		schemas[i] = t.Schema
	}

	rv := MultiResourceTemplate[ParentIdent, Ident]{
		Schema:    tfmapper.MergeSchemas(schemas...),
		Templates: template,
	}

	return &rv
}

type ContextFunc func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
type ContextLocator[ParentIdent any, Ident any] funcsupport.Function[*ResourceTemplate[ParentIdent, Ident, any], ContextFunc]

func (mt *MultiResourceTemplate[ParentIdent, Ident]) WithEachTemplate(ctx context.Context, state *schema.ResourceData, m interface{}, fLoc ContextLocator[ParentIdent, Ident]) diag.Diagnostics {
	rv := diag.Diagnostics{}

	for _, t := range mt.Templates {
		f := fLoc(t)
		if f != nil {
			dg := f(ctx, state, m)
			if len(dg) > 0 {
				rv = append(dg)
			}
			if dg.HasError() {
				return rv
			}
		}
	}

	return rv
}

func (mt *MultiResourceTemplate[ParentIdent, Ident]) Read(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	return mt.WithEachTemplate(ctx, state, m, func(in *ResourceTemplate[ParentIdent, Ident, any]) ContextFunc {
		return in.Read
	})
}

func (mt *MultiResourceTemplate[ParentIdent, Ident]) Create(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	return mt.WithEachTemplate(ctx, state, m, func(in *ResourceTemplate[ParentIdent, Ident, any]) ContextFunc {
		return in.Create
	})
}

func (mt *MultiResourceTemplate[ParentIdent, Ident]) Update(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	return mt.WithEachTemplate(ctx, state, m, func(in *ResourceTemplate[ParentIdent, Ident, any]) ContextFunc {
		return in.Update
	})
}

func (mt *MultiResourceTemplate[ParentIdent, Ident]) Delete(ctx context.Context, state *schema.ResourceData, m interface{}) diag.Diagnostics {
	return mt.WithEachTemplate(ctx, state, m, func(in *ResourceTemplate[ParentIdent, Ident, any]) ContextFunc {
		return in.Delete
	})
}
