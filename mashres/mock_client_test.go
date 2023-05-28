package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	v3client.Client
	mock.Mock
}

func (mc *MockClient) ListRolesFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.Role, error) {
	args := mc.Called(ctx, params, fields)
	return args.Get(0).([]masherytypes.Role), args.Error(1)
}

func (mc *MockClient) ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string, fields []string) ([]masherytypes.EmailTemplateSet, error) {
	args := mc.Called(ctx, params, fields)
	return args.Get(0).([]masherytypes.EmailTemplateSet), args.Error(1)
}

func (mc *MockClient) CreateService(ctx context.Context, service masherytypes.Service) (*masherytypes.Service, error) {
	return servicePointerWithError(mc.Called(ctx, service))
}

func servicePointerWithError(args mock.Arguments) (*masherytypes.Service, error) {
	var rv *masherytypes.Service = nil
	if rawRw := args.Get(0); rawRw != nil {
		rv = rawRw.(*masherytypes.Service)
	}

	return rv, args.Error(1)
}

func (mc *MockClient) SetServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, roles []masherytypes.RolePermission) error {
	args := mc.Called(ctx, id, roles)
	return args.Error(0)
}

func (mc *MockClient) DeleteService(ctx context.Context, serviceId masherytypes.ServiceIdentifier) error {
	args := mc.Called(ctx, serviceId)
	return args.Error(0)
}

func (mc *MockClient) CountEndpointsOf(ctx context.Context, serviceId masherytypes.ServiceIdentifier) (int64, error) {
	args := mc.Called(ctx, serviceId)
	return args.Get(0).(int64), args.Error(1)
}

func (mc *MockClient) GetService(ctx context.Context, id masherytypes.ServiceIdentifier) (*masherytypes.Service, error) {
	return servicePointerWithError(mc.Called(ctx, id))
}

func (mc *MockClient) GetServiceRoles(ctx context.Context, serviceId masherytypes.ServiceIdentifier) ([]masherytypes.RolePermission, error) {
	args := mc.Called(ctx, serviceId)
	return args.Get(0).([]masherytypes.RolePermission), args.Error(1)
}

func (mc *MockClient) DeleteServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier) error {
	args := mc.Called(ctx, id)
	return args.Error(0)
}

func (mc *MockClient) UpdateService(ctx context.Context, service masherytypes.Service) (*masherytypes.Service, error) {
	return servicePointerWithError(mc.Called(ctx, service))
}
