package mashschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DummyResourceMapperType struct {
	DataSourceMapperImpl
}

type DummyObject struct {
	Id      string
	Payload string
}

var DummyDataSourceMapper *DummyResourceMapperType

func init() {
	DummyDataSourceMapper = &DummyResourceMapperType{
		DataSourceMapperImpl: DataSourceMapperImpl{
			v3ObjectName: "dummy object",
			persistOne: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				ptr := rv.(*DummyObject)
				d.SetId(ptr.Id)

				data := map[string]interface{}{
					"a": ptr.Payload,
				}
				return SetResourceFields(data, d)
			},
			persistMany: nil,
		},
	}

	DummyDataSourceMapper.SchemaBuilder().
		AddComputedString("a", "single dummy string")
}
