package mashschema_test

import (
	"testing"
)

func TestV3ErrorSetConversion(t *testing.T) {
	//now := masherytypes.MasheryJSONTime(time.Now())
	//d := NewResourceData(&mashschema.ServiceErrorSetSchema)
	//
	//v3 := masherytypes.MasheryErrorSet{
	//	AddressableV3Object: masherytypes.AddressableV3Object{
	//		Id:      "id",
	//		Name:    "name",
	//		Created: &now,
	//		Updated: &now,
	//	},
	//	Type:      "type",
	//	JSONP:     true,
	//	JSONPType: "jsonpType",
	//}
	//
	//diags := mashschema.V3ErrorSetToResourceData(&v3, d)
	//if len(diags) > 0 {
	//	t.Errorf("full conversion has encountered %d errors where none were expected", len(diags))
	//}

	//reverse := mashschema.V3ErrorSetUpsertable(d)
	//if len(diags) > 0 {
	//	t.Errorf("Reverse conversion has encountered %d errors where none were expected", len(diags))
	//}

	//mashery.assertSameString(t, `Name`, &v3.Name, &reverse.Name)
	//mashery.assertSameString(t, `Type`, &v3.Type, &reverse.Type)
	//mashery.assertSameBool(t, `JSONP`, &v3.JSONP, &reverse.JSONP)
	//mashery.assertSameString(t, `JSONPType`, &v3.JSONPType, &reverse.JSONPType)
}

func TestV3ErrorMessageExtraction(t *testing.T) {
	//d := NewResourceData(&mashschema.ServiceErrorSetSchema)
	//d.SetId("serviceId::errorSetId")

	//passedSet := map[string]interface{}{
	//	mashschema.MashSvcErrorSetMessage: []map[string]interface{}{
	//		mashschema.V3ErrorMessageForResourceData(masherytypes.MasheryErrorMessage{
	//			Id:           "ERR_403_anything",
	//			Code:         999,
	//			Status:       "status",
	//			DetailHeader: "detailHeader",
	//			ResponseBody: "responseBody",
	//		}),
	//	},
	//}
	//if diags := SetResourceFields(passedSet, d); len(diags) > 0 {
	//	t.Errorf("could not set source message")
	//}
	//messages := mashschema.V3ErrorSetMessages(d)
	//assert.Equal(t, 1, len(messages), "Should have one message")
	//
	//assert.Equal(t, "ERR_403_anything", messages[0].Id, "Id should be ERR_403_aaa")
	//assert.Equal(t, "status", messages[0].Status, "Status should be status")
	//assert.Equal(t, "detailHeader", messages[0].DetailHeader, "DetailHeader should be detailHeader")
	//assert.Equal(t, "responseBody", messages[0].ResponseBody, "ResponseBody should be responseBody")
	//assert.Equal(t, 403, messages[0].Code, "Code should be 403")
}
