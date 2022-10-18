package mashery_test

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashery"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestWillReturnDiagnosticIfQueryFails(t *testing.T) {

	mockClient := &MasheryPlanMethodMockClient{}
	dummyMapper := mashschema.DummyDataSourceMapper
	resourceOperation := &mashery.DatasourceTemplate{
		Mapper: dummyMapper,
		Query: func(ctx context.Context, cl v3client.Client, query map[string]string) ([]interface{}, error) {
			return nil, errors.New("mock error")
		},
	}

	d := dummyMapper.TestResourceData()

	dg := resourceOperation.TemplateQuery(context.TODO(), d, mockClient)
	assert.Equal(t, 1, len(dg))
	assert.Equal(t, diag.Error, dg[0].Severity)
	assert.Equal(t, "querying dummy object failed", dg[0].Summary)
	assert.Equal(t, "querying dummy objects using parameters map[] has failed: mock error", dg[0].Detail)

}
