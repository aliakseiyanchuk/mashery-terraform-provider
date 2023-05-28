package mashres

import (
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestQueryingEmailTemplateSetWillSucceedOnUniqueMatch(t *testing.T) {
	params := map[string]string{"name": "email-set-name"}
	h := CreateTestDatasource(EmailTemplateSetDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenEmailTemplateSetIsReturned(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.willHaveStateId(t, "set-id")
}

func TestQueryingEmailTemplateSetWillReturnDiagnosticOnAPIError(t *testing.T) {
	params := map[string]string{"name": "email-set-name"}
	h := CreateTestDatasource(EmailTemplateSetDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	givenQueryingEmailTemplatesReturnsError(h, params)

	h.thenExecutingDataSourceQueryWillYieldDiagnostic(t, "query has returned an error: sample rejection for email template set")
}

func TestQueryingEmailTemplateSetWillReturnDiagnosticOnNoMatch(t *testing.T) {
	params := map[string]string{"name": "email-set-name"}
	h := CreateTestDatasource(EmailTemplateSetDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	h.givenStateFieldSetTo(t, mashschema.MashDataSourceRequired, true)
	givenQueryingEmailTemplatesReturnsNothing(h, params)

	h.thenExecutingDataSourceQueryWillYieldDiagnostic(t, "no matching object was found, however the configuration requires a match")
}

func TestQueryingEmailTemplateSetWillProcessOptionalMatch(t *testing.T) {
	params := map[string]string{"name": "email-set-name"}
	h := CreateTestDatasource(EmailTemplateSetDataSource)

	h.givenStateFieldSetTo(t, mashschema.MashDataSourceSearch, params)
	h.givenStateFieldSetTo(t, mashschema.MashDataSourceRequired, false)
	givenQueryingEmailTemplatesReturnsNothing(h, params)

	h.thenExecutingDataSourceQuery(t)
	h.willHaveStateId(t, "")
}

func givenEmailTemplateSetIsReturned(h *DatasourceTemplateMockHelper, params map[string]string) {
	returnedSet := []masherytypes.EmailTemplateSet{
		{
			AddressableV3Object: masherytypes.AddressableV3Object{
				Id:        "set-id",
				Name:      params["name"],
				Created:   nil,
				Updated:   nil,
				Retrieved: time.Time{},
			},
			Type:           "Some-Type",
			EmailTemplates: nil,
		},
	}

	h.mockClientWill().
		On("ListEmailTemplateSetsFiltered", mock.Anything, params, mashschema.EmptyStringArray).
		Return(returnedSet, nil)
}

func givenQueryingEmailTemplatesReturnsError(h *DatasourceTemplateMockHelper, params map[string]string) {
	h.mockClientWill().
		On("ListEmailTemplateSetsFiltered", mock.Anything, params, mashschema.EmptyStringArray).
		Return([]masherytypes.EmailTemplateSet{}, errors.New("sample rejection for email template set"))
}

func givenQueryingEmailTemplatesReturnsNothing(h *DatasourceTemplateMockHelper, params map[string]string) {
	h.mockClientWill().
		On("ListEmailTemplateSetsFiltered", mock.Anything, params, mashschema.EmptyStringArray).
		Return([]masherytypes.EmailTemplateSet{}, nil)
}
