package mashschema

import (
	"regexp"
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

var ErrorParsePattern = regexp.MustCompile("ERR_(\\d{3})_.*")
