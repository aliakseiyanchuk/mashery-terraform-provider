package mashschema

import (
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strconv"
)

var ValidErrorSetMessageId = []string{
	"ERR_400_BAD_REQUEST", "ERR_403_NOT_AUTHORIZED", "ERR_403_DEVELOPER_INACTIVE", "ERR_403_DEVELOPER_OVER_QPS",
	"ERR_403_DEVELOPER_OVER_RATE", "ERR_403_DEVELOPER_UNKNOWN_REFERER", "ERR_403_SERVICE_OVER_QPS", "ERR_403_SERVICE_REQUIRES_SSL",
	"ERR_414_REQUEST_URI_TOO_LONG", "ERR_502_BAD_GATEWAY", "ERR_503_SERVICE_UNAVAILABLE", "ERR_504_GATEWAY_TIMEOUT",
	"ERR_400_UNSUPPORTED_PARAMETER", "ERR_400_UNSUPPORTED_SIGNATURE_METHOD", "ERR_400_MISSING_REQUIRED_CONSUMER_KEY",
	"ERR_400_MISSING_REQUIRED_REQUEST_TOKEN", "ERR_400_MISSING_REQUIRED_ACCESS_TOKEN", "ERR_400_DUPLICATED_OAUTH_PROTOCOL_PARAMETER",
	"ERR_401_TIMESTAMP_IS_INVALID", "ERR_401_INVALID_SIGNATURE", "ERR_401_INVALID_OR_EXPIRED_TOKEN", "ERR_401_INVALID_CONSUMER_KEY",
	"ERR_401_INVALID_NONCE",
}

const (
	MashSvcErrorSetId        = "error_set_id"
	MashSvcErrorSetType      = "type"
	MashSvcErrorSetJsonp     = "jsonp"
	MashSvcErrorSetJsonpType = "jsonp_type"
	MashSvcErrorSetMessage   = "error_message"

	MashSvcErrorSetMessageErrorId      = "error"
	MashSvcErrorSetMessageStatus       = "status"
	MashSvcErrorSetMessageDetailHeader = "detail_header"
	MashSvcErrorSetMessageResponseBody = "response_body"
)

var ErrorSetMapper *ErrorSetMapperImpl

type ErrorSetMapperImpl struct {
	ResourceMapperImpl
}

func (esm *ErrorSetMapperImpl) ErrorSetIdentityDiag(d *schema.ResourceData) (masherytypes.ErrorSetIdentifier, diag.Diagnostics) {
	rv := masherytypes.ErrorSetIdentifier{}
	if CompoundIdFrom(&rv, d.Id()) {
		return rv, nil
	} else {
		return rv, diag.Diagnostics{ErrorSetMapper.lackingIdentificationDiagnostic("id")}
	}
}
func (esm *ErrorSetMapperImpl) ErrorSetIdentity(d *schema.ResourceData) (masherytypes.ErrorSetIdentifier, error) {
	rv := masherytypes.ErrorSetIdentifier{}
	if CompoundIdFrom(&rv, d.Id()) {
		return rv, nil
	} else {
		return rv, errors.New("error set identity is incomplete")
	}
}

func (esm *ErrorSetMapperImpl) GetServiceIdentifier(d *schema.ResourceData) string {
	return ExtractString(d, MashSvcId, "")
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

func (esm *ErrorSetMapperImpl) convertExistingErrorSubset(errSet *masherytypes.ErrorSet, d *schema.ResourceData) []map[string]interface{} {
	var rv []map[string]interface{}

	defined := esm.UpsertableErrorMessages(d)
	for _, inboundMsg := range *errSet.ErrorMessages {
		if ptr := esm.FindMessageById(&defined, inboundMsg.Id); ptr != nil {
			rv = append(rv, esm.PersistErrorMessage(inboundMsg))
		}
	}

	return rv
}

func (esm *ErrorSetMapperImpl) PersistTyped(errSet *masherytypes.ErrorSet, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashObjName:              errSet.Name,
		MashSvcErrorSetJsonp:     errSet.JSONP,
		MashSvcErrorSetJsonpType: errSet.JSONPType,
		MashSvcErrorSetType:      errSet.Type,
	}

	if d.Get(MashSvcErrorSetMessage) != nil && errSet.ErrorMessages != nil {
		data[MashSvcErrorSetMessage] = esm.convertExistingErrorSubset(errSet, d)
	}

	return SetResourceFields(data, d)
}

