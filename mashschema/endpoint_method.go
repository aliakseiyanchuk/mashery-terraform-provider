package mashschema

import (
	"context"
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

type ServiceEndpointMethodIdentifier struct {
	ServiceEndpointIdentifier
	MethodId string
}

func (semi *ServiceEndpointMethodIdentifier) Self() interface{} {
	return semi
}

var ServiceEndpointMethodMapper *ServiceEndpointMethodMapperImpl

type ServiceEndpointMethodMapperImpl struct {
	MapperImpl
}

func (semm *ServiceEndpointMethodMapperImpl) CreateIdentifierTyped() *ServiceEndpointMethodIdentifier {
	return &ServiceEndpointMethodIdentifier{}
}

func (semm *ServiceEndpointMethodMapperImpl) CreateIdentifier(d *schema.ResourceData) *ServiceEndpointMethodIdentifier {
	rv := semm.CreateIdentifierTyped()
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (semm *ServiceEndpointMethodMapperImpl) EndpointIdentifier(d *schema.ResourceData) (*ServiceEndpointIdentifier, diag.Diagnostics) {
	rv := &ServiceEndpointIdentifier{}
	CompoundIdFrom(rv, ExtractString(d, MashEndpointId, ""))

	if !IsIdentified(rv) {
		return rv, diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Incomplete identifier",
			Detail:        "Endpoint identifier supplies incomplete data, or is malformed",
			AttributePath: cty.GetAttrPath(MashEndpointId),
		}}
	} else {
		return rv, nil
	}
}

func (semm *ServiceEndpointMethodMapperImpl) UpsertableTyped(d *schema.ResourceData) (*masherytypes.MasheryMethod, diag.Diagnostics) {
	ident := ServiceEndpointMethodIdentifier{}
	CompoundIdFrom(&ident, d.Id())

	return &masherytypes.MasheryMethod{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   ident.MethodId,
			Name: ExtractString(d, MashObjName, ""),
		},
		SampleJsonResponse: ExtractString(d, MashServiceEndpointMethodSampleJson, ""),
		SampleXmlResponse:  ExtractString(d, MashServiceEndpointMethodSampleXml, ""),
	}, nil
}

func (semm *ServiceEndpointMethodMapperImpl) PersistTyped(ctx context.Context, inp *masherytypes.MasheryMethod, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashServiceEndpointMethodId:         inp.Id,
		MashObjCreated:                      inp.Created.ToString(),
		MashObjUpdated:                      inp.Updated.ToString(),
		MashServiceEndpointMethodSampleJson: nullForEmptyString(inp.SampleJsonResponse),
		MashServiceEndpointMethodSampleXml:  nullForEmptyString(inp.SampleXmlResponse),
	}

	return semm.SetResourceFields(ctx, data, d)
}

func initEndpointMethodSchemaBoilerplate() {
	addComputedString(&ServiceEndpointMethodMapper.schema, ServiceEndpointMethodRef, "V3 identifier of this method")
	addComputedString(&ServiceEndpointMethodMapper.schema, MashObjCreated, "Date/time the object was created")
	addComputedString(&ServiceEndpointMethodMapper.schema, MashObjUpdated, "Date/time the object was updated")

	addOptionalString(&ServiceEndpointMethodMapper.schema, MashServiceEndpointMethodSampleJson, "Sample JSON response")
	addOptionalString(&ServiceEndpointMethodMapper.schema, MashServiceEndpointMethodSampleXml, "Sample XML response")
}

func init() {
	ServiceEndpointMethodMapper = &ServiceEndpointMethodMapperImpl{
		MapperImpl: MapperImpl{
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
		},
	}

	initEndpointMethodSchemaBoilerplate()

	ServiceEndpointMethodMapper.identifier = func() interface{} {
		return &ServiceEndpointMethodIdentifier{}
	}

	ServiceEndpointMethodMapper.persistFunc = func(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics {
		return ServiceEndpointMethodMapper.PersistTyped(ctx, rv.(*masherytypes.MasheryMethod), d)
	}

	ServiceEndpointMethodMapper.upsertFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return ServiceEndpointMethodMapper.UpsertableTyped(d)
	}
}
