package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestV3MasheryEndpointMethodToResourceState(t *testing.T) {
	now := v3client.MasheryJSONTime(time.Now())
	d := NewResourceData(&mashery.EndpointMethodSchema)
	assertOk(t, d.Set(mashery.MashObjName, "name"))

	refInput := v3client.MasheryMethod{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "methodId",
			Name:    "name",
			Created: &now,
			Updated: &now,
		},
		SampleJsonResponse: "a-json",
		SampleXmlResponse:  "an-xml",
	}

	setChecks := mashery.V3EndpointMethodToResourceState(&refInput, d)
	LogErrorDiagnostics(t, "setting full data", &setChecks)

	assertResourceContainsKey(t, d, mashery.MashObjCreated)
	assertResourceContainsKey(t, d, mashery.MashObjUpdated)
	assertResourceHasStringKey(t, d, mashery.MashServiceEndpointMethodRef, "methodId")

	if d.HasChange(mashery.MashObjName) {
		t.Errorf("Object name should not have been altered while being set")
	}

	methIdent := mashery.ServiceEndpointMethodIdentifier{
		ServiceEndpointIdentifier: mashery.ServiceEndpointIdentifier{
			ServiceId:  "serviceId",
			EndpointId: "endpointId",
		},
		MethodId: "methodId",
	}

	d.SetId(methIdent.Id())

	refOutput := mashery.MashEndpointMethodUpsertable(d)

	assertSameString(t, "Id", &refInput.Id, &refOutput.Id)
	assertSameString(t, "Name", &refInput.Name, &refOutput.Name)
	assertSameString(t, "SampleJsonResponse", &refInput.SampleJsonResponse, &refOutput.SampleJsonResponse)
	assertSameString(t, "SampleXmlResponse", &refInput.SampleXmlResponse, &refOutput.SampleXmlResponse)

}
