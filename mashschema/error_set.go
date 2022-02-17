package mashschema

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strconv"
)

var validErrorSetMessageId = []string{
	"ERR_400_BAD_REQUEST", "ERR_403_NOT_AUTHORIZED", "ERR_403_DEVELOPER_INACTIVE", "ERR_403_DEVELOPER_OVER_QPS",
	"ERR_403_DEVELOPER_OVER_RATE", "ERR_403_DEVELOPER_UNKNOWN_REFERER", "ERR_403_SERVICE_OVER_QPS", "ERR_403_SERVICE_REQUIRES_SSL",
	"ERR_414_REQUEST_URI_TOO_LONG", "ERR_502_BAD_GATEWAY", "ERR_503_SERVICE_UNAVAILABLE", "ERR_504_GATEWAY_TIMEOUT",
	"ERR_400_UNSUPPORTED_PARAMETER", "ERR_400_UNSUPPORTED_SIGNATURE_METHOD", "ERR_400_MISSING_REQUIRED_CONSUMER_KEY",
	"ERR_400_MISSING_REQUIRED_REQUEST_TOKEN", "ERR_400_MISSING_REQUIRED_ACCESS_TOKEN", "ERR_400_DUPLICATED_OAUTH_PROTOCOL_PARAMETER",
	"ERR_401_TIMESTAMP_IS_INVALID", "ERR_401_INVALID_SIGNATURE", "ERR_401_INVALID_OR_EXPIRED_TOKEN", "ERR_401_INVALID_CONSUMER_KEY",
	"ERR_401_INVALID_NONCE",
}

const (
	MashSvcErrorSetType      = "type"
	MashSvcErrorSetJsonp     = "jsonp"
	MashSvcErrorSetJsonpType = "jsonp_type"
	MashSvcErrorSetMessage   = "error_message"

	MashSvrErrorSetMessageStatus       = "status"
	MashSvrErrorSetMessageDetailHeader = "detail_header"
	MashSvrErrorSetMessageResponseBody = "response_body"
)

var ServiceErrorSetSchema = map[string]*schema.Schema{
	MashObjId: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Id of this error set",
	},
	MashSvcId: {
		Type:     schema.TypeString,
		Required: true,
	},
	MashObjName: {
		Type:     schema.TypeString,
		Required: true,
	},
	MashSvcErrorSetJsonpType: {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	MashSvcErrorSetJsonp: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	MashSvcErrorSetType: {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "text/xml",
	},
	MashSvcErrorSetMessage: {
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				MashObjId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Id of this message",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return validateStringValueInSet(i, path, &validErrorSetMessageId)
					},
				},
				MashSvrErrorSetMessageStatus: {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "String passed in the status field",
				},
				MashSvrErrorSetMessageDetailHeader: {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Detailed header contents",
				},
				MashSvrErrorSetMessageResponseBody: {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Response body value",
				},
			},
		},
	},
}

type ErrorSetIdentifier struct {
	ServiceIdentifier
	ErrorSetId string
}

func (esi *ErrorSetIdentifier) Self() interface{} {
	return esi
}

var ErrorSetMapper *ErrorSetMapperImpl

type ErrorSetMapperImpl struct {
	MapperImpl
}

func (esm *ErrorSetMapperImpl) GetServiceIdentifier(d *schema.ResourceData) string {
	return ExtractString(d, MashSvcId, "")
}

func (esm *ErrorSetMapperImpl) GetIdentifier(d *schema.ResourceData) *ErrorSetIdentifier {
	ident := &ErrorSetIdentifier{}
	CompoundIdFrom(ident, d.Id())

	return ident
}

func (esm *ErrorSetMapperImpl) ErrorSetMessagesChanged(d *schema.ResourceData) bool {
	return d.HasChanges(MashSvcErrorSetMessage)
}

func (esm *ErrorSetMapperImpl) FindMessageById(inp *[]masherytypes.MasheryErrorMessage, id string) *masherytypes.MasheryErrorMessage {
	for _, val := range *inp {
		if id == val.Id {
			return &val
		}
	}

	return nil
}

func (esm *ErrorSetMapperImpl) convertExistingErrorSubset(errSet *masherytypes.MasheryErrorSet, d *schema.ResourceData) []map[string]interface{} {
	var rv []map[string]interface{}

	defined := esm.UpsertableErrorMessages(d)
	for _, inboundMsg := range *errSet.ErrorMessages {
		if ptr := esm.FindMessageById(&defined, inboundMsg.Id); ptr != nil {
			rv = append(rv, esm.PersistErrorMessage(inboundMsg))
		}
	}

	return rv
}

