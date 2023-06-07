package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestServiceErrorSetBuilderWillProduceSchema(t *testing.T) {
	sch := ServiceErrorSetResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(sch) > 0)
}

func TestServiceErrorSetIdentityMapping(t *testing.T) {
	autoTestIdentity(t, ServiceErrorSetResourceSchemaBuilder, masherytypes.ErrorSetIdentifier{
		ErrorSetId: "defg",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "abc",
		},
	})

	autoTestParentIdentity(t, ServiceErrorSetResourceSchemaBuilder, masherytypes.ServiceIdentifier{
		ServiceId: "defg",
	})
}

func TestServiceErrorSetBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceErrorSetResourceSchemaBuilder, func() masherytypes.ErrorSet {
		return masherytypes.ErrorSet{}
	})
}

func TestServiceErrorSetBuilderMappingsSaves(t *testing.T) {
	mapper, state := ServiceErrorSetResourceSchemaBuilder.MapperAndTestData()

	errSet := masherytypes.ErrorSet{
		ErrorMessages: &[]masherytypes.MasheryErrorMessage{
			{
				Code:         403,
				DetailHeader: "Account Inactive",
				ResponseBody: "",
				Status:       "Forbidden",
				Id:           "ERR_403_DEVELOPER_INACTIVE",
			},
			{
				Code:         414,
				DetailHeader: "detail",
				ResponseBody: "response body",
				Status:       "Request-URI Too Long",
				Id:           "ERR_414_REQUEST_URI_TOO_LONG",
			},
		},
	}

	mapper.RemoteToSchema(&errSet, state)
	v := state.Get(mashschema.MashSvcErrorSetMessage).(*schema.Set)
	assert.Equal(t, 1, v.Len())
}

func TestServiceErrorSetBuilderValidatesDefaults(t *testing.T) {
	mapper, state := ServiceErrorSetResourceSchemaBuilder.MapperAndTestData()

	obj1 := ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         403,
		DetailHeader: "Account Inactive",
		ResponseBody: "",
		Status:       "Forbidden",
		Id:           "ERR_403_DEVELOPER_INACTIVE",
	})
	obj2 := ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         414,
		DetailHeader: "detail",
		ResponseBody: "response body",
		Status:       "Request-URI Too Long",
		Id:           "ERR_414_REQUEST_URI_TOO_LONG",
	})

	err := state.Set(mashschema.MashSvcErrorSetMessage, []interface{}{obj1, obj2})
	assert.Nil(t, err)

	dg := mapper.IsStateValid(state)
	assert.True(t, len(dg) > 0)
	assert.Equal(t, "invalid input for field error_message", dg[0].Summary)
	assert.Equal(t, "field validation has returned the following error: error message for ERR_403_DEVELOPER_INACTIVE matches Mashery-default", dg[0].Detail)
}

func TestServiceErrorSetBuilderValidatesInvalidIds(t *testing.T) {
	mapper, state := ServiceErrorSetResourceSchemaBuilder.MapperAndTestData()

	obj1 := ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         403,
		DetailHeader: "Account Inactive",
		ResponseBody: "",
		Status:       "Forbidden",
		Id:           "ERR_4034_DEVELOPER_INACTIVE",
	})
	obj2 := ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         414,
		DetailHeader: "detail",
		ResponseBody: "response body",
		Status:       "Request-URI Too Long",
		Id:           "ERR_414_REQUEST_URI_TOO_LONG",
	})

	err := state.Set(mashschema.MashSvcErrorSetMessage, []interface{}{obj1, obj2})
	assert.Nil(t, err)

	dg := mapper.IsStateValid(state)
	assert.True(t, len(dg) > 0)
	assert.Equal(t, "invalid input for field error_message", dg[0].Summary)
	assert.Equal(t, "field validation has returned the following error: unsupported error code: ERR_4034_DEVELOPER_INACTIVE", dg[0].Detail)
}

func TestServiceErrorSetBuilderValidatesCustomMessages(t *testing.T) {
	mapper, state := ServiceErrorSetResourceSchemaBuilder.MapperAndTestData()

	obj2 := ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         414,
		DetailHeader: "detail",
		ResponseBody: "response body",
		Status:       "Request-URI Too Long",
		Id:           "ERR_414_REQUEST_URI_TOO_LONG",
	})

	err := state.Set(mashschema.MashSvcErrorSetMessage, []interface{}{obj2})
	assert.Nil(t, err)

	dg := mapper.IsStateValid(state)
	assert.True(t, len(dg) == 0)
}
