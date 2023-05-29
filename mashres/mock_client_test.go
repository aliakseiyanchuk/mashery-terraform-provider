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

func (mc *MockClient) CreateServiceOAuthSecurityProfile(ctx context.Context, service masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, error) {
	args := mc.Called(ctx, service)

	var rv *masherytypes.MasheryOAuth
	if rawRV := args.Get(0); rawRV != nil {
		rv = rawRV.(*masherytypes.MasheryOAuth)
	}

	return rv, args.Error(1)
}

func (mc *MockClient) GetServiceCache(ctx context.Context, id masherytypes.ServiceIdentifier) (*masherytypes.ServiceCache, error) {
	args := mc.Called(ctx, id)

	return serviceCacheAndErrorFrom(args)
}

func serviceCacheAndErrorFrom(args mock.Arguments) (*masherytypes.ServiceCache, error) {
	var rv *masherytypes.ServiceCache
	if rawRV := args.Get(0); rawRV != nil {
		rv = rawRV.(*masherytypes.ServiceCache)
	}

	return rv, args.Error(1)
}

func (mc *MockClient) CreateServiceCache(ctx context.Context, id masherytypes.ServiceIdentifier, service masherytypes.ServiceCache) (*masherytypes.ServiceCache, error) {
	args := mc.Called(ctx, id, service)
	return serviceCacheAndErrorFrom(args)
}

func (mc *MockClient) CreateEndpoint(ctx context.Context, serviceId masherytypes.ServiceIdentifier, endp masherytypes.Endpoint) (*masherytypes.Endpoint, error) {
	args := mc.Called(ctx, serviceId, endp)

	var rv *masherytypes.Endpoint = nil
	if args.Get(0) != nil {
		rv = args.Get(0).(*masherytypes.Endpoint)
	}

	return rv, args.Error(1)
}

func (mc *MockClient) CreatePackage(ctx context.Context, pack masherytypes.Package) (*masherytypes.Package, error) {
	args := mc.Called(ctx, pack)

	var rv *masherytypes.Package = nil
	if args.Get(0) != nil {
		rv = args.Get(0).(*masherytypes.Package)
	}

	return rv, args.Error(1)
}
func (mc *MockClient) CreatePlan(ctx context.Context, packageId masherytypes.PackageIdentifier, plan masherytypes.Plan) (*masherytypes.Plan, error) {
	args := mc.Called(ctx, packageId, plan)

	var rv *masherytypes.Plan = nil
	if args.Get(0) != nil {
		rv = args.Get(0).(*masherytypes.Plan)
	}

	return rv, args.Error(1)
}
