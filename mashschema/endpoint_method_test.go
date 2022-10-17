package mashschema_test

import (
	"testing"
)

func TestV3MasheryEndpointMethodToResourceState(t *testing.T) {
	//now := masherytypes.MasheryJSONTime(time.Now())
	//d := TestResourceData(&mashschema.EndpointMethodSchema)
	////mashery.assertOk(t, d.Set(MashObjName, "name"))
	//
	//refInput := masherytypes.MasheryMethod{
	//	AddressableV3Object: masherytypes.AddressableV3Object{
	//		Id:      "methodId",
	//		Name:    "name",
	//		Created: &now,
	//		Updated: &now,
	//	},
	//	SampleJsonResponse: "a-json",
	//	SampleXmlResponse:  "an-xml",
	//}
	//
	//setChecks := mashschema.V3EndpointMethodToResourceState(&refInput, d)
	//LogErrorDiagnostics(t, "setting full data", &setChecks)

	//mashery.assertResourceContainsKey(t, d, MashObjCreated)
	//mashery.assertResourceContainsKey(t, d, MashObjUpdated)
	//mashery.assertResourceHasStringKey(t, d, mashschema.ServiceEndpointMethodRef, "methodId")
	//
	//if d.HasChange(MashObjName) {
	//	t.Errorf("Object name should not have been altered while being set")
	//}

	//methIdent := mashschema.ServiceEndpointMethodIdentifier{
	//	ServiceEndpointIdentifier: mashschema.ServiceEndpointIdentifier{
	//		ServiceId:  "serviceId",
	//		EndpointId: "endpointId",
	//	},
	//	MethodId: "methodId",
	//}
	//
	//d.SetId(methIdent.Id())

	//refOutput := mashschema.MashEndpointMethodUpsertable(d)

	//mashery.assertSameString(t, "Id", &refInput.Id, &refOutput.Id)
	//mashery.assertSameString(t, "Name", &refInput.Name, &refOutput.Name)
	//mashery.assertSameString(t, "SampleJsonResponse", &refInput.SampleJsonResponse, &refOutput.SampleJsonResponse)
	//mashery.assertSameString(t, "SampleXmlResponse", &refInput.SampleXmlResponse, &refOutput.SampleXmlResponse)

}
