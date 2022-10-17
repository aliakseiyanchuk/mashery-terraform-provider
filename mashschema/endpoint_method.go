package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashServiceEndpointMethodSampleJson = "sample_json"
	MashServiceEndpointMethodSampleXml  = "sample_xml"
	MashServiceEndpointMethodId         = "method_id" // Shared, to be moved
	ServiceEndpointMethodRef            = "method_id" // Shared, to be moved
)

var ServiceEndpointMethodMapper *ServiceEndpointMethodMapperImpl

type ServiceEndpointMethodMapperImpl struct {
	ResourceMapperImpl
}

func (semm *ServiceEndpointMethodMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.ServiceEndpointMethod, masherytypes.ServiceEndpointIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}
	ctxIdent := masherytypes.ServiceEndpointIdentifier{}

	if !CompoundIdFrom(&ctxIdent, ExtractString(d, MashEndpointId, "")) {
		rvd = append(rvd, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "endpoint v3 identitty is not complete",
			Detail:        "Endpoint v3Identity supplies incomplete data, or is malformed",
			AttributePath: cty.GetAttrPath(MashEndpointId),
		})
	}

	ident := masherytypes.ServiceEndpointMethodIdentifier{}
	primaryIdentFull := CompoundIdFrom(&ident, d.Id())

	parentSelector := func() masherytypes.ServiceEndpointIdentifier {
		if primaryIdentFull {
			return ident.ServiceEndpointIdentifier
		} else {
			return ctxIdent
		}
	}

	return masherytypes.ServiceEndpointMethod{
		BaseMethod: masherytypes.BaseMethod{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:   ident.MethodId,
				Name: ExtractString(d, MashObjName, ""),
			},
			SampleJsonResponse: ExtractString(d, MashServiceEndpointMethodSampleJson, ""),
			SampleXmlResponse:  ExtractString(d, MashServiceEndpointMethodSampleXml, ""),
		},

		ParentEndpointId: parentSelector(),
	}, ctxIdent, nil
}

func (semm *ServiceEndpointMethodMapperImpl) PersistTyped(inp masherytypes.ServiceEndpointMethod, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashServiceEndpointMethodId:         inp.Id,
		MashObjCreated:                      inp.Created.ToString(),
		MashObjUpdated:                      inp.Updated.ToString(),
		MashServiceEndpointMethodSampleJson: nullForEmptyString(inp.SampleJsonResponse),
		MashServiceEndpointMethodSampleXml:  nullForEmptyString(inp.SampleXmlResponse),
	}

	return semm.persistMap(inp.Identifier(), data, d)
}

func initEndpointMethodSchemaBoilerplate() {
	addComputedString(&ServiceEndpointMethodMapper.schema, ServiceEndpointMethodRef, "V3 v3Identity of this method")
	addComputedString(&ServiceEndpointMethodMapper.schema, MashObjCreated, "Date/time the object was created")
	addComputedString(&ServiceEndpointMethodMapper.schema, MashObjUpdated, "Date/time the object was updated")

	addOptionalString(&ServiceEndpointMethodMapper.schema, MashServiceEndpointMethodSampleJson, "Sample JSON response")
	addOptionalString(&ServiceEndpointMethodMapper.schema, MashServiceEndpointMethodSampleXml, "Sample XML response")
}

func init() {
	ServiceEndpointMethodMapper = &ServiceEndpointMethodMapperImpl{
		ResourceMapperImpl: ResourceMapperImpl{
			schema: map[string]*schema.Schema{
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
			},

			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv := masherytypes.ServiceEndpointMethodIdentifier{}
				rvd := diag.Diagnostics{}

				if !CompoundIdFrom(&rv, d.Id()) {
					rvd = append(rvd, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "endpoint method's identity is incomplete",
					})
				}

				return rv, rvd
			},

			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return ServiceEndpointMethodMapper.Upsertable(d)
			},

			persistFunc: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				ptr := rv.(*masherytypes.ServiceEndpointMethod)
				return ServiceEndpointMethodMapper.PersistTyped(*ptr, d)
			},
		},
	}

	initEndpointMethodSchemaBoilerplate()
}
