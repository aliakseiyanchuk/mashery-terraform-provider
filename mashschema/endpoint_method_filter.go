package mashschema

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strings"
)

const (
	MashEndpointMethodFilterId           = "filter_id"
	ServiceEndpointMethodFilterRef       = "service_filter_id"
	MashServiceEndpointMethodFilterNotes = "notes"

	MashServiceEndpointMethodFilterXmlFields  = "xml_fields"
	MashServiceEndpointMethodFilterJsonFields = "json_fields"

	endpointFilterFieldSeparator = " | "
)

var JsonRegexp = regexp.MustCompile("^(?:\\/[^\\/\\| ]+)+(?: \\| (?:\\/[^\\/\\| ]+)+)*$")

type ServiceEndpointMethodFilterIdentifier struct {
	ServiceEndpointMethodIdentifier
	FilterId string
}

func (semfi *ServiceEndpointMethodFilterIdentifier) Self() interface{} {
	return semfi
}

var ServiceEndpointMethodFilterMapper *ServiceEndpointMethodFilterMapperImpl

type ServiceEndpointMethodFilterMapperImpl struct {
	MapperImpl
}

func (sefm *ServiceEndpointMethodFilterMapperImpl) GetMethodIdentifier(d *schema.ResourceData) (*ServiceEndpointMethodIdentifier, diag.Diagnostics) {
	rv := &ServiceEndpointMethodIdentifier{}
	CompoundIdFrom(rv, ExtractString(d, ServiceEndpointMethodRef, ""))

	var rvd diag.Diagnostics = nil
	if !IsIdentified(rv) {
		rvd = CompoundIdMalformedDiagnostic(cty.GetAttrPath(ServiceEndpointMethodRef))
	}

	return rv, rvd
}

func (sefm *ServiceEndpointMethodFilterMapperImpl) GetIdentifier(d *schema.ResourceData) *ServiceEndpointMethodFilterIdentifier {
	rv := &ServiceEndpointMethodFilterIdentifier{}
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (sefm *ServiceEndpointMethodFilterMapperImpl) SetIdentifier(pIdent *ServiceEndpointMethodIdentifier, fl *masherytypes.MasheryResponseFilter, d *schema.ResourceData) {
	v := &ServiceEndpointMethodFilterIdentifier{
		ServiceEndpointMethodIdentifier: ServiceEndpointMethodIdentifier{
			ServiceEndpointIdentifier: ServiceEndpointIdentifier{
				ServiceIdentifier: ServiceIdentifier{
					ServiceId: pIdent.ServiceId,
				},
				EndpointId: pIdent.EndpointId,
			},
			MethodId: pIdent.MethodId,
		},
		FilterId: fl.Id,
	}

	d.SetId(CompoundId(v))
}

func (sefm *ServiceEndpointMethodFilterMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.MasheryResponseFilter, diag.Diagnostics) {
	emfi := ServiceEndpointMethodFilterIdentifier{}
	CompoundIdFrom(&emfi, d.Id())

	rvd := diag.Diagnostics{}

	rv := masherytypes.MasheryResponseFilter{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   emfi.FilterId,
			Name: ExtractString(d, MashObjName, "Terraform filter"),
		},
		Notes:            ExtractString(d, MashServiceEndpointMethodFilterNotes, "Managed by terraform"),
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

func (sefm *ServiceEndpointMethodFilterMapperImpl) PersistTyped(ctx context.Context, inp *masherytypes.MasheryResponseFilter, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashEndpointMethodFilterId:                inp.Id,
		MashObjName:                               inp.Name,
		MashObjCreated:                            inp.Created.ToString(),
		MashObjUpdated:                            inp.Updated.ToString(),
		MashServiceEndpointMethodFilterNotes:      inp.Notes,
		MashServiceEndpointMethodFilterXmlFields:  nilArrayForEmptyString(inp.XmlFilterFields, endpointFilterFieldSeparator),
		MashServiceEndpointMethodFilterJsonFields: nilArrayForEmptyString(inp.JsonFilterFields, endpointFilterFieldSeparator),
	}

	return sefm.SetResourceFields(ctx, data, d)
}

func initEndpointMethodFilterSchemaBoilerplate() {
	addComputedString(&ServiceEndpointMethodFilterMapper.schema, MashEndpointMethodFilterId, "V3 Id of this filter")
	addRequiredString(&ServiceEndpointMethodFilterMapper.schema, MashObjName, "Filter name")

	// Created and updated fields
	addComputedString(&ServiceEndpointMethodFilterMapper.schema, MashObjCreated, "Date/time the object was created")
	addComputedString(&ServiceEndpointMethodFilterMapper.schema, MashObjUpdated, "Date/time the object was updated")
}

func init() {
	ServiceEndpointMethodFilterMapper = &ServiceEndpointMethodFilterMapperImpl{
		MapperImpl: MapperImpl{
			schema: map[string]*schema.Schema{
				MashServiceEndpointMethodFilterNotes: {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "Managed by Terraform",
					Description: "Notes added to this method filter",
				},
				ServiceEndpointMethodRef: {
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
			},
		},
	}
	initEndpointMethodFilterSchemaBoilerplate()
}
