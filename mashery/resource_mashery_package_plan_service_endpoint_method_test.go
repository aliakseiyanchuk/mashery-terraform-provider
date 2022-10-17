package mashery_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashery"
	"terraform-provider-mashery/mashschema"
	"testing"

	mock "github.com/stretchr/testify/mock"
)

type MasheryPlanMethodMockClient struct {
	v3client.PluggableClient
	mock.Mock
}

func (m *MasheryPlanMethodMockClient) DeletePackagePlanMethod(ctx context.Context, id masherytypes.MasheryPlanServiceEndpointMethod) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestDeletePlanMethod(t *testing.T) {
	d := mashschema.NewResourceData(&mashschema.PlanMethodSchema)
	d.SetId("packageId::planId::serviceId::endpointId::methodGUID")

	expId := masherytypes.MasheryPlanServiceEndpointMethod{
		MasheryPlanServiceEndpoint: masherytypes.MasheryPlanServiceEndpoint{
			MasheryPlanService: masherytypes.MasheryPlanService{
				PackageId: "packageId",
				PlanId:    "planId",
				ServiceId: "serviceId",
			},
			EndpointId: "endpointId",
		},
		MethodId: "methodGUID",
	}

	tm := MasheryPlanMethodMockClient{}
	tm.On("DeletePackagePlanMethod", mock.Anything, expId).Return(nil)

	dgs := mashery.DeletePlanMethod(context.TODO(), d, &tm)
	mashschema.LogErrorDiagnostics(t, "DeletePlanMethod", &dgs)

	tm.AssertExpectations(t)
}
