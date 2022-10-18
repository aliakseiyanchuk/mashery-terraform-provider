package mashery_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/mock"
)

type MasheryPlanMethodMockClient struct {
	v3client.Client
	mock.Mock
}

func (mock *MasheryPlanMethodMockClient) CreateService(ctx context.Context, service masherytypes.Service) (*masherytypes.Service, error) {
	args := mock.Called(ctx, service)
	return args.Get(0).(*masherytypes.Service), args.Error(1)
}
