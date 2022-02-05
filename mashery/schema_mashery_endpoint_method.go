package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashServiceEndpointMethodSampleJson = "sample_json"
	MashServiceEndpointMethodSampleXml  = "sample_xml"
	MashServiceEndpointMethodId         = "method_id" // Shared, to be moved
	MashServiceEndpointMethodRef        = "method_id" // Shared, to be moved
)

var EndpointMethodSchema = map[string]*schema.Schema{
	MashEndpointId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Endpoint to which the method must be attached",
	},
	MashObjName: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Method name, as it would be detected by Mashery",
	},
}

type ServiceEndpointMethodIdentifier struct {
	ServiceEndpointIdentifier
	MethodId string
}

func (emi *ServiceEndpointMethodIdentifier) Id() string {
	return CreateCompoundId(emi.ServiceId, emi.EndpointId, emi.MethodId)
}

func (emi *ServiceEndpointMethodIdentifier) From(id string) {
	ParseCompoundId(id, &emi.ServiceId, &emi.EndpointId, &emi.MethodId)
}

func (emi *ServiceEndpointMethodIdentifier) IsIdentified() bool {
	return emi.ServiceEndpointIdentifier.IsIdentified() && len(emi.MethodId) > 0
}

func MashEndpointMethodUpsertable(d *schema.ResourceData) masherytypes.MasheryMethod {
	ident := ServiceEndpointMethodIdentifier{}
	ident.From(d.Id())

	return masherytypes.MasheryMethod{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   ident.MethodId,
			Name: extractString(d, MashObjName, ""),
		},
		SampleJsonResponse: extractString(d, MashServiceEndpointMethodSampleJson, ""),
		SampleXmlResponse:  extractString(d, MashServiceEndpointMethodSampleXml, ""),
	}
}

func V3EndpointMethodToResourceState(inp *masherytypes.MasheryMethod, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashServiceEndpointMethodId:         inp.Id,
		MashObjCreated:                      inp.Created.ToString(),
		MashObjUpdated:                      inp.Updated.ToString(),
		MashServiceEndpointMethodSampleJson: nullForEmptyString(inp.SampleJsonResponse),
		MashServiceEndpointMethodSampleXml:  nullForEmptyString(inp.SampleXmlResponse),
	}

	return SetResourceFields(data, d)
}

func initEndpointMethodSchemaBoilerplate() {
	addComputedString(&EndpointMethodSchema, MashServiceEndpointMethodRef, "V3 identifier of this method")
	addComputedString(&EndpointMethodSchema, MashObjCreated, "Date/time the object was created")
	addComputedString(&EndpointMethodSchema, MashObjUpdated, "Date/time the object was updated")

	addOptionalString(&EndpointMethodSchema, MashServiceEndpointMethodSampleJson, "Sample JSON response")
	addOptionalString(&EndpointMethodSchema, MashServiceEndpointMethodSampleXml, "Sample XML response")
}

func init() {
	initEndpointMethodSchemaBoilerplate()
}
