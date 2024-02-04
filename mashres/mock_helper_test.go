package mashres

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"testing"
)

type MockHelperBase struct {
	schema *schema.Resource

	data   *schema.ResourceData
	mockCl MockClient
}

func (mhb *MockHelperBase) willHaveFieldSetTo(t *testing.T, name string, s interface{}) {
	if v := mhb.data.Get(name); v != nil {
		assert.Equal(t, s, v)
	} else {
		assert.Failf(t, "field is not found in state", "field %s is not found in state", name)
	}
}

func (mhb *MockHelperBase) willHaveStringArrayFieldSetTo(t *testing.T, name string, s []string) {
	if v := mhb.data.Get(name); v != nil {
		arr := mashschema.ExtractStringArray(mhb.data, name, &mashschema.EmptyStringArray)
		assert.True(t, reflect.DeepEqual(s, arr))
	} else {
		assert.Failf(t, "field is not found in state", "field %s is not found in state", name)
	}
}

func (mhb *MockHelperBase) givenStateFieldSetTo(t *testing.T, field string, value interface{}) {
	err := mhb.data.Set(field, value)
	assert.Nil(t, err)
}

func (mhb *MockHelperBase) givenStateFieldSetToWrappedJSON(t *testing.T, field string, i interface{}) {
	err := mhb.data.Set(field, tfmapper.WrapJSON(i))
	assert.Nil(t, err)
}

func (mhb *MockHelperBase) willHaveStateId(t *testing.T, v string) {
	assert.Equal(t, v, mhb.data.Id())
}

func (mhb *MockHelperBase) willHaveStateIdDefined(t *testing.T) {
	assert.True(t, len(mhb.data.Id()) > 3)
}

func (mhb *MockHelperBase) mockClientWill() *MockClient {
	return &mhb.mockCl
}

// --------------------------------------------------------------------------------------------------------------
// DatasourceTemplateMockHelper

type DatasourceTemplateMockHelper[ParentIdent any, Ident any, MType any] struct {
	MockHelperBase
	template DatasourceTemplate[ParentIdent, Ident, MType]
}

func (dtmh *DatasourceTemplateMockHelper[ParentIdent, Ident, MType]) thenExecutingDataSourceQuery(t *testing.T) {
	dg := dtmh.template.Query(context.TODO(), dtmh.data, &dtmh.mockCl)
	assert.True(t, len(dg) == 0)

	dtmh.mockCl.AssertExpectations(t)
}

func (dtmh *DatasourceTemplateMockHelper[ParentIdent, Ident, MType]) thenExecutingDataSourceQueryWillYieldDiagnostic(t *testing.T, dgText string) {
	dg := dtmh.template.Query(context.TODO(), dtmh.data, &dtmh.mockCl)
	assert.True(t, len(dg) == 1)
	assert.Equal(t, diag.Error, dg[0].Severity)
	assert.Equal(t, dgText, dg[0].Summary)

	dtmh.mockCl.AssertExpectations(t)
}

// --------------------------------------------------------------------------------------------------------------
// Static messages

type ResourceTemplateMockHelper[Parent any, Ident any, MTYpe any] struct {
	MockHelperBase
	template *ResourceTemplate[Parent, Ident, MTYpe]
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) givenIdentity(t *testing.T, ident Ident) {
	err := rthm.template.Mapper.AssignIdentity(ident, rthm.data)
	assert.Nil(t, err)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) givenIdentityString(str string) {
	rthm.data.SetId(str)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) givenParentIdentity(t *testing.T, ident Parent) {
	err := rthm.template.Mapper.TestSetPrentIdentity(ident, rthm.data)
	assert.Nil(t, err)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingCreate(t *testing.T) {
	dg := rthm.template.Create(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertEmptyDiagnostic(t, dg)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingRead(t *testing.T) {
	dg := rthm.template.Read(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertEmptyDiagnostic(t, dg)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingReadWillYieldDiagnostic(t *testing.T, text string) {
	dg := rthm.template.Read(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertDiagnostic(t, dg, text)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingUpdate(t *testing.T) {
	dg := rthm.template.Update(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertEmptyDiagnostic(t, dg)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingUpdateWillYieldDiagnostic(t *testing.T, text string) {
	dg := rthm.template.Update(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertDiagnostic(t, dg, text)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingDelete(t *testing.T) {
	dg := rthm.template.Delete(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertEmptyDiagnostic(t, dg)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) assertEmptyDiagnostic(t *testing.T, dg diag.Diagnostics) {
	assert.True(t, len(dg) == 0, "Received diagnostics: %s", dg)
	rthm.mockCl.AssertExpectations(t)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingDeleteWillYieldDiagnostic(t *testing.T, text string) {
	dg := rthm.template.Delete(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertDiagnostic(t, dg, text)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) assertDiagnostic(t *testing.T, dg diag.Diagnostics, text string) {
	assert.True(t, len(dg) == 1)
	assert.Equal(t, diag.Error, dg[0].Severity)
	assert.Equal(t, text, dg[0].Summary)

	rthm.mockCl.AssertExpectations(t)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenExecutingCreateWillYieldDiagnostic(t *testing.T, text string) {
	dg := rthm.template.Create(context.TODO(), rthm.data, &rthm.mockCl)
	rthm.assertDiagnostic(t, dg, text)
}

type IdentityValidator[Ident any] func(t *testing.T, id Ident)

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenAssignedIdIs(t *testing.T, f IdentityValidator[Ident]) {
	ident, err := rthm.template.Mapper.Identity(rthm.data)
	assert.Nil(t, err)
	f(t, ident)
}

func (rthm *DatasourceTemplateMockHelper[Parent, Ident, MTYpe]) thenAssignedIdIs(t *testing.T, f IdentityValidator[Ident]) {
	ident, err := rthm.template.DatasourceMapper().Identity(rthm.data)
	assert.Nil(t, err)
	f(t, ident)
}

func (rthm *ResourceTemplateMockHelper[Parent, Ident, MTYpe]) thenAssignedIdIsEmpty(t *testing.T) {
	assert.Equal(t, "", rthm.data.Id())
}

// --------------------------------------------------------------------------------------------------------------
// Static messages.

func CreateTestDatasource[ParentIdent any, Ident any, MType any](tmpl DatasourceTemplate[ParentIdent, Ident, MType]) *DatasourceTemplateMockHelper[ParentIdent, Ident, MType] {
	rv := DatasourceTemplateMockHelper[ParentIdent, Ident, MType]{
		template: tmpl,
		MockHelperBase: MockHelperBase{
			schema: tmpl.DataSourceSchema(),
			data:   tmpl.TestState(),
			mockCl: MockClient{},
		},
	}

	return &rv
}

func CreateTestResource[Parent any, Ident any, MType any](tmpl *ResourceTemplate[Parent, Ident, MType]) *ResourceTemplateMockHelper[Parent, Ident, MType] {
	rv := ResourceTemplateMockHelper[Parent, Ident, MType]{
		template: tmpl,
		MockHelperBase: MockHelperBase{
			schema: tmpl.ResourceSchema(),
			data:   tmpl.TestState(),
			mockCl: MockClient{},
		},
	}

	return &rv
}
