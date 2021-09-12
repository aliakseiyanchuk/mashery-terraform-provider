package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestJsonPattern(t *testing.T) {
	if b := mashery.JsonRegexp.MatchString("/a,/b"); !b {
		t.Error("Not matching!")
	}
}

func TestV3EndpointMethodFilterToResourceData(t *testing.T) {
	now := v3client.MasheryJSONTime(time.Now())

	d := NewResourceData(&mashery.EndpointMethodFilterSchema)
	assertOk(t, d.Set(mashery.MashObjName, "DefaultFilter"))
	assertOk(t, d.Set(mashery.MashServiceEndpointMethodRef, "service::endpoint::methodId"))

	ref := v3client.MasheryResponseFilter{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "filterId",
			Name:    "DefaultFilter",
			Created: &now,
			Updated: &now,
		},
		Notes:            "notes",
		XmlFilterFields:  "xml,fields,a",
		JsonFilterFields: "json,fields,b",
	}

	setChecks := mashery.V3EndpointMethodFilterToResourceData(&ref, d)
	LogErrorDiagnostics(t, "setting full data", &setChecks)

	assertResourceContainsKey(t, d, mashery.MashObjCreated)
	assertResourceContainsKey(t, d, mashery.MashObjUpdated)
	assertResourceHasStringKey(t, d, mashery.MashEndpointMethodFilterId, "filterId")

	mashFilterIdent := mashery.ServiceEndpointMethodFilterIdentifier{
		ServiceEndpointMethodIdentifier: mashery.ServiceEndpointMethodIdentifier{
			ServiceEndpointIdentifier: mashery.ServiceEndpointIdentifier{
				ServiceId:  "serviceId",
				EndpointId: "endpointId",
			},
			MethodId: "methodId",
		},
		FilterId: "filterId",
	}
	d.SetId(mashFilterIdent.Id())

	refOutput, dg := mashery.MashEndpointMethodFilterUpsertable(d)

	LogErrorDiagnostics(t, "Upsert", &dg)

	assertSameString(t, "Id", &ref.Id, &refOutput.Id)
	assertSameString(t, "Name", &ref.Name, &refOutput.Name)
	assertSameString(t, "Notes", &ref.Notes, &refOutput.Notes)
	assertSameString(t, "XmlFilterFields", &ref.XmlFilterFields, &refOutput.XmlFilterFields)
	assertSameString(t, "JsonFilterFields", &ref.JsonFilterFields, &refOutput.JsonFilterFields)
}
