package mashres

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/mashschemag"
	"testing"
)

func TestCreateErrorSetMessageWillReturnDiagOnDefaultMessage(t *testing.T) {
	h := CreateTestResource(ServiceErrorSetResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	obj1 := mashschemag.ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         403,
		DetailHeader: "Account Inactive",
		ResponseBody: "",
		Status:       "Forbidden",
		Id:           "ERR_403_DEVELOPER_INACTIVE",
	})
	obj2 := mashschemag.ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         414,
		DetailHeader: "detail",
		ResponseBody: "response body",
		Status:       "Request-URI Too Long",
		Id:           "ERR_414_REQUEST_URI_TOO_LONG",
	})

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashSvcErrorSetMessage, []interface{}{obj1, obj2})
	h.thenExecutingCreateWillYieldDiagnostic(t, "invalid input for field error_message")

	// ---------------------------------------------------
	// Double-check that unsupported upgrades will be disabled.

	ident := masherytypes.ErrorSetIdentifier{
		ErrorSetId: "setId",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "sid",
		},
	}
	h.givenIdentity(t, ident)
	h.thenExecutingUpdateWillYieldDiagnostic(t, "invalid input for field error_message")
}

func TestCreateErrorSetMessageWillReturnDiagOnUnsupportedErrorMessageCode(t *testing.T) {
	h := CreateTestResource(ServiceErrorSetResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	obj1 := mashschemag.ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         403,
		DetailHeader: "Account Inactive",
		ResponseBody: "",
		Status:       "Forbidden",
		Id:           "ERR_4033_DEVELOPER_INACTIVE",
	})
	obj2 := mashschemag.ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         414,
		DetailHeader: "detail",
		ResponseBody: "response body",
		Status:       "Request-URI Too Long",
		Id:           "ERR_414_REQUEST_URI_TOO_LONG",
	})

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashSvcErrorSetMessage, []interface{}{obj1, obj2})
	h.thenExecutingCreateWillYieldDiagnostic(t, "invalid input for field error_message")

	// ---------------------------------------------------
	// Double-check that unsupported upgrades will be disabled.

	ident := masherytypes.ErrorSetIdentifier{
		ErrorSetId: "setId",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "sid",
		},
	}
	h.givenIdentity(t, ident)
	h.thenExecutingUpdateWillYieldDiagnostic(t, "invalid input for field error_message")
}

func TestCreateErrorSetMessageWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceErrorSetResource)

	serviceIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	}

	obj1 := mashschemag.ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         403,
		DetailHeader: "Account Inactive",
		ResponseBody: "Custom Response Body",
		Status:       "Forbidden",
		Id:           "ERR_403_DEVELOPER_INACTIVE",
	})
	obj2 := mashschemag.ErrorMessageToTerraform(masherytypes.MasheryErrorMessage{
		Code:         414,
		DetailHeader: "detail",
		ResponseBody: "response body",
		Status:       "Request-URI Too Long",
		Id:           "ERR_414_REQUEST_URI_TOO_LONG",
	})

	h.givenParentIdentity(t, serviceIdent)
	h.givenStateFieldSetTo(t, mashschema.MashSvcErrorSetMessage, []interface{}{obj1, obj2})
	givenCreateErrorSetSucceeds(h, serviceIdent)
	h.thenExecutingCreate(t)

	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ErrorSetIdentifier) {
		assert.Equal(t, "abc", id.ServiceId)
		assert.Equal(t, "setId", id.ErrorSetId)
	})
}

func TestReadingErrorSetMessageWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceErrorSetResource)

	errIdent := masherytypes.ErrorSetIdentifier{
		ErrorSetId: "setId",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "abc",
		},
	}

	h.givenIdentity(t, errIdent)
	h.givenParentIdentity(t, errIdent.ServiceIdentifier)
	givenReadErrorSetSucceeds(h, errIdent)
	h.thenExecutingRead(t)

	// Then it succeeds
}

func givenCreateErrorSetSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet], serviceId masherytypes.ServiceIdentifier) {
	rv := masherytypes.ErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "setId",
			Name: "the-set-name",
		},
		Type:            "type",
		JSONP:           false,
		JSONPType:       "",
		ErrorMessages:   nil,
		ParentServiceId: serviceId,
	}

	h.mockClientWill().
		On("CreateErrorSet", mock.Anything, serviceId, mock.Anything).
		Return(rv, nil).
		Once()
}

func givenReadErrorSetSucceeds(h *ResourceTemplateMockHelper[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet], errId masherytypes.ErrorSetIdentifier) {
	rv := masherytypes.ErrorSet{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "setId",
			Name: "the-set-name",
		},
		Type:      "type",
		JSONP:     false,
		JSONPType: "",
		ErrorMessages: &[]masherytypes.MasheryErrorMessage{
			{
				Id:           "ERR_400_UNSUPPORTED_PARAMETER",
				Status:       "dStatus",
				DetailHeader: "detail",
				ResponseBody: "response body",
				Code:         400,
			},
		},
		ParentServiceId: errId.ServiceIdentifier,
	}

	h.mockClientWill().
		On("GetErrorSet", mock.Anything, errId).
		Return(rv, true, nil).
		Once()
}