func (esm *ErrorSetMapperImpl) PersistTyped(ctx context.Context, errSet *masherytypes.MasheryErrorSet, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashObjName:              errSet.Name,
		MashSvcErrorSetJsonp:     errSet.JSONP,
		MashSvcErrorSetJsonpType: errSet.JSONPType,
		MashSvcErrorSetType:      errSet.Type,
	}

	if d.Get(MashSvcErrorSetMessage) != nil && errSet.ErrorMessages != nil {
		data[MashSvcErrorSetMessage] = esm.convertExistingErrorSubset(errSet, d)
	}

	return esm.SetResourceFields(ctx, data, d)
}

func (esm *ErrorSetMapperImpl) UpsertableTyped(d *schema.ResourceData) masherytypes.MasheryErrorSet {
	setIdent := ErrorSetIdentifier{}
	CompoundIdFrom(&setIdent, d.Id())

	rv := masherytypes.MasheryErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   setIdent.ErrorSetId,
			Name: ExtractString(d, MashObjName, "Terraform-managed error set"),
		},
		JSONPType: ExtractString(d, MashSvcErrorSetJsonpType, ""),
		JSONP:     extractBool(d, MashSvcErrorSetJsonp, false),
		Type:      ExtractString(d, MashSvcErrorSetType, "text/xml"),
	}

	return rv
}

func extractKeyFromMap(inp map[string]interface{}, key string, receiver *string) {
	if valRaw := inp[key]; valRaw != nil {
		if str, ok := valRaw.(string); ok {
			*receiver = str
		}
	}
}

func extractIntKeyFromMap(inp map[string]interface{}, key string, receiver *int) {
	if valRaw := inp[key]; valRaw != nil {
		if intVal, ok := valRaw.(int); ok {
			*receiver = intVal
		}
	}
}

var errorParsePattern = regexp.MustCompile("ERR_(\\d{3})_.*")

func (esm *ErrorSetMapperImpl) SetIdentifier(serviceId string, errSet *masherytypes.MasheryErrorSet, d *schema.ResourceData) {
	ident := &ErrorSetIdentifier{
		ServiceIdentifier: ServiceIdentifier{
			ServiceId: serviceId,
		},
		ErrorSetId: errSet.Id,
	}

	d.SetId(CompoundId(ident))
}

func (esm *ErrorSetMapperImpl) UpsertableErrorMessage(d interface{}) masherytypes.MasheryErrorMessage {
	rv := masherytypes.MasheryErrorMessage{}

	if mp, ok := d.(map[string]interface{}); ok {
		extractKeyFromMap(mp, MashObjId, &rv.Id)
		if errorParsePattern.MatchString(rv.Id) {
			match := errorParsePattern.FindStringSubmatch(rv.Id)
			i, _ := strconv.ParseInt(match[1], 10, 32)
			rv.Code = int(i)
		}

		extractKeyFromMap(mp, MashSvrErrorSetMessageStatus, &rv.Status)
		extractKeyFromMap(mp, MashSvrErrorSetMessageDetailHeader, &rv.DetailHeader)
		extractKeyFromMap(mp, MashSvrErrorSetMessageResponseBody, &rv.ResponseBody)
	}

	return rv
}

func (esm *ErrorSetMapperImpl) PersistErrorMessage(inp masherytypes.MasheryErrorMessage) map[string]interface{} {
	return map[string]interface{}{
		MashObjId:                          inp.Id,
		MashSvrErrorSetMessageStatus:       inp.Status,
		MashSvrErrorSetMessageDetailHeader: inp.DetailHeader,
		MashSvrErrorSetMessageResponseBody: inp.ResponseBody,

		// Code is inferred from Id.
	}
}

func (esm *ErrorSetMapperImpl) UpsertableErrorMessages(d *schema.ResourceData) []masherytypes.MasheryErrorMessage {
	if rawSet, ok := d.GetOk(MashSvcErrorSetMessage); ok {
		if msgSet, ok := rawSet.(*schema.Set); ok {
			rv := make([]masherytypes.MasheryErrorMessage, msgSet.Len())
			for i, iMsg := range msgSet.List() {
				rv[i] = esm.UpsertableErrorMessage(iMsg)
			}

			return rv
		}
	}

	return []masherytypes.MasheryErrorMessage{}
}
