package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestV3ErrorSetConversion(t *testing.T) {
	now := v3client.MasheryJSONTime(time.Now())
	d := NewResourceData(&mashery.ServiceErrorSetSchema)

	v3 := v3client.MasheryErrorSet{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "id",
			Name:    "name",
			Created: &now,
			Updated: &now,
		},
		Type:      "type",
		JSONP:     true,
		JSONPType: "jsonpType",
	}

	diags := mashery.V3ErrorSetToResourceData(&v3, d)
	if len(diags) > 0 {
		t.Errorf("full conversion has encountered %d errors where none were expected", len(diags))
	}

	reverse := mashery.V3ErrorSetUpsertable(d)
	if len(diags) > 0 {
		t.Errorf("Reverse conversion has encountered %d errors where none were expected", len(diags))
	}

	assertSameString(t, `Name`, &v3.Name, &reverse.Name)
	assertSameString(t, `Type`, &v3.Type, &reverse.Type)
	assertSameBool(t, `JSONP`, &v3.JSONP, &reverse.JSONP)
	assertSameString(t, `JSONPType`, &v3.JSONPType, &reverse.JSONPType)
}

func TestV3ErrorMessageExtraction(t *testing.T) {
	d := NewResourceData(&mashery.ServiceErrorSetSchema)
	d.SetId("serviceId::errorSetId")

	passedSet := map[string]interface{}{
		mashery.MashSvcErrorSetMessage: []map[string]interface{}{
			mashery.V3ErrorMessageForResourceData(v3client.MasheryErrorMessage{
				Id:           "ERR_403_anything",
				Code:         999,
				Status:       "status",
				DetailHeader: "detailHeader",
				ResponseBody: "responseBody",
			}),
		},
	}
	if diags := mashery.SetResourceFields(passedSet, d); len(diags) > 0 {
		t.Errorf("could not set source message")
	}
	messages := mashery.V3ErrorSetMessages(d)
	assert.Equal(t, 1, len(messages), "Should have one message")

	assert.Equal(t, "ERR_403_anything", messages[0].Id, "Id should be ERR_403_aaa")
	assert.Equal(t, "status", messages[0].Status, "Status should be status")
	assert.Equal(t, "detailHeader", messages[0].DetailHeader, "DetailHeader should be detailHeader")
	assert.Equal(t, "responseBody", messages[0].ResponseBody, "ResponseBody should be responseBody")
	assert.Equal(t, 403, messages[0].Code, "Code should be 403")
}
