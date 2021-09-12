package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
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

var ErrorSetMessage = map[string]*schema.Schema{
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
}

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
			Schema: ErrorSetMessage,
		},
	},
}

type ErrorSetIdentifier struct {
	ServiceId  string
	ErrorSetId string
}

func (esi *ErrorSetIdentifier) Id() string {
	return CreateCompoundId(esi.ServiceId, esi.ErrorSetId)
}

func (esi *ErrorSetIdentifier) From(id string) {
	ParseCompoundId(id, &esi.ServiceId, &esi.ErrorSetId)
}

func (esi *ErrorSetIdentifier) IsIdentified() bool {
	return len(esi.ServiceId) > 0 && len(esi.ErrorSetId) > 0
}

func findMessageById(inp *[]v3client.MasheryErrorMessage, id string) *v3client.MasheryErrorMessage {
	for _, val := range *inp {
		if id == val.Id {
			return &val
		}
	}

	return nil
}

func convertExistingErrorSubset(errSet *v3client.MasheryErrorSet, d *schema.ResourceData) []map[string]interface{} {
	rv := []map[string]interface{}{}

	defined := V3ErrorSetMessages(d)
	for _, inboundMsg := range *errSet.ErrorMessages {
		if ptr := findMessageById(&defined, inboundMsg.Id); ptr != nil {
			rv = append(rv, V3ErrorMessageForResourceData(inboundMsg))
		}
	}

	return rv
}

func V3ErrorSetToResourceData(errSet *v3client.MasheryErrorSet, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashObjName:              errSet.Name,
		MashSvcErrorSetJsonp:     errSet.JSONP,
		MashSvcErrorSetJsonpType: errSet.JSONPType,
		MashSvcErrorSetType:      errSet.Type,
	}

	if d.Get(MashSvcErrorSetMessage) != nil && errSet.ErrorMessages != nil {
		data[MashSvcErrorSetMessage] = convertExistingErrorSubset(errSet, d)
	}

	return SetResourceFields(data, d)
}

func V3ErrorSetUpsertable(d *schema.ResourceData) v3client.MasheryErrorSet {
	setIdent := ErrorSetIdentifier{}
	setIdent.From(d.Id())

	rv := v3client.MasheryErrorSet{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:   setIdent.ErrorSetId,
			Name: extractString(d, MashObjName, "Terraform-managed error set"),
		},
		JSONPType: extractString(d, MashSvcErrorSetJsonpType, ""),
		JSONP:     extractBool(d, MashSvcErrorSetJsonp, false),
		Type:      extractString(d, MashSvcErrorSetType, "text/xml"),
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

func V3ErrorSetMessage(d interface{}) v3client.MasheryErrorMessage {
	rv := v3client.MasheryErrorMessage{}

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

func V3ErrorMessageForResourceData(inp v3client.MasheryErrorMessage) map[string]interface{} {
	return map[string]interface{}{
		MashObjId:                          inp.Id,
		MashSvrErrorSetMessageStatus:       inp.Status,
		MashSvrErrorSetMessageDetailHeader: inp.DetailHeader,
		MashSvrErrorSetMessageResponseBody: inp.ResponseBody,

		// Code is inferred from Id.
	}
}

func V3ErrorSetMessages(d *schema.ResourceData) []v3client.MasheryErrorMessage {
	if rawSet, ok := d.GetOk(MashSvcErrorSetMessage); ok {
		if msgSet, ok := rawSet.(*schema.Set); ok {
			rv := make([]v3client.MasheryErrorMessage, msgSet.Len())
			for i, iMsg := range msgSet.List() {
				rv[i] = V3ErrorSetMessage(iMsg)
			}

			return rv
		}
	}

	return []v3client.MasheryErrorMessage{}
}