func (esm *ErrorSetMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.ErrorSet, masherytypes.ServiceIdentifier, diag.Diagnostics) {
	rvd := diag.Diagnostics{}

	setIdent := masherytypes.ErrorSetIdentifier{}
	identComplete := CompoundIdFrom(&setIdent, d.Id())

	svcIdent := masherytypes.ServiceIdentifier{}
	if !CompoundIdFrom(&svcIdent, ExtractString(d, MashSvcId, "")) {
		rvd = append(rvd, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "service identifier is incomplete",
		})
	}

	parentSelector := func() masherytypes.ServiceIdentifier {
		if identComplete {
			return setIdent.ServiceIdentifier
		} else {
			return svcIdent
		}
	}

	rv := masherytypes.ErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   setIdent.ErrorSetId,
			Name: ExtractString(d, MashObjName, "Terraform-managed error set"),
		},
		JSONPType: ExtractString(d, MashSvcErrorSetJsonpType, ""),
		JSONP:     ExtractBool(d, MashSvcErrorSetJsonp, false),
		Type:      ExtractString(d, MashSvcErrorSetType, "text/xml"),

		ParentServiceId: parentSelector(),
	}

	return rv, svcIdent, rvd
}

func ExtractKeyFromMap(inp map[string]interface{}, key string, receiver *string) {
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

var ErrorParsePattern = regexp.MustCompile("ERR_(\\d{3})_.*")

func (esm *ErrorSetMapperImpl) UpsertableErrorMessage(d interface{}) masherytypes.MasheryErrorMessage {
	rv := masherytypes.MasheryErrorMessage{}

	if mp, ok := d.(map[string]interface{}); ok {
		ExtractKeyFromMap(mp, MashObjId, &rv.Id)
		if ErrorParsePattern.MatchString(rv.Id) {
			match := ErrorParsePattern.FindStringSubmatch(rv.Id)
			i, _ := strconv.ParseInt(match[1], 10, 32)
			rv.Code = int(i)
		}

		ExtractKeyFromMap(mp, MashSvcErrorSetMessageStatus, &rv.Status)
		ExtractKeyFromMap(mp, MashSvcErrorSetMessageDetailHeader, &rv.DetailHeader)
		ExtractKeyFromMap(mp, MashSvcErrorSetMessageResponseBody, &rv.ResponseBody)
	}

	return rv
}

func (esm *ErrorSetMapperImpl) PersistErrorMessage(inp masherytypes.MasheryErrorMessage) map[string]interface{} {
	return map[string]interface{}{
		MashObjId:                          inp.Id,
		MashSvcErrorSetMessageStatus:       inp.Status,
		MashSvcErrorSetMessageDetailHeader: inp.DetailHeader,
		MashSvcErrorSetMessageResponseBody: inp.ResponseBody,

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

func init() {
	ErrorSetMapper = &ErrorSetMapperImpl{
		ResourceMapperImpl{
			v3ObjectName: "error set",
			schema: map[string]*schema.Schema{
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
									return ValidateStringValueInSet(i, path, &ValidErrorSetMessageId)
								},
							},
							MashSvcErrorSetMessageStatus: {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "String passed in the status field",
							},
							MashSvcErrorSetMessageDetailHeader: {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "Detailed header contents",
							},
							MashSvcErrorSetMessageResponseBody: {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "Response body value",
							},
						},
					},
				},
			},

			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				return ErrorSetMapper.ErrorSetIdentityDiag(d)
			},
		},
	}
}
