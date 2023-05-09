package tfmapper

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringFieldMapperExpectedOperation(t *testing.T) {
	builder := NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]()

	fieldMapper := StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.ApiKeyValueLocationKey
		},
		FieldMapperBase: FieldMapperBase{
			Key: "tf_path",
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "some description",
			},
		},
	}

	mapper := builder.Add(&fieldMapper).Mapper()

	other := masherytypes.Endpoint{
		ApiKeyValueLocationKey: "api_key",
	}

	state := builder.TestResourceData()
	dg := mapper.RemoteToSchema(&other, state)
	assert.Equal(t, 0, len(dg))

	reverse := masherytypes.Endpoint{}
	mapper.SchemaToRemote(state, &reverse)

	assert.Equal(t, other.ApiKeyValueLocationKey, reverse.ApiKeyValueLocationKey)
}

func TestStringFieldMapperUsingDefaultValues(t *testing.T) {
	builder := NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]()

	fieldMapper := StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.ApiKeyValueLocationKey
		},
		FieldMapperBase: FieldMapperBase{
			Key: "tf_path",
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Default:     "abc",
				Optional:    true,
				Description: "some description",
			},
		},
	}

	mapper := builder.Add(&fieldMapper).Mapper()

	state := builder.TestResourceData()

	upsert := masherytypes.Endpoint{}
	mapper.SchemaToRemote(state, &upsert)

	assert.Equal(t, "abc", upsert.ApiKeyValueLocationKey)
}

func TestStringFieldMapperWillReturnErrorOnErrorMapping(t *testing.T) {
	builder := NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]()

	fieldMapper := StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.ApiKeyValueLocationKey
		},
		FieldMapperBase: FieldMapperBase{
			Key: "tf_path",
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "some description",
			},
		},
	}

	mapper := builder.Add(&fieldMapper).Mapper()

	other := masherytypes.Endpoint{
		ApiKeyValueLocationKey: "api_key",
	}

	state := builder.TestResourceData()
	dg := mapper.RemoteToSchema(&other, state)
	assert.Equal(t, 1, len(dg))
}
