package mashres

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DatasourceTemplateMockHelper struct {
	template DatasourceTemplate
	schema   *schema.Resource

	data   *schema.ResourceData
	mockCl MockClient
}

func (dtmh *DatasourceTemplateMockHelper) givenStateFieldSetTo(t *testing.T, field string, value interface{}) {
	err := dtmh.data.Set(field, value)
	assert.Nil(t, err)
}

func (dtmh *DatasourceTemplateMockHelper) willHaveStateId(t *testing.T, v string) {
	assert.Equal(t, dtmh.data.Id(), v)
}

func (dtmh *DatasourceTemplateMockHelper) mockClientWill() *MockClient {
	return &dtmh.mockCl
}

func (dtmh *DatasourceTemplateMockHelper) thenExecutingDataSourceQuery(t *testing.T) {
	dg := RoleDataSource.Query(context.TODO(), dtmh.data, &dtmh.mockCl)
	assert.True(t, len(dg) == 0)

	dtmh.mockCl.AssertExpectations(t)
}

func (dtmh *DatasourceTemplateMockHelper) thenExecutingDataSourceQueryWillYieldDiagnostic(t *testing.T, dgText string) {
	dg := RoleDataSource.Query(context.TODO(), dtmh.data, &dtmh.mockCl)
	assert.True(t, len(dg) == 1)
	assert.Equal(t, diag.Error, dg[0].Severity)
	assert.Equal(t, dgText, dg[0].Summary)

	dtmh.mockCl.AssertExpectations(t)
}

func (dtmh *DatasourceTemplateMockHelper) willHaveFieldSetTo(t *testing.T, name string, s interface{}) {
	if v := dtmh.data.Get(name); v != nil {
		assert.Equal(t, s, v)
	} else {
		assert.Failf(t, "field is not found in state", "field %s is not found in state", name)
	}
}

func CreateTestDatasource(tmpl DatasourceTemplate) *DatasourceTemplateMockHelper {
	rv := DatasourceTemplateMockHelper{
		template: tmpl,
		schema:   tmpl.DataSourceSchema(),
		data:     tmpl.TestState(),
		mockCl:   MockClient{},
	}

	return &rv
}
