package mashery_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashery"
	"testing"

	mock "github.com/stretchr/testify/mock"
)

type MasheryPlanMethodMockClient struct {
	v3client.PluggableClient
	mock.Mock
}

func (m *MasheryPlanMethodMockClient) DeletePackagePlanMethod(ctx context.Context, id v3client.MasheryPlanServiceEndpointMethod) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestDeletePlanMethod(t *testing.T) {
	d := NewResourceData(&mashery.PlanMethodSchema)
	d.SetId("packageId::planId::serviceId::endpointId::methodGUID")

	expId := v3client.MasheryPlanServiceEndpointMethod{
		MasheryPlanServiceEndpoint: v3client.MasheryPlanServiceEndpoint{
			MasheryPlanService: v3client.MasheryPlanService{
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
	LogErrorDiagnostics(t, "DeletePlanMethod", &dgs)

	tm.AssertExpectations(t)
}
