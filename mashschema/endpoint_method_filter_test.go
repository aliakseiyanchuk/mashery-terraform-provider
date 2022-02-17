package mashschema_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestJsonPattern(t *testing.T) {
	if b := mashschema.JsonRegexp.MatchString("/a,/b"); !b {
		t.Error("Not matching!")
	}
}

func TestV3EndpointMethodFilterToResourceData(t *testing.T) {
	now := masherytypes.MasheryJSONTime(time.Now())

	d := mashschema.ServiceEndpointMethodFilterMapper.NewResourceData()
	//mashery.assertOk(t, d.Set(mashschema.MashObjName, "DefaultFilter"))
	//mashery.assertOk(t, d.Set(mashschema.ServiceEndpointMethodRef, "service::endpoint::methodId"))

	ref := masherytypes.MasheryResponseFilter{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "filterId",
			Name:    "DefaultFilter",
			Created: &now,
			Updated: &now,
		},
		Notes:            "notes",
		XmlFilterFields:  "xml,fields,a",
		JsonFilterFields: "json,fields,b",
	}

	setChecks := mashschema.ServiceEndpointMethodFilterMapper.PersistTyped(context.TODO(), &ref, d)
	LogErrorDiagnostics(t, "setting full data", &setChecks)

	//mashery.assertResourceContainsKey(t, d, MashObjCreated)
	//mashery.assertResourceContainsKey(t, d, MashObjUpdated)
	//mashery.assertResourceHasStringKey(t, d, mashschema.MashEndpointMethodFilterId, "filterId")

	//mashFilterIdent := mashschema.ServiceEndpointMethodFilterIdentifier{
	//	ServiceEndpointMethodIdentifier: mashschema.ServiceEndpointMethodIdentifier{
	//		ServiceEndpointIdentifier: mashschema.ServiceEndpointIdentifier{
	//			ServiceIdentifier: mashschema.ServiceIdentifier{
	//				ServiceId: "serviceId",
	//			},
	//			EndpointId: "endpointId",
	//		},
	//		MethodId: "methodId",
	//	},
	//	FilterId: "filterId",
	//}
	//d.SetId(mashFilterIdent.Id())

	//refOutput, dg := mashschema.MashEndpointMethodFilterUpsertable(d)

	//LogErrorDiagnostics(t, "Upsert", &dg)

	//mashery.assertSameString(t, "Id", &ref.Id, &refOutput.Id)
	//mashery.assertSameString(t, "Name", &ref.Name, &refOutput.Name)
	//mashery.assertSameString(t, "Notes", &ref.Notes, &refOutput.Notes)
	//mashery.assertSameString(t, "XmlFilterFields", &ref.XmlFilterFields, &refOutput.XmlFilterFields)
	//mashery.assertSameString(t, "JsonFilterFields", &ref.JsonFilterFields, &refOutput.JsonFilterFields)
}
