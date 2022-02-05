package mashery

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"regexp"
	"strings"
)

const (
	MashEndpointMethodFilterId           = "filter_id"
	MashServiceEndpointMethodFilterRef   = "service_filter_id"
	MashServiceEndpointMethodFilterNotes = "notes"

	MashServiceEndpointMethodFilterXmlFields  = "xml_fields"
	MashServiceEndpointMethodFilterJsonFields = "json_fields"

	endpointFilterFieldSeparator = " | "
)

var JsonRegexp = regexp.MustCompile("^(?:\\/[^\\/\\| ]+)+(?: \\| (?:\\/[^\\/\\| ]+)+)*$")

var EndpointMethodFilterSchema = map[string]*schema.Schema{
	MashServiceEndpointMethodFilterNotes: {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "Managed by Terraform",
		Description: "Notes added to this method filter",
	},
	MashServiceEndpointMethodRef: {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Endpoint method to which this filter is attached",
		ValidateDiagFunc: validateDiagInputIsEndpointMethodIdentifier,
	},
	MashServiceEndpointMethodFilterXmlFields: {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Fields to strip from the response body of an XML response.",
		Elem:        stringElem(),
	},
	MashServiceEndpointMethodFilterJsonFields: {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Fields to strip from the response body of a JSON response.",
		Elem:        stringElem(),
	},
}

func validateDiagInputIsEndpointMethodIdentifier(i interface{}, path cty.Path) diag.Diagnostics {
	if str, ok := i.(string); ok {
		mid := ServiceEndpointMethodIdentifier{}
		mid.From(str)

		if !mid.IsIdentified() {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Incomplete identifier",
				Detail:        "Endpoint method identifier is incomplete or malformed",
				AttributePath: path,
			}}
		} else {
			return diag.Diagnostics{}
		}
	} else {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Unexpected type",
			Detail:        fmt.Sprintf("Input should be string, but was %s", reflect.TypeOf(i)),
			AttributePath: path,
		}}
	}
}

type ServiceEndpointMethodFilterIdentifier struct {
	ServiceEndpointMethodIdentifier
	FilterId string
}

func (semfi *ServiceEndpointMethodFilterIdentifier) Inherit(p ServiceEndpointMethodIdentifier) {
	semfi.ServiceId = p.ServiceId
	semfi.EndpointId = p.EndpointId
	semfi.MethodId = p.MethodId
}

func (emfi *ServiceEndpointMethodFilterIdentifier) Id() string {
	return CreateCompoundId(emfi.ServiceId, emfi.EndpointId, emfi.MethodId, emfi.FilterId)
}

func (emfi *ServiceEndpointMethodFilterIdentifier) From(id string) {
	ParseCompoundId(id, &emfi.ServiceId, &emfi.EndpointId, &emfi.MethodId, &emfi.FilterId)
}

func (emfi *ServiceEndpointMethodFilterIdentifier) IsIdentified() bool {
	return emfi.ServiceEndpointMethodIdentifier.IsIdentified() && len(emfi.FilterId) > 0
}

func MashEndpointMethodFilterUpsertable(d *schema.ResourceData) (masherytypes.MasheryResponseFilter, diag.Diagnostics) {
	emfi := ServiceEndpointMethodFilterIdentifier{}
	emfi.From(d.Id())

	rvd := diag.Diagnostics{}

	rv := masherytypes.MasheryResponseFilter{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   emfi.FilterId,
			Name: extractString(d, MashObjName, "Terraform filter"),
		},
		Notes:            extractString(d, MashServiceEndpointMethodFilterNotes, "Managed by terraform"),
		XmlFilterFields:  strings.Join(ExtractStringArray(d, MashServiceEndpointMethodFilterXmlFields, &EmptyStringArray), endpointFilterFieldSeparator),
		JsonFilterFields: strings.Join(ExtractStringArray(d, MashServiceEndpointMethodFilterJsonFields, &EmptyStringArray), endpointFilterFieldSeparator),
	}

	if len(rv.JsonFilterFields) > 0 && !JsonRegexp.MatchString(rv.JsonFilterFields) {
		rvd = append(rvd, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid JSON filter element(s)",
			Detail:        fmt.Sprintf("Resulting filter (%s) is malformed", rv.JsonFilterFields),
			AttributePath: cty.GetAttrPath(MashServiceEndpointMethodFilterJsonFields),
		})
	}

	if len(rv.XmlFilterFields) > 0 && !JsonRegexp.MatchString(rv.XmlFilterFields) {
		rvd = append(rvd, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid XML filter element(s)",
			Detail:        fmt.Sprintf("Resulting filter (%s) is malformed", rv.JsonFilterFields),
			AttributePath: cty.GetAttrPath(MashServiceEndpointMethodFilterXmlFields),
		})
	}

	return rv, rvd
}

func V3EndpointMethodFilterToResourceData(inp *masherytypes.MasheryResponseFilter, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashEndpointMethodFilterId:                inp.Id,
		MashObjName:                               inp.Name,
		MashObjCreated:                            inp.Created.ToString(),
		MashObjUpdated:                            inp.Updated.ToString(),
		MashServiceEndpointMethodFilterNotes:      inp.Notes,
		MashServiceEndpointMethodFilterXmlFields:  nilArrayForEmptyString(inp.XmlFilterFields, endpointFilterFieldSeparator),
		MashServiceEndpointMethodFilterJsonFields: nilArrayForEmptyString(inp.JsonFilterFields, endpointFilterFieldSeparator),
	}

	return SetResourceFields(data, d)
}

func initEndpointMethodFilterSchemaBoilerplate() {
	addComputedString(&EndpointMethodFilterSchema, MashEndpointMethodFilterId, "V3 Id of this filter")
	addRequiredString(&EndpointMethodFilterSchema, MashObjName, "Filter name")

	// Created and updated fields
	addComputedString(&EndpointMethodFilterSchema, MashObjCreated, "Date/time the object was created")
	addComputedString(&EndpointMethodFilterSchema, MashObjUpdated, "Date/time the object was updated")
}

func init() {
	initEndpointMethodFilterSchemaBoilerplate()
}
