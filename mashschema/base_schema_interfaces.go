package mashschema

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TFResourceSchema Shortcut for map of schema pointers to make the code a bit more descriptive
type TFResourceSchema = map[string]*schema.Schema
type UpsertFunc func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics)
type PersistFunc func(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics
type IdentifierFunc func() interface{}

// Mapper Base interface for converting HCL constructs into the Mashery V3 object upsertable calls.
type Mapper interface {

	// TerraformSchema get the Terraform schema for this resource
	TerraformSchema() TFResourceSchema

	Upsertable(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics)

	SetState(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics

	// CreateIdentifier create the identifier object that will encode identifier
	CreateIdentifier() interface{}
}

type MapperImpl struct {
	Mapper

	schema      TFResourceSchema
	upsertFunc  UpsertFunc
	persistFunc PersistFunc
	identifier  IdentifierFunc
}

func (mi *MapperImpl) GetQuery(d *schema.ResourceData) map[string]string {
	return extractStringMap(d, MashDataSourceSearch)
}

func (mi *MapperImpl) IsMatchRequired(d *schema.ResourceData) bool {
	return extractBool(d, MashDataSourceRequired, true)
}

func (mi *MapperImpl) CreateIdentifier() interface{} {
	if mi.identifier != nil {
		return mi.identifier()
	} else {
		return nil
	}
}

func (mi *MapperImpl) NewResourceData() *schema.ResourceData {
	res := schema.Resource{
		Schema: mi.schema,
	}

	return res.TestResourceData()
}

func (m *MapperImpl) TerraformSchema() TFResourceSchema {
	return m.schema
}

func (m *MapperImpl) Upsertable(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, d)
	} else {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Summary: "upsert function not defined",
			},
		}
	}
}

func (m *MapperImpl) SetState(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics {
	if m.persistFunc != nil {
		return m.persistFunc(ctx, rv, d)
	} else {
		return diag.Diagnostics{
			diag.Diagnostic{
				Summary: "persist function not defined",
			},
		}
	}
}

// SetResourceFields Set resource data fields, recording any errors occurring while being set.
func (m *MapperImpl) SetResourceFields(ctx context.Context, data map[string]interface{}, res *schema.ResourceData) diag.Diagnostics {
	rv := diag.Diagnostics{}

	for k, v := range data {
		if err := res.Set(k, v); err != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("failed to set field %s", k),
				Detail:   err.Error(),
			})
		}
	}

	return rv
}
