package mashschema

import (
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

var ServiceEndpointMethodFilterMapper *ServiceEndpointMethodFilterMapperImpl

type ServiceEndpointMethodFilterMapperImpl struct {
	ResourceMapperImpl
}

func (sefm *ServiceEndpointMethodFilterMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.ServiceEndpointMethodFilter, masherytypes.ServiceEndpointMethodIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}

	ctxIdent := masherytypes.ServiceEndpointMethodIdentifier{}
	if !CompoundIdFrom(&ctxIdent, ExtractString(d, MashServiceEndpointMethodId, "")) {
		rvd = append(rvd, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "endpoint method identifier is incomplete",
		})
	}

	emfi := masherytypes.ServiceEndpointMethodFilterIdentifier{}
	identCompete := CompoundIdFrom(&emfi, d.Id())

	parentSelector := func() masherytypes.ServiceEndpointMethodIdentifier {
		if identCompete {
			return emfi.ServiceEndpointMethodIdentifier
		} else {
			return ctxIdent
		}
	}

	rv := masherytypes.ServiceEndpointMethodFilter{
		ResponseFilter: masherytypes.ResponseFilter{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:   emfi.FilterId,
				Name: ExtractString(d, MashObjName, "Terraform filter"),
			},
			Notes:            ExtractString(d, MashServiceEndpointMethodFilterNotes, "Managed by terraform"),
			XmlFilterFields:  strings.Join(ExtractStringArray(d, MashServiceEndpointMethodFilterXmlFields, &EmptyStringArray), endpointFilterFieldSeparator),
			JsonFilterFields: strings.Join(ExtractStringArray(d, MashServiceEndpointMethodFilterJsonFields, &EmptyStringArray), endpointFilterFieldSeparator),
		},

		ServiceEndpointMethod: parentSelector(),
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

	return rv, ctxIdent, rvd
}

func (sefm *ServiceEndpointMethodFilterMapperImpl) PersistTyped(inp masherytypes.ServiceEndpointMethodFilter, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashEndpointMethodFilterId:                inp.Id,
		MashObjName:                               inp.Name,
		MashObjCreated:                            inp.Created.ToString(),
		MashObjUpdated:                            inp.Updated.ToString(),
		MashServiceEndpointMethodFilterNotes:      inp.Notes,
		MashServiceEndpointMethodFilterXmlFields:  nilArrayForEmptyString(inp.XmlFilterFields, endpointFilterFieldSeparator),
		MashServiceEndpointMethodFilterJsonFields: nilArrayForEmptyString(inp.JsonFilterFields, endpointFilterFieldSeparator),
	}

	return sefm.persistMap(inp.Identifier(), data, d)
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
		ResourceMapperImpl: ResourceMapperImpl{
			v3ObjectName: "endpoint method filter",
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

			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv := masherytypes.ServiceEndpointMethodFilterIdentifier{}
				if CompoundIdFrom(&rv, d.Id()) {
					return rv, nil
				} else {
					return rv, diag.Diagnostics{diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "endpoint method filter identifier is incomplete",
					}}
				}
			},

			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return ServiceEndpointMethodFilterMapper.Upsertable(d)
			},

			persistFunc: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				ptr := rv.(*masherytypes.ServiceEndpointMethodFilter)
				return ServiceEndpointMethodFilterMapper.PersistTyped(*ptr, d)
			},
		},
	}

	initEndpointMethodFilterSchemaBoilerplate()
}
