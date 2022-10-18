package mashschema

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

// TFResourceSchema Shortcut for map of schema pointers to make the code a bit more descriptive
type TFResourceSchema = map[string]*schema.Schema
type UpsertFunc func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics)
type PersistFunc func(rv interface{}, d *schema.ResourceData) diag.Diagnostics
type PersistManyFunc func(rv []interface{}, d *schema.ResourceData) diag.Diagnostics
type V3IdentifierFunc func(d *schema.ResourceData) (interface{}, diag.Diagnostics)

type V3ObjectIdentifier interface{}
type Upsertable interface{}

type DataSourceMapper interface {
	// TerraformSchema get the Terraform schema for this resource
	TerraformSchema() TFResourceSchema

	SetState(rv []interface{}, d *schema.ResourceData) diag.Diagnostics
	SetStateOf(rv interface{}, d *schema.ResourceData) diag.Diagnostics
}

type DataSourceMapperImpl struct {
	DataSourceMapper
	schema      TFResourceSchema
	persistOne  PersistFunc
	persistMany PersistManyFunc
}

func (dsmi *DataSourceMapperImpl) SetState(rv []interface{}, d *schema.ResourceData) diag.Diagnostics {
	if dsmi.persistMany != nil {
		return dsmi.persistMany(rv, d)
	} else {
		return diag.Diagnostics{
			diag.Diagnostic{
				Summary: "persistMany function not defined",
			},
		}
	}
}

func (dsmi *DataSourceMapperImpl) SetStateOf(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
	if dsmi.persistOne != nil {
		return dsmi.persistOne(rv, d)
	} else {
		return diag.Diagnostics{
			diag.Diagnostic{
				Summary: "persistOne function not defined",
			},
		}
	}
}

// TestResourceData create a test resource data for flattening and recoverign.
func (dsmi *DataSourceMapperImpl) TestResourceData() *schema.ResourceData {
	res := schema.Resource{
		Schema: dsmi.schema,
	}

	return res.TestResourceData()
}

// ResourceMapper Base interface for converting HCL constructs into the Mashery V3 object upsertable calls.
type ResourceMapper interface {

	// TerraformSchema get the Terraform schema for this resource
	TerraformSchema() TFResourceSchema

	V3ObjectName() string

	// V3Identity Creates the v3Identity for this object. Will return non-nil diagnostics
	// if the object is not positively identified
	V3Identity(d *schema.ResourceData) (V3ObjectIdentifier, diag.Diagnostics)

	// Upsertable creates an update-insertable object that can be passed onto the client. Where applicable,
	// the {@link V3ObjectIdentifier} will contain the identifier of the context
	Upsertable(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics)

	// SetState Set the state the V3 object in the terraform schema
	SetState(rv Upsertable, d *schema.ResourceData) diag.Diagnostics

	TestResourceData() *schema.ResourceData

	TestResourceDataWith(init map[string]interface{}) (*schema.ResourceData, diag.Diagnostics)
}

type ResourceMapperImpl struct {
	ResourceMapper

	v3ObjectName string
	schema       TFResourceSchema
	upsertFunc   UpsertFunc
	persistFunc  PersistFunc
	v3Identity   V3IdentifierFunc
}

func (rmi *ResourceMapperImpl) persistMap(inp interface{}, fields map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {
	d.SetId(CompoundId(inp))
	return SetResourceFields(fields, d)
}

func (rmi *ResourceMapperImpl) lackingIdentificationDiagnostic(fields ...string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Lacking identification",
		Detail:   fmt.Sprintf("field(s) %s must be set to identify V3 %s object and must match object schema", strings.Join(fields, ", "), rmi.v3ObjectName),
	}
}

func (rmi *ResourceMapperImpl) V3ObjectName() string {
	return rmi.v3ObjectName
}

func (mi *ResourceMapperImpl) TestResourceData() *schema.ResourceData {
	res := schema.Resource{
		Schema: mi.schema,
	}

	return res.TestResourceData()
}

func (mi *ResourceMapperImpl) TestResourceDataWith(init map[string]interface{}) (*schema.ResourceData, diag.Diagnostics) {
	res := schema.Resource{
		Schema: mi.schema,
	}

	rv := res.TestResourceData()
	rvd := SetResourceFields(init, rv)

	return rv, rvd
}

func (m *ResourceMapperImpl) TerraformSchema() TFResourceSchema {
	return m.schema
}

func (m *ResourceMapperImpl) V3Identity(d *schema.ResourceData) (V3ObjectIdentifier, diag.Diagnostics) {
	if m.v3Identity != nil {
		return m.v3Identity(d)
	} else {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Summary: "upsert function was not defined",
			},
		}
	}
}

func (m *ResourceMapperImpl) Upsertable(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
	if m.upsertFunc != nil {
		return m.upsertFunc(d)
	} else {
		return nil, nil, diag.Diagnostics{
			diag.Diagnostic{
				Summary: "upsert function was not defined",
			},
		}
	}
}

func (m *ResourceMapperImpl) SetState(rv Upsertable, d *schema.ResourceData) diag.Diagnostics {
	if m.persistFunc != nil {
		return m.persistFunc(rv, d)
	} else {
		return diag.Diagnostics{
			diag.Diagnostic{
				Summary: "persist function not defined",
			},
		}
	}
}

// SetResourceFields Set resource data fields, recording any errors occurring while being set.
func SetResourceFields(data map[string]interface{}, res *schema.ResourceData) diag.Diagnostics {
	rv := diag.Diagnostics{}

	for k, v := range data {
		if err := res.Set(k, v); err != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("failed to set field %s", k),
				Detail:   fmt.Sprintf("settings field %s encoutnered an error: %s", k, err.Error()),
			})
		}
	}

	return rv
}
